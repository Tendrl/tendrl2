package sds_sync

import (
	"encoding/json"
	"github.com/Tendrl/tendrl2/services/etcd"
	"github.com/Tendrl/tendrl2/services/gd2"
	"log"
	"time"
	//"github.com/gluster/glusterd2/pkg/restclient"
)

func init() {
}

func SyncAll(endpoint string) error {
	syncPeers(endpoint)
	syncCluster(endpoint)
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	return nil
}
func syncCluster(endpoint string) {
	cluster_id := "cb012f6c-9cc1-4390-8d19-885dbf98dd4f"
	version := gd2.Version(endpoint)
	log.Println(version)
	cluster := make(map[string]interface{})
	cluster["integration_id"] = cluster_id
	cluster["current_job"] = "{}"
	cluster["short_name"] = ""
	cluster["last_sync"] = time.Now().String()
	cluster["is_managed"] = "yes"
	cluster["sds_version"] = version["glusterd-version"]
	cluster["sds_name"] = "gluster"
	cluster["gd2-api-version"] = version["api-version"]
	log.Println(cluster)
	json_cluster, _ := json.Marshal(cluster)
	if _, err := etcd.Set("/data", string(json_cluster)); err != nil {
		panic(err)
	}
}

func syncPeers(endpoint string) error {
	// export GD2_ENDPOINT="http://192.168.121.222:24007"

	peers, err := gd2.Client(endpoint).Peers()
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(peers); i++ {
		json_peer, _ := json.Marshal(peers[i])
		resp, err := etcd.Set("/peers/"+peers[i].ID.String()+"/data", string(json_peer)) // TODO Add TTL 2 min
		log.Println(resp, err)
	}
	return nil
}
