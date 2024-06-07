package surfstore

import (
	context "context"
	"fmt"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type MetaStore struct {
	FileMetaMap        map[string]*FileMetaData
	BlockStoreAddrs    []string
	ConsistentHashRing *ConsistentHashRing
	UnimplementedMetaStoreServer
}

func (m *MetaStore) GetFileInfoMap(ctx context.Context, _ *emptypb.Empty) (*FileInfoMap, error) {
	return &FileInfoMap{FileInfoMap: m.FileMetaMap}, nil
}

func (m *MetaStore) UpdateFile(ctx context.Context, fileMetaData *FileMetaData) (*Version, error) {
	fn := fileMetaData.Filename
	fmd, ok := m.FileMetaMap[fn]
	clientVersion := fileMetaData.Version
	if !ok {
		if clientVersion != 1 {
			return nil, fmt.Errorf("UpdateFile error")
		}
		// create
		m.FileMetaMap[fn] = fileMetaData
	} else if clientVersion == fmd.Version+1 {
		// update
		m.FileMetaMap[fn] = fileMetaData
	} else {
		clientVersion = -1
	}
	return &Version{Version: clientVersion}, nil
}

// Given a list of block hashes, find out which block server they belong to.
// Returns a mapping from block server address to block hashes.
func (m *MetaStore) GetBlockStoreMap(ctx context.Context, blockHashesIn *BlockHashes) (*BlockStoreMap, error) {
	blockStoreMap := BlockStoreMap{BlockStoreMap: make(map[string]*BlockHashes)}
	for _, blockStoreAddr := range m.BlockStoreAddrs {
		blockStoreMap.BlockStoreMap[blockStoreAddr] = &BlockHashes{Hashes: []string{}}
	}
	for _, blockHash := range blockHashesIn.Hashes {
		responsibleServer := m.ConsistentHashRing.GetResponsibleServer(blockHash)
		blockStoreMap.BlockStoreMap[responsibleServer].Hashes = append(blockStoreMap.BlockStoreMap[responsibleServer].Hashes, blockHash)
	}
	return &blockStoreMap, nil
}

func (m *MetaStore) GetBlockStoreAddrs(ctx context.Context, _ *emptypb.Empty) (*BlockStoreAddrs, error) {
	return &BlockStoreAddrs{BlockStoreAddrs: m.BlockStoreAddrs}, nil
}

// This line guarantees all method for MetaStore are implemented
var _ MetaStoreInterface = new(MetaStore)

func NewMetaStore(blockStoreAddrs []string) *MetaStore {
	return &MetaStore{
		FileMetaMap:        map[string]*FileMetaData{},
		BlockStoreAddrs:    blockStoreAddrs,
		ConsistentHashRing: NewConsistentHashRing(blockStoreAddrs),
	}
}
