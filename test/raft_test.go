package SurfTest

import (
	context "context"
	"cse224/proj5/pkg/surfstore"
	"fmt"
	"testing"
	"time"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func TestRaftSetLeader(t *testing.T) {
	//Setup
	cfgPath := "./config_files/3nodes.txt"
	test := InitTest(cfgPath)
	defer EndTest(test)

	// TEST
	leaderIdx := 0
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})

	// // heartbeat
	for _, server := range test.Clients {
		server.SendHeartbeat(test.Context, &emptypb.Empty{})
	}

	for idx, server := range test.Clients {
		// all should have the leaders term
		state, _ := server.GetInternalState(test.Context, &emptypb.Empty{})
		if state == nil {
			t.Fatalf("Could not get state")
		}
		if state.Term != int64(1) {
			t.Fatalf("Server %d should be in term %d", idx, 1)
		}
		if idx == leaderIdx {
			// server should be the leader
			if state.Status != surfstore.ServerStatus_LEADER {
				t.Fatalf("Server %d should be the leader", idx)
			}
		} else {
			// server should not be the leader
			if state.Status == surfstore.ServerStatus_LEADER {
				t.Fatalf("Server %d should not be the leader", idx)
			}
		}
	}

	leaderIdx = 2
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})

	// heartbeat
	for _, server := range test.Clients {
		server.SendHeartbeat(test.Context, &emptypb.Empty{})
	}

	for idx, server := range test.Clients {
		// all should have the leaders term
		state, _ := server.GetInternalState(test.Context, &emptypb.Empty{})
		if state == nil {
			t.Fatalf("Could not get state")
		}
		if state.Term != int64(2) {
			t.Fatalf("Server should be in term %d", 2)
		}
		if idx == leaderIdx {
			// server should be the leader
			if state.Status != surfstore.ServerStatus_LEADER {
				t.Fatalf("Server %d should be the leader", idx)
			}
		} else {
			// server should not be the leader
			if state.Status == surfstore.ServerStatus_LEADER {
				t.Fatalf("Server %d should not be the leader", idx)
			}
		}
	}
}

func TestRaftFollowersGetUpdates(t *testing.T) {
	//Setup
	cfgPath := "./config_files/3nodes.txt"
	test := InitTest(cfgPath)
	defer EndTest(test)

	// TEST
	leaderIdx := 0
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	filemeta1 := &surfstore.FileMetaData{
		Filename:      "testFile1",
		Version:       1,
		BlockHashList: nil,
	}

	test.Clients[leaderIdx].UpdateFile(test.Context, filemeta1)
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	goldenMeta := make(map[string]*surfstore.FileMetaData)
	goldenMeta[filemeta1.Filename] = filemeta1

	goldenLog := make([]*surfstore.UpdateOperation, 0)
	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term:         1,
		FileMetaData: nil,
	})
	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term:         1,
		FileMetaData: filemeta1,
	})
	var leader bool
	term := int64(1)

	for idx, server := range test.Clients {
		if idx == leaderIdx {
			leader = bool(true)
		} else {
			leader = bool(false)
		}
		_, err := CheckInternalState(&leader, &term, goldenLog, goldenMeta, server, test.Context)
		if err != nil {
			t.Fatalf("Error checking state for server %d: %s", idx, err.Error())
		}
	}
}

func TestRaftNetworkPartitionWithConcurrentRequests(t *testing.T) {
	t.Log("leader1 gets 1 request while the majority of the cluster is unreachable. As a result of a (one way) network partition, leader1 ends up with the minority. leader2 from the majority is elected")
	// 	// A B C D E
	cfgPath := "./config_files/5nodes.txt"
	test := InitTest(cfgPath)
	defer EndTest(test)

	filemeta1 := &surfstore.FileMetaData{
		Filename:      "testFile1",
		Version:       1,
		BlockHashList: nil,
	}
	filemeta2 := &surfstore.FileMetaData{
		Filename:      "testFile2",
		Version:       1,
		BlockHashList: nil,
	}

	A := 0
	C := 2
	// D := 3
	E := 4

	// A is leader
	leaderIdx := A
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	// Partition happens
	// C D E are unreachable
	unreachableFromServers := &surfstore.UnreachableFromServers{
		ServerIds: []int64{0, 1},
	}
	for i := C; i <= E; i++ {
		test.Clients[i].MakeServerUnreachableFrom(test.Context, unreachableFromServers)
	}

	blockChan := make(chan bool)

	// A gets an entry and pushes to A and B
	go func() {
		// This should block though and fail to commit when getting the RPC response from the new leader
		_, _ = test.Clients[leaderIdx].UpdateFile(test.Context, filemeta1)
		blockChan <- false
	}()

	go func() {
		<-time.NewTimer(10 * time.Second).C
		blockChan <- true
	}()

	if !(<-blockChan) {
		t.Fatalf("Request did not block")
	}

	// C becomes leader
	leaderIdx = C
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})
	// C D E are restored
	for i := C; i <= E; i++ {
		test.Clients[i].MakeServerUnreachableFrom(test.Context, &surfstore.UnreachableFromServers{ServerIds: []int64{}})
	}

	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})
	time.Sleep(time.Second)

	// Every node should have an empty log
	goldenLog := make([]*surfstore.UpdateOperation, 0)
	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term:         1,
		FileMetaData: nil,
	})
	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term:         2,
		FileMetaData: nil,
	})
	// Leaders should not commit entries that were created in previous terms.
	goldenMeta := make(map[string]*surfstore.FileMetaData)
	term := int64(2)

	for idx, server := range test.Clients {

		_, err := CheckInternalState(nil, &term, goldenLog, goldenMeta, server, test.Context)
		if err != nil {
			t.Fatalf("Error checking state for server %d: %s", idx, err.Error())
		}
	}

	go func() {
		test.Clients[leaderIdx].UpdateFile(test.Context, filemeta1)
		test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})
	}()
	go func() {
		test.Clients[leaderIdx].UpdateFile(test.Context, filemeta2)
		test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})
	}()

	//Enough time for things to settle
	time.Sleep(2 * time.Second)
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	goldenMeta[filemeta1.Filename] = filemeta1
	goldenMeta[filemeta2.Filename] = filemeta2

	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term:         2,
		FileMetaData: filemeta1,
	})
	goldenLog = append(goldenLog, &surfstore.UpdateOperation{
		Term:         2,
		FileMetaData: filemeta2,
	})

	for idx, server := range test.Clients {

		state, err := server.GetInternalState(test.Context, &emptypb.Empty{})
		if err != nil {
			t.Fatalf("could not get internal state: %s", err.Error())
		}
		if state == nil {
			t.Fatalf("state is nil")
		}

		if len(state.Log) != 4 {
			t.Fatalf("Should have 4 logs")
		}

		if err != nil {
			t.Fatalf("Error checking state for server %d: %s", idx, err.Error())
		}
		fmt.Println(idx, server)
		fmt.Println(state)
		for _, l := range state.Log {
			fmt.Println(l)
		}
	}
}

func TestRaftNetworkDuplicate(t *testing.T) {
	// 	// A B C D E
	cfgPath := "./config_files/5nodes.txt"
	test := InitTest(cfgPath)
	defer EndTest(test)

	filemeta1 := &surfstore.FileMetaData{
		Filename:      "testFile1",
		Version:       1,
		BlockHashList: nil,
	}
	filemeta2 := &surfstore.FileMetaData{
		Filename:      "testFile2",
		Version:       1,
		BlockHashList: nil,
	}

	A := 0
	C := 2
	// D := 3
	E := 4

	// A is leader
	leaderIdx := A
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	// Partition happens
	// C D E are unreachable
	unreachableFromServers := &surfstore.UnreachableFromServers{
		ServerIds: []int64{0, 1},
	}
	for i := 3; i <= E; i++ {
		test.Clients[i].MakeServerUnreachableFrom(test.Context, unreachableFromServers)
	}

	blockChan := make(chan bool)

	// A gets an entry and pushes to A and B
	go func() {
		// This should block though and fail to commit when getting the RPC response from the new leader
		_, _ = test.Clients[leaderIdx].UpdateFile(test.Context, filemeta1)
		test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})
		blockChan <- false
	}()

	go func() {
		<-time.NewTimer(5 * time.Second).C
		blockChan <- true
	}()

	if !(<-blockChan) {
		fmt.Printf("Request did not block")
	}

	// C becomes leader
	leaderIdx = C
	test.Clients[leaderIdx].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})
	time.Sleep(10 * time.Second)
	// C D E are restored
	for i := 3; i <= E; i++ {
		test.Clients[i].MakeServerUnreachableFrom(test.Context, &surfstore.UnreachableFromServers{ServerIds: []int64{}})
	}

	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})
	time.Sleep(time.Second)

	go func() {
		test.Clients[leaderIdx].UpdateFile(test.Context, filemeta1)
		test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})
	}()
	go func() {
		test.Clients[leaderIdx].UpdateFile(test.Context, filemeta2)
		test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})
	}()

	//Enough time for things to settle
	time.Sleep(2 * time.Second)
	test.Clients[leaderIdx].SendHeartbeat(test.Context, &emptypb.Empty{})

	for idx, server := range test.Clients {
		state, _ := server.GetInternalState(test.Context, &emptypb.Empty{})
		fmt.Println(idx, server)
		fmt.Println(state)
		for _, l := range state.Log {
			fmt.Println(l)
		}
	}
}

func TestSimpleCrash(t *testing.T) {
	// A B C
	cfgPath := "./config_files/3nodes.txt"
	test := InitTest(cfgPath)
	defer EndTest(test)

	// A is leader
	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[0].SendHeartbeat(test.Context, &emptypb.Empty{})

	// C crash
	test.Clients[2].Crash(test.Context, &emptypb.Empty{})

	// A update file, crash before sending heartbeat
	filemeta1 := &surfstore.FileMetaData{
		Filename:      "testFile1",
		Version:       1,
		BlockHashList: nil,
	}
	test.Clients[0].UpdateFile(test.Context, filemeta1)
	// test.Clients[0].SendHeartbeat(test.Context, &emptypb.Empty{})
	time.Sleep(30 * time.Second)
	
	// // C restore
	test.Clients[2].Restore(test.Context, &emptypb.Empty{})
	time.Sleep(time.Second)
	showClients(test.Clients, test.Context)
}

func TestServerCrash(t *testing.T) {
	cfgPath := "./config_files/5nodes.txt"
	test := InitTest(cfgPath)
	defer EndTest(test)

	// A is leader
	test.Clients[0].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[0].SendHeartbeat(test.Context, &emptypb.Empty{})
	time.Sleep(time.Second)

	// C crash
	test.Clients[2].Crash(test.Context, &emptypb.Empty{})

	// A update file, crash before sending heartbeat
	filemeta1 := &surfstore.FileMetaData{
		Filename:      "testFile1",
		Version:       1,
		BlockHashList: nil,
	}
	test.Clients[0].UpdateFile(test.Context, filemeta1)
	test.Clients[0].Crash(test.Context, &emptypb.Empty{})
	time.Sleep(time.Second)

	// C restored
	test.Clients[2].Restore(test.Context, &emptypb.Empty{})

	// B is leader and get request
	test.Clients[1].SetLeader(test.Context, &emptypb.Empty{})
	test.Clients[1].SendHeartbeat(test.Context, &emptypb.Empty{})
	time.Sleep(time.Second)

	filemeta2 := &surfstore.FileMetaData{
		Filename:      "testFile2",
		Version:       1,
		BlockHashList: nil,
	}
	test.Clients[1].UpdateFile(test.Context, filemeta2)
	time.Sleep(time.Second)

	// A restored
	test.Clients[0].Restore(test.Context, &emptypb.Empty{})
	test.Clients[1].SendHeartbeat(test.Context, &emptypb.Empty{})
	time.Sleep(time.Second)

	showClients(test.Clients, test.Context)
}

func showClients(clients []surfstore.RaftSurfstoreClient, ctx context.Context) {
	for idx, server := range clients {
		state, _ := server.GetInternalState(ctx, &emptypb.Empty{})
		fmt.Println("Server: ", idx, "=================================")
		fmt.Println("Status: ", state.Status)
		fmt.Println("CommitIndex: ", state.CommitIndex)
		fmt.Println("Term: ", state.Term)
		fmt.Println("Log: {")
		for _, l := range state.Log {
			fmt.Println("  ", l)
		}
		fmt.Println("}")
	}
}