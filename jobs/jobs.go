package jobs

import (
	"encoding/json"
	"github.com/Tendrl/tendrl2/services/etcd"
	"github.com/coreos/etcd/clientv3"
	"log"
)

type Job struct {
	Id   string `json:"omitempty"`
	Type string `json:"type"`
	Data string `json:"data"`
}

func Find(id string) (*clientv3.GetResponse, error) {
	cluster_id := "cb012f6c-9cc1-4390-8d19-885dbf98dd4f"
	return etcd.Getv3("/clusters/" + cluster_id + "/jobs/" + id)
}

func All() (*clientv3.GetResponse, error) {
	cluster_id := "cb012f6c-9cc1-4390-8d19-885dbf98dd4f"
	return etcd.Getv3("/clusters/" + cluster_id + "/jobs/")
}

func Enqueue(data string) error {
	cluster_id := "cb012f6c-9cc1-4390-8d19-885dbf98dd4f"
	queue := etcd.Queue("/clusters/" + cluster_id + "/jobs/")
	job := Job{
		Type: "ansible",
		Data: data,
	}
	json_job, err := json.Marshal(job)
	if err != nil {
		log.Println(err)
	}
	queue.Enqueue(string(json_job))
	return nil
}
