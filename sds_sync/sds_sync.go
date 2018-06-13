package sds_sync

import (
	"fmt"
	//"github.com/ddliu/go-httpclient"
	"encoding/json"
	"github.com/Tendrl/tendrl2/services/etcd"
	"github.com/Tendrl/tendrl2/services/gd2"
	//"github.com/gluster/glusterd2/pkg/restclient"
)

func init() {
}

//type GetPeerList []Peers

func SyncPeers() string {
	//httpclient.Defaults(httpclient.Map{
	//httpclient.OPT_USERAGENT: "tendrl2-go",
	//"Accept-Language":        "en-us",
	//})

	//res, err := httpclient.Get("http://192.168.121.222/v1/peers", map[string]string{
	//"q": "news",
	//})

	// export GD2_ENDPOINT="http://192.168.121.222:24007"
	peers, _ := gd2.Client.Peers()
	for i := 0; i < len(peers); i++ {
		json_peer, _ := json.Marshal(peers[i])
		resp, err := etcd.Set("/peers/"+peers[i].ID.String()+"/data", string(json_peer)) // TODO Add TTL 1 min
		fmt.Print(fmt.Sprintf("%+v\n", resp), err)

	}
	return "Hello, world."
}

func BuildHi() string {
	return "Hi, world."
}
