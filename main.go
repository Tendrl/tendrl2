package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Tendrl/tendrl2/jobs"
	"github.com/Tendrl/tendrl2/sds_sync"
	"github.com/Tendrl/tendrl2/services/etcd"
	"goji.io"
	"goji.io/pat"
)

func newJob(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	err := jobs.Enqueue(string(body))

	if err != nil {
		panic(err)
	}
}

func main() {
	cluster_id := "cb012f6c-9cc1-4390-8d19-885dbf98dd4f"
	mux := goji.NewMux()
	mux.HandleFunc(pat.Post("/jobs"), newJob)
	log.Println("listening")
	go http.ListenAndServe("localhost:8000", mux)
	if err := sds_sync.SyncAll(os.Getenv("GD2_ENDPOINT")); err != nil {
		log.Println("Sync Error", err)
	}
	queue := etcd.Queue("/clusters/" + cluster_id + "/jobs")
	if err := queue.Enqueue("first job"); err != nil {
		panic(err)
	}
	log.Println("Listening for jobs...")
	for true {
		job_json, err := queue.Dequeue()
		if err != nil {
			panic(err)
		}
		log.Println("Got job:", job_json)
	}
}
