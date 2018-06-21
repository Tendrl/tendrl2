package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Tendrl/tendrl2/jobs"
	"goji.io"
	"goji.io/pat"
)

func newJob(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	job_data := jobs.JobData{}
	if err := json.Unmarshal([]byte(body), &job_data); err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "{\"error\":\"", err, "\"}")
	}
	job := jobs.New(job_data)
	job.Enqueue()
	fmt.Fprintf(w, job.Json())
}

func getJob(w http.ResponseWriter, r *http.Request) {
	id := pat.Param(r, "job_id")
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(w, "{\"error\":\"", r, "\"}")
		}
	}()
	job, err := jobs.Find(id)
	if err != nil {
		fmt.Fprintf(w, "{\"error\":\""+err.Error()+"\"}")
		return
	}
	fmt.Fprintf(w, job.Json())
}

func main() {
	cluster_id := "my_cluster"
	mux := goji.NewMux()
	mux.HandleFunc(pat.Post("/clusters/"+cluster_id+"/jobs"), newJob)
	mux.HandleFunc(pat.Get("/clusters/"+cluster_id+"/jobs/:job_id"), getJob)
	log.Println("Listening for jobs...")
	go http.ListenAndServe("0.0.0.0:8000", mux)
	log.Println("Starting job worker...")
	for true {
		job_id := jobs.Work()
		log.Println("Working on job:", job_id)
	}
}
