package sds_sync

import (
	"encoding/json"
	"github.com/Tendrl/tendrl2/services/etcd"
	"github.com/Tendrl/tendrl2/services/gd2"
	"testing"
)

func TestSdsSync(t *testing.T) {
	peers, _ := gd2.Client.Peers()
	cluster_id := "cb012f6c-9cc1-4390-8d19-885dbf98dd4f"
	actual, _ := etcd.Get("/clusters/" + cluster_id + "/peers/" + peers[0].ID.String() + "/data")
	syncPeers()
	actual, _ := etcd.Get("/peers")
	expected, err := json.Marshal(peers[0])
	if string(expected) != string(actual) {
		t.Errorf("Error occured while syncing GD2 Peers: %d != %d", expected, actual)
	}
}
