package webrtc

import (
	"sync"
)

type Peers struct {
	ListLock    sync.RWMutex
	Connections []PeerConnectionState
	TrackLocals map[string]*webrtc.TrakLocalStaticRTP
}

func (p *Peers) DispatchKeyFrame() {

}
