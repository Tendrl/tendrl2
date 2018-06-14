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
	mux := goji.NewMux()
	mux.HandleFunc(pat.Post("/jobs"), newJob)
	log.Println("listening")
	go http.ListenAndServe("localhost:8000", mux)
	if err := sds_sync.SyncAll(os.Getenv("GD2_ENDPOINT")); err != nil {
		log.Println("Sync Error", err)
	}
	queue := etcd.Queue("/jobs/pending")
	if err := queue.Enqueue("first job"); err != nil {
		panic(err)
	}
	for true {
		log.Println("dequeueing")
		job_json, err := queue.Dequeue()
		if err != nil {
			panic(err)
		}
		log.Println(job_json)
	}
}
