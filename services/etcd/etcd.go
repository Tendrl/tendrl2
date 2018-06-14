package etcd

import (
	"context"
	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/namespace"
	"github.com/coreos/etcd/contrib/recipes"
	"log"
	"time"
)

var (
	Clientv2 client.KeysAPI
	Clientv3 *clientv3.Client
)

func init() {
	cfg := client.Config{
		Endpoints: []string{"http://localhost:2379"},
		Transport: client.DefaultTransport,
		//set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	Clientv2 = client.NewKeysAPI(c)
	Clientv3, err = clientv3.New(clientv3.Config{
		Endpoints: []string{"http://localhost:2379"},
		//set timeout per request to fail fast when the target endpoint is unavailable
		DialTimeout: time.Second,
	})
	Clientv3.KV = namespace.NewKV(Clientv3.KV, "tendrl2/")
	Clientv3.Watcher = namespace.NewWatcher(Clientv3.Watcher, "tendrl2/")
	Clientv3.Lease = namespace.NewLease(Clientv3.Lease, "tendrl2/")
	//Clientv3.Get(
}

func Get(path string) (*client.Response, error) {
	namespaced_path := "/tendrl2" + path
	resp, err := Clientv2.Get(context.Background(), namespaced_path, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Println("Get is done. Metadata is %q\n", resp)
		// print value
		log.Println("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	}
	return resp, err
}

func Set(path string, value string) (*client.Response, error) {
	namespaced_path := "/tendrl2" + path
	resp, err := Clientv2.Set(context.Background(), namespaced_path, value, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		log.Println("Get is done. Metadata is %q\n", resp)
		// print value
		log.Println("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	}
	return resp, err
}

func Queue(path string) *recipe.Queue {
	return recipe.NewQueue(Clientv3, path)
}

func Getv3(path string) (*clientv3.GetResponse, error) {
	resp, err := Clientv3.Get(context.Background(), path)
	if err != nil {
		log.Fatal(err)
	}
	return resp, err
}
