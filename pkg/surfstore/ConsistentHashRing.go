package surfstore

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

// serverHash -> serverName
type ConsistentHashRing struct {
	ServerMap map[string]string
}

func (c ConsistentHashRing) GetResponsibleServer(blockId string) string {
	// 1. sort hash values (key in hash ring)
	hashes := []string{}
	for h := range c.ServerMap {
		hashes = append(hashes, h)
	}
	sort.Strings(hashes)
	// 2. find the first server with larger hash value than blockHash
	responsibleServer := ""
	for i := 0; i < len(hashes); i++ {
		if hashes[i] > blockId {
			responsibleServer = c.ServerMap[hashes[i]]
			break
		}
	}
	if responsibleServer == "" {
		responsibleServer = c.ServerMap[hashes[0]]
	}
	return responsibleServer
}

func (c ConsistentHashRing) Hash(addr string) string {
	h := sha256.New()
	h.Write([]byte(addr))
	return hex.EncodeToString(h.Sum(nil))

}

func NewConsistentHashRing(serverAddrs []string) *ConsistentHashRing {
	consistentHashRing := ConsistentHashRing{ServerMap: make(map[string]string)}
	for _, serverAddr := range serverAddrs {
		serverName := "blockstore" + serverAddr
		serverHash := consistentHashRing.Hash(serverName)
		consistentHashRing.ServerMap[serverHash] = serverAddr
	}
	return &consistentHashRing
}
