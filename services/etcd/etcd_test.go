package etcd

import "testing"

func TestEtcdGet(t *testing.T) {
	Get("/foobar")
}
