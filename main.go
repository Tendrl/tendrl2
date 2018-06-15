package main

import (
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

func main() {
	cluster_id := "cb012f6c-9cc1-4390-8d19-885dbf98dd4f"
	mux := goji.NewMux()
	mux.HandleFunc(pat.Post("/clusters/"+cluster_id+"/jobs"), newJob)
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
