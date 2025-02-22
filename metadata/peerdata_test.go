package metadata

import (
	"os"
	"testing"

	"github.com/theQRL/zond/config"
	"github.com/theQRL/zond/protos"
)

func CreateDirectoryIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func TestLoad(t *testing.T) {
	userConfig := config.GetUserConfig()

	err := CreateDirectoryIfNotExists(userConfig.DataDir())
	if err != nil {
		t.Error("Error creating data directory ", err.Error())
	}
	defer os.Remove(userConfig.BaseDir)
	peerData := &PeerData{
		pbData: &protos.PeerData{
			ConnectedPeers:    make([]*protos.PeerInfo, 0),
			DisconnectedPeers: make([]*protos.PeerInfo, 0),
		},
		connectedPeers:    make([]*PeerInfo, 0),
		disconnectedPeers: make([]*PeerInfo, 0),
		bootstrapNodes:    make(map[string]bool),
	}

	err = peerData.Load()

	if err != nil {
		t.Error("got unexpected error while loading peerdata", err)
	}
}

func TestSave(t *testing.T) {
	userConfig := config.GetUserConfig()

	err := CreateDirectoryIfNotExists(userConfig.DataDir())
	if err != nil {
		t.Error("Error creating data directory ", err.Error())
	}
	defer os.Remove(userConfig.BaseDir)
	peerData := &PeerData{
		pbData: &protos.PeerData{
			ConnectedPeers:    make([]*protos.PeerInfo, 0),
			DisconnectedPeers: make([]*protos.PeerInfo, 0),
		},
		connectedPeers:    make([]*PeerInfo, 0),
		disconnectedPeers: make([]*PeerInfo, 0),
		bootstrapNodes:    make(map[string]bool),
	}

	err = peerData.Save()
	if err != nil {
		t.Error("got unexpected error while saving peerdata")
	}
}

func TestAddConnectedPeers(t *testing.T) {
	userConfig := config.GetUserConfig()

	err := CreateDirectoryIfNotExists(userConfig.DataDir())
	if err != nil {
		t.Error("Error creating data directory ", err.Error())
	}
	defer os.Remove(userConfig.BaseDir)
	peerData := &PeerData{
		pbData: &protos.PeerData{
			ConnectedPeers:    make([]*protos.PeerInfo, 0),
			DisconnectedPeers: make([]*protos.PeerInfo, 0),
		},
		connectedPeers:    make([]*PeerInfo, 0),
		disconnectedPeers: make([]*PeerInfo, 0),
		bootstrapNodes:    make(map[string]bool),
	}

	peerData.AddConnectedPeers("peerAddress1/peerAddress2/peerAddress3/peerAddress4/peerAddress5/peerAddress6/peerAddress7")

	connectedPeers := peerData.ConnectedPeers()
	if len(connectedPeers) == 0 {
		t.Error("no connected peers found")
	}

	if connectedPeers[0].MultiAddr() != "peerAddress1/peerAddress2/peerAddress3/peerAddress4/peerAddress5/peerAddress6/peerAddress7" {
		t.Error("not able to add peer to connected peers list")
	}
}

func TestAddDisconnectedPeers(t *testing.T) {
	userConfig := config.GetUserConfig()

	err := CreateDirectoryIfNotExists(userConfig.DataDir())
	if err != nil {
		t.Error("Error creating data directory ", err.Error())
	}
	defer os.Remove(userConfig.BaseDir)
	peerData := &PeerData{
		pbData: &protos.PeerData{
			ConnectedPeers:    make([]*protos.PeerInfo, 0),
			DisconnectedPeers: make([]*protos.PeerInfo, 0),
		},
		connectedPeers:    make([]*PeerInfo, 0),
		disconnectedPeers: make([]*PeerInfo, 0),
		bootstrapNodes:    make(map[string]bool),
	}

	peerData.AddDisconnectedPeers("peerAddress1/peerAddress2/peerAddress3/peerAddress4/peerAddress5/peerAddress6/peerAddress7")

	disconnectedPeers := peerData.DisconnectedPeers()
	if len(disconnectedPeers) == 0 {
		t.Error("no connected peers found")
	}

	if disconnectedPeers[0].MultiAddr() != "peerAddress1/peerAddress2/peerAddress3/peerAddress4/peerAddress5/peerAddress6/peerAddress7" {
		t.Error("not able to add peer to connected peers list")
	}
}

func TestRemovePeer(t *testing.T) {
	userConfig := config.GetUserConfig()

	err := CreateDirectoryIfNotExists(userConfig.DataDir())
	if err != nil {
		t.Error("Error creating data directory ", err.Error())
	}
	defer os.Remove(userConfig.BaseDir)
	peerData := &PeerData{
		pbData: &protos.PeerData{
			ConnectedPeers:    make([]*protos.PeerInfo, 0),
			DisconnectedPeers: make([]*protos.PeerInfo, 0),
		},
		connectedPeers:    make([]*PeerInfo, 0),
		disconnectedPeers: make([]*PeerInfo, 0),
		bootstrapNodes:    make(map[string]bool),
	}

	peerData.AddConnectedPeers("peerAddress1/peerAddress2/peerAddress3/peerAddress4/peerAddress5/peerAddress6/peerAddress7")

	err = peerData.RemovePeer("peerAddress1/peerAddress2/peerAddress3/peerAddress4/peerAddress5/peerAddress6/peerAddress7")
	if err != nil {
		t.Error("unexpected error while removing peer", err)
	}

	isPeerInList := peerData.IsPeerInList("peerAddress1/peerAddress2/peerAddress3/peerAddress4/peerAddress5/peerAddress6/peerAddress7")
	if isPeerInList {
		t.Error("unable to remove peer from peer list")
	}
}
