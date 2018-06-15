package jobs

import (
	"encoding/json"
	"github.com/Tendrl/tendrl2/services/etcd"
	"github.com/coreos/etcd/client"
	"github.com/google/uuid"
	"log"
	//"os/exec"
)

var (
	cluster_id = "my_cluster"
)

type Job struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type"`
	Data string `json:"data"`
}

func Find(id string) (*client.Response, error) {
	return etcd.Get("/clusters/" + cluster_id + "/jobs/processing/" + id + "/data")
}

func All() (*client.Response, error) {
	return etcd.Get("/clusters/" + cluster_id + "/jobs/processing")
}

func Enqueue(data string) string {
	queue := etcd.Queue("/clusters/" + cluster_id + "/jobs/pending")
	log.Println("enqueueing")
	job := Job{
		Id:   uuid.New().String(),
		Type: "ansible",
		Data: data,
	}
	json_job, err := json.Marshal(job)
	if err != nil {
		panic(err)
	}
	queue.Enqueue(string(json_job))
	return job.Id
}

func Dequeue() Job {
	queue := etcd.Queue("/clusters/" + cluster_id + "/jobs/pending")
	job_json, err := queue.Dequeue()
	if err != nil {
		panic(err)
	}
	log.Println("Got job:", job_json)
	job := Job{}
	if err := json.Unmarshal([]byte(job_json), &job); err != nil {
		panic(err)
	}
	return job
}

func (job Job) Register() error {
	json_job, _ := json.Marshal(job)
	_, err := etcd.Set("/clusters/"+cluster_id+"/jobs/processing/"+job.Id+"/data", string(json_job))
	return err
}

func Work() string {
	job := Dequeue()
	job.Register()
	go job.Run()
	return job.Id
}

func (job Job) Run() {
	//inventory := []string{"192.168.121.139", "192.168.121.222"}
	//arguments := "provision --provision-with prepare_gluster.yml"
	//exec.Command("ansible-playbook", append([]string{"-i"}, inventory, []string{"ansible/prepare_gluster.yml"}))
	log.Println("ansible-playbook", append([]string{"-i", "192.168.121.139", "192.168.121.222", "ansible/prepare_gluster.yml"}))
}
