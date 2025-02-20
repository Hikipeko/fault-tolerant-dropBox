package surfstore

import (
	context "context"
	"errors"
	"fmt"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type RPCClient struct {
	MetaStoreAddrs []string
	BaseDir        string
	BlockSize      int
}

func (surfClient *RPCClient) GetBlock(blockHash string, blockStoreAddr string, block *Block) error {
	// connect to the server
	conn, err := grpc.Dial(blockStoreAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	c := NewBlockStoreClient(conn)

	// perform the call
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	b, err := c.GetBlock(ctx, &BlockHash{Hash: blockHash})
	if err != nil {
		conn.Close()
		return err
	}
	block.BlockData = b.BlockData
	block.BlockSize = b.BlockSize

	// close the connection
	return conn.Close()
}

func (surfClient *RPCClient) PutBlock(block *Block, blockStoreAddr string, succ *bool) error {
	// connect to the server
	conn, err := grpc.Dial(blockStoreAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	c := NewBlockStoreClient(conn)

	// perform the call
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	s, err := c.PutBlock(ctx, block)
	if err != nil {
		conn.Close()
		return err
	}
	*succ = s.Flag
	return conn.Close()
}

func (surfClient *RPCClient) MissingBlocks(blockHashesIn []string, blockStoreAddr string, blockHashesOut *[]string) error {
	// connect to the server
	conn, err := grpc.Dial(blockStoreAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	c := NewBlockStoreClient(conn)

	// perform the call
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	bh, err := c.MissingBlocks(ctx, &BlockHashes{Hashes: blockHashesIn})
	if err != nil {
		conn.Close()
		return err
	}
	*blockHashesOut = bh.Hashes
	return conn.Close()
}

func (surfClient *RPCClient) GetBlockHashes(blockStoreAddr string, blockHashes *[]string) error {
	conn, err := grpc.Dial(blockStoreAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	c := NewBlockStoreClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	hashes, err := c.GetBlockHashes(ctx, &emptypb.Empty{})
	if err != nil {
		conn.Close()
		return err
	}
	*blockHashes = hashes.Hashes
	return conn.Close()
}

func (surfClient *RPCClient) GetFileInfoMap(serverFileInfoMap *map[string]*FileMetaData) error {
	for _, server := range surfClient.MetaStoreAddrs {
		conn, err := grpc.Dial(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		c := NewRaftSurfstoreClient(conn)

		// perform the call
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		m, err := c.GetFileInfoMap(ctx, &emptypb.Empty{})

		//Handle errors appropriately
		if err != nil {
			if errors.Is(err, ErrNotLeader) || errors.Is(err, ErrServerCrashed) || errors.Is(err, ErrServerCrashedUnreachable) {
				conn.Close()
				continue
			}
			conn.Close()
			return err
		}
		*serverFileInfoMap = m.FileInfoMap
		return conn.Close()
	}
	return fmt.Errorf("could not find a leader")
}

func (surfClient *RPCClient) UpdateFile(fileMetaData *FileMetaData, latestVersion *int32) error {
	for _, server := range surfClient.MetaStoreAddrs {
		conn, err := grpc.Dial(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		c := NewRaftSurfstoreClient(conn)

		// perform the call
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		v, err := c.UpdateFile(ctx, fileMetaData)

		if err != nil {
			if errors.Is(err, ErrNotLeader) || errors.Is(err, ErrServerCrashed) || errors.Is(err, ErrServerCrashedUnreachable) {
				conn.Close()
				continue
			}
			conn.Close()
			return err
		}

		*latestVersion = v.Version
		return conn.Close()
	}
	return fmt.Errorf("could not find a leader")
}

func (surfClient *RPCClient) GetBlockStoreMap(blockHashesIn []string, blockStoreMap *map[string][]string) error {
	for _, server := range surfClient.MetaStoreAddrs {
		conn, err := grpc.Dial(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		c := NewRaftSurfstoreClient(conn)

		// perform the call
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		bsm, err := c.GetBlockStoreMap(ctx, &BlockHashes{Hashes: blockHashesIn})

		if err != nil {
			if errors.Is(err, ErrNotLeader) || errors.Is(err, ErrServerCrashed) || errors.Is(err, ErrServerCrashedUnreachable) {
				conn.Close()
				continue
			}
			conn.Close()
			return err
		}

		for blockServer, blockHashes := range bsm.BlockStoreMap {
			(*blockStoreMap)[blockServer] = blockHashes.Hashes
		}
		return conn.Close()
	}
	return fmt.Errorf("could not find a leader")
}

func (surfClient *RPCClient) GetBlockStoreAddrs(blockStoreAddrs *[]string) error {
	for _, server := range surfClient.MetaStoreAddrs {
		conn, err := grpc.Dial(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		c := NewRaftSurfstoreClient(conn)

		// perform the call
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		a, err := c.GetBlockStoreAddrs(ctx, &emptypb.Empty{})

		if err != nil {
			if errors.Is(err, ErrNotLeader) || errors.Is(err, ErrServerCrashed) || errors.Is(err, ErrServerCrashedUnreachable) {
				conn.Close()
				continue
			}
			conn.Close()
			return err
		}

		*blockStoreAddrs = a.BlockStoreAddrs
		return conn.Close()
	}
	return fmt.Errorf("could not find a leader")
}

// This line guarantees all method for RPCClient are implemented
var _ ClientInterface = new(RPCClient)

// Create an Surfstore RPC client
func NewSurfstoreRPCClient(addrs []string, baseDir string, blockSize int) RPCClient {
	return RPCClient{
		MetaStoreAddrs: addrs,
		BaseDir:        baseDir,
		BlockSize:      blockSize,
	}
}
