# tendrl2
Tendrl 2 (Gluster 2, Ceph 3 monitoring and management service)


## Whats new in Tendrl 2
- Targeted to work seamlessly with various Gluster deployments methods like standalone VMs, OpenShift/Kubernetes.
- Tendrl http API to focus on enabling easy orchestration and reporting of Gluster workflows (install, expand, CRUD for bricks, volumes etc) via ansible playbooks (eg: gluster-ansible).
- Tendrl 2 will use prometheus.io for collecting time series data from storage nodes and the Gluster cluster.
- Tendrl 2 is mainly written in Golang to enable re-use of code/patterns from Gluster 4 and OpenShift components which are written in Golang.
- Tendrl 2 will provide feature/metrics parity with Tendrl 1
