package surfstore

import (
	context "context"
	"fmt"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type BlockStore struct {
	BlockMap map[string]*Block
	UnimplementedBlockStoreServer
}

func (bs *BlockStore) GetBlock(ctx context.Context, blockHash *BlockHash) (*Block, error) {
	b, ok := bs.BlockMap[blockHash.Hash]
	if !ok {
		return nil, fmt.Errorf("GetBlock error")
	}
	return b, nil
}

func (bs *BlockStore) PutBlock(ctx context.Context, block *Block) (*Success, error) {
	h := GetBlockHashString(block.BlockData[:block.BlockSize])
	bs.BlockMap[h] = block
	return &Success{Flag: true}, nil
}

// Given a list of hashes “in”, returns a list containing the
// hashes that are not stored in the key-value store
func (bs *BlockStore) MissingBlocks(ctx context.Context, blockHashesIn *BlockHashes) (*BlockHashes, error) {
	res := BlockHashes{Hashes: []string{}}
	for _, h := range blockHashesIn.Hashes {
		_, ok := bs.BlockMap[h]
		if !ok {
			res.Hashes = append(res.Hashes, h)
		}
	}
	return &res, nil
}

// Return a list containing all blockHashes on this block server
func (bs *BlockStore) GetBlockHashes(ctx context.Context, _ *emptypb.Empty) (*BlockHashes, error) {
	hashes := []string{}
	for hash := range bs.BlockMap {
		hashes = append(hashes, hash)
	}
	return &BlockHashes{Hashes: hashes}, nil
}

// This line guarantees all method for BlockStore are implemented
var _ BlockStoreInterface = new(BlockStore)

func NewBlockStore() *BlockStore {
	return &BlockStore{
		BlockMap: map[string]*Block{},
	}
}
