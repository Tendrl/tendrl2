package sds_sync

import (
	"encoding/json"
	"github.com/Tendrl/tendrl2/services/etcd"
	"github.com/Tendrl/tendrl2/services/gd2"
	"testing"
)

func TestSdsSync(t *testing.T) {
	peers, _ := gd2.Client.Peers()
	SyncPeers()
	actual, _ := etcd.Get("/peers")
	expected, err := json.Marshal(peers)
	if string(expected) != string(actual) {
		t.Errorf("Error occured while syncing GD2 Peers: %d != %d", expected, actual)
	}
}
