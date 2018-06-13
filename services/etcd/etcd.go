package etcd

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/client"
	"log"
	"time"
)

var (
	Client client.KeysAPI
)

func init() {
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		//set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	Client = client.NewKeysAPI(c)
}

func Get(path string) (*client.Response, error) {
	namespaced_path := "/tendrl2" + path
	resp, err := Client.Get(context.Background(), namespaced_path, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		fmt.Sprintf("Get is done. Metadata is %q\n", resp)
		// print value
		fmt.Sprintf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	}
	return resp, err
}

func Set(path string, value string) (*client.Response, error) {
	namespaced_path := "/tendrl2" + path
	resp, err := Client.Set(context.Background(), namespaced_path, value, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		// print common key info
		fmt.Sprintf("Get is done. Metadata is %q\n", resp)
		// print value
		fmt.Sprintf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	}
	return resp, err
}
