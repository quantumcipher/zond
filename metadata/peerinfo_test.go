package metadata

import (
	"testing"

	"github.com/theQRL/zond/ntp"
)

func TestNewPeerInfo(t *testing.T) {
	timestamp := ntp.GetNTP().Time()
	multiAddr := "peerAddress1/peerAddress2/peerAddress3/peerAddress4/peerAddress5/peerAddress6/peerAddress7"

	peerInfo := NewPeerInfo(multiAddr, timestamp)

	if peerInfo.MultiAddr() != multiAddr {
		t.Errorf("expected multiaddress (%v), got (%v)", multiAddr, peerInfo.MultiAddr())
	}
}
