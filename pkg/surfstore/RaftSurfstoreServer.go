package surfstore

import (
	context "context"
	"log"
	"sync"
	"time"

	grpc "google.golang.org/grpc"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// TODO Add fields you need here
type RaftSurfstore struct {
	serverStatus      ServerStatus
	serverStatusMutex *sync.RWMutex
	sendHeartbestMutex *sync.RWMutex
	term              int64
	log               []*UpdateOperation
	id                int64
	metaStore         *MetaStore
	commitIndex       int64

	raftStateMutex *sync.RWMutex

	rpcConns   []*grpc.ClientConn
	grpcServer *grpc.Server

	//New Additions
	peers           []string
	pendingRequests []*chan PendingRequest
	lastApplied     int64
	nextIndex       []int64

	/*--------------- Chaos Monkey --------------*/
	unreachableFrom map[int64]bool
	UnimplementedRaftSurfstoreServer
}

// 1. Append cmd to log
// 2. Issue AppendEntries
// 3. If majority replied success, apply to state machine (commited)
// 4. Return result to cliemt
func (s *RaftSurfstore) GetFileInfoMap(ctx context.Context, empty *emptypb.Empty) (*FileInfoMap, error) {
	// Ensure that the majority of servers are up
	if err := s.checkStatus(); err != nil {
		return nil, err
	}
	s.sendPersistentHeartbeats(ctx)
	return s.metaStore.GetFileInfoMap(ctx, empty)
}

func (s *RaftSurfstore) GetBlockStoreMap(ctx context.Context, hashes *BlockHashes) (*BlockStoreMap, error) {
	// Ensure that the majority of servers are up
	if err := s.checkStatus(); err != nil {
		return nil, err
	}
	s.sendPersistentHeartbeats(ctx)
	return s.metaStore.GetBlockStoreMap(ctx, hashes)
}

func (s *RaftSurfstore) GetBlockStoreAddrs(ctx context.Context, empty *emptypb.Empty) (*BlockStoreAddrs, error) {
	// Ensure that the majority of servers are up
	if err := s.checkStatus(); err != nil {
		return nil, err
	}
	s.sendPersistentHeartbeats(ctx)
	return s.metaStore.GetBlockStoreAddrs(ctx, empty)
}

func (s *RaftSurfstore) UpdateFile(ctx context.Context, filemeta *FileMetaData) (*Version, error) {
	// Ensure that the request gets replicated on majority of the servers.
	// Commit the entries and then apply to the state machine
	log.Println("Server", s.id, "[UpdateFile]", filemeta)
	if err := s.checkStatus(); err != nil {
		return nil, err
	}
	entry := UpdateOperation{
		Term:         s.term,
		FileMetaData: filemeta,
	}
	s.log = append(s.log, &entry)

	// for simplicity, make updatefile blocking
	s.sendPersistentHeartbeats(ctx)

	// Ensure that leader commits first and then applies to the state machine
	// TODO: need to commit multiple cases
	s.commitIndex = int64(len(s.log)) - 1
	s.lastApplied = int64(len(s.log)) - 1

	if entry.FileMetaData != nil {
		return s.metaStore.UpdateFile(ctx, entry.FileMetaData)
	}
	// called by setLeader
	return nil, nil
}

func (s *RaftSurfstore) AppendEntries(ctx context.Context, input *AppendEntryInput) (*AppendEntryOutput, error) {
	//check the status
	s.raftStateMutex.RLock()
	sTerm := s.term
	sStatus := s.serverStatus
	sID := s.id
	unreachable, ok := s.unreachableFrom[input.LeaderId]
	if !ok {
		unreachable = false
	}
	s.raftStateMutex.RUnlock()

	// 0. If server is crashed or unreachable
	if sStatus == ServerStatus_CRASHED {
		log.Println("Server", sID, "[AppendEntries] Server Crashed")
		return nil, status.Error(500, ErrServerCrashed.Error())
	}
	if unreachable {
		log.Println("Server", sID, "[AppendEntries] Server Unreachable")
		return nil, status.Error(501, ErrServerCrashedUnreachable.Error())
	}
	appendEntriesOutput := AppendEntryOutput{
		Term:         sTerm,
		ServerId:     sID,
		Success:      true,
		MatchedIndex: -1,
	}

	// 1. Reply false if term < currentTerm (§5.1)
	if input.Term < sTerm {
		log.Println("Server", sID, "[AppendEntries] Sending output:", "Term", appendEntriesOutput.Term, "Id", appendEntriesOutput.ServerId, "Success", appendEntriesOutput.Success, "Matched Index", appendEntriesOutput.MatchedIndex)
		appendEntriesOutput.Success = false
		return &appendEntriesOutput, nil
	}
	// 6. If received from other legitimate leader, waive leadership
	if input.Term > sTerm {
		s.serverStatusMutex.Lock()
		s.serverStatus = ServerStatus_FOLLOWER
		s.serverStatusMutex.Unlock()

		s.raftStateMutex.Lock()
		s.term = input.Term
		appendEntriesOutput.Term = input.Term
		s.raftStateMutex.Unlock()
	}
	// 2. Reply false if log doesn’t contain an entry at prevLogIndex or whose term
	// doesn't match prevLogTerm (§5.3)
	s.raftStateMutex.RLock()
	if int(input.PrevLogIndex) >= len(s.log) {
		appendEntriesOutput.Success = false
	} else if input.PrevLogIndex != -1 {
		prevLogTerm := s.log[input.PrevLogIndex].Term
		if prevLogTerm != input.PrevLogTerm {
			appendEntriesOutput.Success = false
		}
	}
	s.raftStateMutex.RUnlock()
	if !appendEntriesOutput.Success {
		log.Println("Server", sID, "[AppendEntries] Sending output:", "Term", appendEntriesOutput.Term, "Id", appendEntriesOutput.ServerId, "Success", appendEntriesOutput.Success, "Matched Index", appendEntriesOutput.MatchedIndex)
		return &appendEntriesOutput, nil
	}
	// 3. If an existing entry conflicts with a new one (same index but different
	// terms), delete the existing entry and all that follow it (§5.3)
	// 4. Append any new entries not already in the log
	s.raftStateMutex.Lock()
	s.log = append(s.log[:input.PrevLogIndex+1], input.Entries...)
	// log.Printf("Server %v [AppendEntries] log updated to %v", s.id, s.log)

	// 5. If leaderCommit > commitIndex, set commitIndex = min(leaderCommit, index
	// of last new entry)
	s.commitIndex = input.LeaderCommit
	for s.lastApplied < s.commitIndex {
		entry := s.log[s.lastApplied+1]
		if entry.FileMetaData != nil {
			log.Println("Server", sID, "[AppendEntries] apply update", entry.FileMetaData)
			_, err := s.metaStore.UpdateFile(ctx, entry.FileMetaData)
			if err != nil {
				s.raftStateMutex.Unlock()
				return nil, err
			}
		}
		s.lastApplied += 1
	}
	log.Println("Server", sID, "[AppendEntries] Sending output:", "Term", appendEntriesOutput.Term, "Id", appendEntriesOutput.ServerId, "Success", appendEntriesOutput.Success, "Matched Index", appendEntriesOutput.MatchedIndex)
	s.raftStateMutex.Unlock()

	return &appendEntriesOutput, nil
}

// 1. term++
// 2. add noop to log
func (s *RaftSurfstore) SetLeader(ctx context.Context, _ *emptypb.Empty) (*Success, error) {
	s.serverStatusMutex.RLock()
	serverStatus := s.serverStatus
	s.serverStatusMutex.RUnlock()

	if serverStatus == ServerStatus_CRASHED {
		return &Success{Flag: false}, ErrServerCrashed
	}

	s.serverStatusMutex.Lock()
	s.serverStatus = ServerStatus_LEADER
	log.Printf("Server %d has been set as leader", s.id)
	s.serverStatusMutex.Unlock()

	s.raftStateMutex.Lock()
	s.term += 1
	lastLogIndex := len(s.log)
	s.nextIndex = make([]int64, len(s.peers))
	for i := range s.nextIndex {
		s.nextIndex[i] = int64(lastLogIndex)
	}

	s.raftStateMutex.Unlock()
	s.UpdateFile(ctx, nil)
	return &Success{Flag: true}, nil
}

// Used to establish leadership / replicate / commit
// If no enough votes, change back to follower
func (s *RaftSurfstore) SendHeartbeat(ctx context.Context, _ *emptypb.Empty) (*Success, error) {
	if err := s.checkStatus(); err != nil {
		return nil, err
	}

	s.sendPersistentHeartbeats(ctx)
	if err := s.checkStatus(); err != nil {
		return nil, err
	}
	return &Success{Flag: true}, nil
}

func (s *RaftSurfstore) majorityUpdated(peerUpdateStatuses []PeerUpdateStatus) bool {
	updatedCnt := 0
	for _, peerUpdateStatus := range peerUpdateStatuses {
		if peerUpdateStatus == PeerUnknown || peerUpdateStatus == PeerUpdating {
			return false
		} else if peerUpdateStatus == PeerUpdated {
			updatedCnt += 1
		}
	}
	return updatedCnt > len(peerUpdateStatuses)/2
}

func (s *RaftSurfstore) sendPersistentHeartbeats(ctx context.Context) {
	numServers := len(s.peers)
	peerResponses := make(chan PeerStatusResponse)
	for idx := range s.peers {
		idx := int64(idx)
		if idx == s.id {
			continue
		}
		go s.sendToFollower(ctx, idx, peerResponses)
	}

	peerUpdateStatuses := make([]PeerUpdateStatus, numServers)
	for i := range peerUpdateStatuses {
		peerUpdateStatuses[i] = PeerUnknown
	}
	peerUpdateStatuses[s.id] = PeerUpdated
	// proceed only if majority updated and no unknown/updating
	for !s.majorityUpdated(peerUpdateStatuses) {
		if err := s.checkStatus(); err != nil {
			return
		}
		response := <-peerResponses
		peerUpdateStatuses[response.id] = response.status
	}
}

// send log[nextIndex:len-1]
func (s *RaftSurfstore) sendToFollower(ctx context.Context, peerId int64, peerResponses chan<- PeerStatusResponse) {
	client := NewRaftSurfstoreClient(s.rpcConns[peerId])
	for {
		if err := s.checkStatus(); err != nil {
			return
		}
		nextIndex := s.nextIndex[peerId] // initial: 0
		prevLogTerm := int64(0)
		if nextIndex > 0 {
			prevLogTerm = s.log[nextIndex-1].Term
		}

		appendEntriesInput := AppendEntryInput{
			Term:         s.term,
			LeaderId:     s.id,
			PrevLogTerm:  prevLogTerm,
			PrevLogIndex: nextIndex - 1,
			Entries:      s.log[nextIndex:len(s.log)],
			LeaderCommit: s.commitIndex,
		}
		log.Printf("Server %v [sendToFollower %v] term: %v, LeaderId: %v, PrevLogTerm: %v, PrevLogIndex: %v, Entries: %v, LeaderCommit: %v", s.id, peerId, appendEntriesInput.Term, appendEntriesInput.LeaderId, appendEntriesInput.PrevLogTerm, appendEntriesInput.PrevLogIndex, appendEntriesInput.Entries, appendEntriesInput.LeaderCommit)

		reply, err := client.AppendEntries(context.Background(), &appendEntriesInput)
		// log.Println("Server", s.id, "[sendToFollower] Receiving output", "Term:", reply.Term, "Id:", reply.ServerId, "Success:", reply.Success, "Matched Index:", reply.MatchedIndex)

		// response processing
		if err != nil {
			// TODO: fix this
			st, _ := status.FromError(err)
			if st.Message() == ErrServerCrashed.Error() || st.Message() == ErrServerCrashedUnreachable.Error() {
				peerResponses <- PeerStatusResponse{id: peerId, status: PeerUnreachable}
			} else {
				log.Println("Error is", err.Error())
				panic("Should not happen")
			}
		} else if reply.Term > s.term {
			// find a new term, become follower
			s.term = reply.Term
			s.serverStatus = ServerStatus_FOLLOWER
			peerResponses <- PeerStatusResponse{id: peerId, status: PeerUpdated}
			break
		} else if !reply.Success {
			// try again with smaller nextIndex
			s.nextIndex[peerId] -= 1
			peerResponses <- PeerStatusResponse{id: peerId, status: PeerUpdating}
		} else {
			s.nextIndex[peerId] = int64(len(s.log))
			peerResponses <- PeerStatusResponse{id: peerId, status: PeerUpdated}
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// ========== DO NOT MODIFY BELOW THIS LINE =====================================

func (s *RaftSurfstore) MakeServerUnreachableFrom(ctx context.Context, servers *UnreachableFromServers) (*Success, error) {
	s.raftStateMutex.Lock()
	if len(servers.ServerIds) == 0 {
		s.unreachableFrom = make(map[int64]bool)
		log.Printf("Server %d is reachable from all servers", s.id)
	} else {
		for _, serverId := range servers.ServerIds {
			s.unreachableFrom[serverId] = true
		}
		log.Printf("Server %d is unreachable from %v", s.id, s.unreachableFrom)
	}

	s.raftStateMutex.Unlock()

	return &Success{Flag: true}, nil
}

func (s *RaftSurfstore) Crash(ctx context.Context, _ *emptypb.Empty) (*Success, error) {
	s.serverStatusMutex.Lock()
	s.serverStatus = ServerStatus_CRASHED
	log.Printf("Server %d is crashed", s.id)
	s.serverStatusMutex.Unlock()

	return &Success{Flag: true}, nil
}

func (s *RaftSurfstore) Restore(ctx context.Context, _ *emptypb.Empty) (*Success, error) {
	s.serverStatusMutex.Lock()
	s.serverStatus = ServerStatus_FOLLOWER
	s.serverStatusMutex.Unlock()

	s.raftStateMutex.Lock()
	s.unreachableFrom = make(map[int64]bool)
	s.raftStateMutex.Unlock()

	log.Printf("Server %d is restored to follower and reachable from all servers", s.id)

	return &Success{Flag: true}, nil
}

func (s *RaftSurfstore) GetInternalState(ctx context.Context, empty *emptypb.Empty) (*RaftInternalState, error) {
	fileInfoMap, _ := s.metaStore.GetFileInfoMap(ctx, empty)
	s.serverStatusMutex.RLock()
	s.raftStateMutex.RLock()
	state := &RaftInternalState{
		Status:      s.serverStatus,
		Term:        s.term,
		CommitIndex: s.commitIndex,
		Log:         s.log,
		MetaMap:     fileInfoMap,
	}
	s.raftStateMutex.RUnlock()
	s.serverStatusMutex.RUnlock()

	return state, nil
}

var _ RaftSurfstoreInterface = new(RaftSurfstore)
