package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Tendrl/tendrl2/jobs"
	"github.com/Tendrl/tendrl2/sds_sync"
	"goji.io"
	"goji.io/pat"
)

func newJob(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	job_id := jobs.Enqueue(string(body))
	fmt.Fprintf(w, "{\"job_id\":\""+job_id+"\"}")
}

func getJobs(w http.ResponseWriter, r *http.Request) {
	dirname := "/var/log/ansible/hosts/"
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
	cluster_id := "cb012f6c-9cc1-4390-8d19-885dbf98dd4f"
	mux := goji.NewMux()
	mux.HandleFunc(pat.Post("/clusters/"+cluster_id+"/jobs"), newJob)
	mux.HandleFunc(pat.Get("/clusters/"+cluster_id+"/jobs"), getJob)
	log.Println("listening")
	go http.ListenAndServe("localhost:8000", mux)
	if err := sds_sync.SyncAll(os.Getenv("GD2_ENDPOINT")); err != nil {
		log.Println("Sync Error", err)
	}
	log.Println("Listening for jobs...")
	for true {
		job_id := jobs.Work()
		log.Println("Working on job:", job_id)
	}
}
