package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Tendrl/tendrl2/jobs"
	"goji.io"
	"goji.io/pat"
)

func newJob(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	job_id := jobs.Enqueue(string(body))
	fmt.Fprintf(w, "{\"job_id\":\""+job_id+"\"}")
}

func getJob(w http.ResponseWriter, r *http.Request) {
	dirname := "job_runs/" + pat.Param(r, "job_id") + "/"
	f, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}
	job_data := map[string]string{}
	for _, file := range files {
		b, err := ioutil.ReadFile(dirname + file.Name())
		if err != nil {
			log.Println(err)
		}
		job_data[file.Name()] = string(b)

	}
	jsonString, _ := json.Marshal(job_data)
	fmt.Fprintf(w, string(jsonString))
}

func main() {
	cluster_id := "my_cluster"
	mux := goji.NewMux()
	mux.HandleFunc(pat.Post("/clusters/"+cluster_id+"/jobs"), newJob)
	mux.HandleFunc(pat.Get("/clusters/"+cluster_id+"/jobs/:job_id"), getJob)
	log.Println("Listening for jobs...")
	go http.ListenAndServe("localhost:8000", mux)
	log.Println("Starting job worker...")
	for true {
		job_id := jobs.Work()
		log.Println("Working on job:", job_id)
	}
}
