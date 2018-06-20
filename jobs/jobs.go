package jobs

import (
	"encoding/json"
	"fmt"
	"github.com/Tendrl/tendrl2/sds_sync"
	"github.com/Tendrl/tendrl2/services/etcd"
	"github.com/coreos/etcd/client"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"
)

var (
	cluster_id = "my_cluster"
)

type JobData struct {
	Inventory []string `json:"inventory"`
	Playbook  string   `json:"playbook"`
}

type Job struct {
	Id         string             `json:"id,omitempty"`
	Type       string             `json:"type"`
	Data       JobData            `json:"data"`
	JobEvents  *map[string]string `json:"job_events,omitempty"`
	ReturnCode string             `json:"rc,omitempty"`
	Status     string             `json:"status,omitempty"`
	Stdout     string             `json:"stdout,omitempty"`
	Created    string             `json:"created_at,omitempty"`
	Modified   string             `json:"modified_at,omitempty"`
}

func Find(id string) (*Job, error) {
	resp, err := etcd.Get("/clusters/" + cluster_id + "/jobs/processing/" + id + "/data")
	if err != nil {
		log.Println("Etcd call failed in finding job:", id)
		return nil, err
	}
	job := Job{}
	json.Unmarshal([]byte(resp.Node.Value), &job)
	dirname := "artifacts/" + id + "/"
	job.JobEvents = fileMap(dirname + "job_events/")

	status, _ := ioutil.ReadFile(dirname + "status")
	job.Status = string(status)

	rc, _ := ioutil.ReadFile(dirname + "rc")
	job.ReturnCode = string(rc)

	stdout, _ := ioutil.ReadFile(dirname + "stdout")
	job.Stdout = string(stdout)

	return &job, nil
}

func fileMap(path string) *map[string]string {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Println("No files in directory "+path, err)
		return nil
	}
	fileMap := map[string]string{}
	for _, file := range files {
		b, err := ioutil.ReadFile(path + file.Name())
		if err != nil {
			log.Println(err)
		}
		fileMap[file.Name()] = string(b)
	}
	return &fileMap
}

func All() (*client.Response, error) {
	return etcd.Get("/clusters/" + cluster_id + "/jobs/processing")
}

func New(data JobData) *Job {
	timestamp := time.Now().Format(time.RFC3339)
	return &Job{
		Id:       uuid.New().String(),
		Type:     "ansible",
		Data:     data,
		Created:  timestamp,
		Modified: timestamp,
	}
}

func (job Job) Enqueue() {
	queue := etcd.Queue("/clusters/" + cluster_id + "/jobs/pending")
	log.Println("enqueueing")
	json_job, err := json.Marshal(job)
	if err != nil {
		panic(err)
	}
	queue.Enqueue(string(json_job))
}

func (job Job) Json() string {
	json_job, err := json.Marshal(job)
	if err != nil {
		log.Println("Error marshalling job")
	}
	return string(json_job)
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

func (job Job) Write() error {
	json_job, _ := json.Marshal(job)
	_, err := etcd.Set("/clusters/"+cluster_id+"/jobs/processing/"+job.Id+"/data", string(json_job))
	return err
}

func (job Job) Register() error {
	job.Modified = time.Now().Format(time.RFC3339)
	return job.Write()
}

func Work() string {
	job := Dequeue()
	job.Register()
	go job.Run()
	return job.Id
}

func (job Job) Run() {
	log.Println("Running job: ", job.Data)
	inventory := "[gluster4-servers]\n" + strings.Join(job.Data.Inventory, "\n")
	ioutil.WriteFile("inventory", []byte(inventory), 0755)

	cmd := exec.Command("ansible-runner", "-p", job.Data.Playbook, "-i", job.Id, "run", ".")

	err := cmd.Start()

	if err != nil {
		log.Println(err.Error())
		log.Println("exit in start")
		return
	}

	err = cmd.Wait()
	if err != nil {
		log.Println(err.Error())
		log.Println("exit in wait")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Sync Error, is the cluster not up yet?", r)
		}
	}()
	if err := sds_sync.SyncAll("http://" + job.Data.Inventory[0] + ":24007"); err != nil {
		log.Println("Sync Error", err)
	}
}
