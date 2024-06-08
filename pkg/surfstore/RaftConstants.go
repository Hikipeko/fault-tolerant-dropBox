package surfstore

import (
	"fmt"
)

type PendingRequest struct {
	success bool
	err     error
}

type PeerUpdateStatus int64

const (
	PeerUnknown     PeerUpdateStatus = 0
	PeerUpdating    PeerUpdateStatus = 1
	PeerUpdated     PeerUpdateStatus = 2
	PeerUnreachable PeerUpdateStatus = 3
)

type PeerStatusResponse struct {
	id     int64
	status PeerUpdateStatus
}

var ErrServerCrashedUnreachable = fmt.Errorf("server is crashed or unreachable")
var ErrServerCrashed = fmt.Errorf("server is crashed")
var ErrNotLeader = fmt.Errorf("server is not the leader")
