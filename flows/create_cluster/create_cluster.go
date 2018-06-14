package create_cluster

import (
	"log"
	"os/exec",
    strings
)

func CreateCluster(hosts []string) {
    if len(hosts) > 1 {
    	log.Printf("CreateCluster.start")
        cmd := exec.Command("ansible-playbook -i %s prepare-gluster.yml", strings.Join(hosts, ", "))
        err := cmd.Run()
        if err {
            log.Printf("[Error] CreateCluster.failed")
            return

        }
        log.Printf("CreateCluster.end")
        log.Printf("CreateCluster.success")

    }
    else{
        log.Printf("CreateCluster.empty_hosts_list")

    }
}
