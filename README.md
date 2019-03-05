# Node Maintenance Operator
node-maintenance-operator is an operator generated from the [operator-sdk](https://github.com/operator-framework/operator-sdk).
The purpose of this operator is to watch for new or deleted custom resources called NodeMaintenance which indicate that a node in the cluster should either:
 - NodeMaintenance CR created: move node into maintenance, cordon the node - set it as unschedulable and evict the pods (which can be evicted) from that node.
  - NodeMaintenance CR deleted: remove pod from maintenance and uncordon the node - set it as schedulable

> *Note*:  The current behavior of the  operator is to mimic `kubectl drain <node name>` as performed in [Kubevirt - evict all VMs and Pods on a node ](https://kubevirt.io/user-guide/docs/latest/administration/node-eviction.html#how-to-evict-all-vms-and-pods-on-a-node)

## Build and run the operator
Before running the operator, the NodeMaintenance CRD must be registered with the Openshift/Kubernetes apiserver:

```sh
$ kubectl create -f deploy/crds/kubevirt_v1alpha3_nodemaintenance_crd.yaml
```

Once this is done, there are two ways to run the operator:

- As Go program outside a cluster
- As a Deployment inside a Openshift/Kubernetes cluster

### 1. Run locally outside the cluster

This method is preferred during development cycle to deploy and test faster.

Set the name of the operator in an environment variable:

```sh
export OPERATOR_NAME=node-maintenance-operator
```

Run the operator locally with the default Kubernetes config file present at `$HOME/.kube/config` or  with specificing kubeconfig via the flag `--kubeconfig=<path/to/kubeconfig>`:

```sh
$ operator-sdk up local --kubeconfig="/home/usr/go/src/kubevirt.io/kubevirt/cluster/os-3.11.0/.kubeconfig"

INFO[0000] Running the operator locally.                
INFO[0000] Using namespace default.                     
{"level":"info","ts":1551793839.3308277,"logger":"cmd","msg":"Go Version: go1.11.4"}
{"level":"info","ts":1551793839.3308823,"logger":"cmd","msg":"Go OS/Arch: linux/amd64"}
{"level":"info","ts":1551793839.330899,"logger":"cmd","msg":"Version of operator-sdk: v0.5.0+git"}
...

```

## Setting Node Maintenance 
### Set Maintenance on - Create a NodeMaintenance CR
To set maintenance on a node a `NodeMaintenance` CR should be created.
A `NodeMaintenance` CR contains:
- Name: The name of the node which will be put into maintenance
- Reason: the reason for the node maintenance

Create the example `NodeMaintenance` CR found at `deploy/crds/kubevirt_v1alpha3_nodemaintenance_cr.yaml`:

```sh
$ cat deploy/crds/kubevirt_v1alpha3_nodemaintenance_cr.yaml

apiVersion: kubevirt.io/v1alpha3
kind: NodeMaintenance
metadata:
  name: node02
spec:
  reason: "Test node maintenance"

$ kubectl apply -f deploy/crds/cache_v1alpha1_memcached_cr.yaml
{"level":"info","ts":1551794418.6742408,"logger":"controller_nodemaintenance","msg":"Reconciling NodeMaintenance","Request.Namespace":"default","Request.Name":"node02"}
{"level":"info","ts":1551794418.674294,"logger":"controller_nodemaintenance","msg":"Applying Maintenance mode on Node: node02 with Reason: Test node maintenance","Request.Namespace":"default","Request.Name":"node02"}
{"level":"info","ts":1551783365.7430992,"logger":"controller_nodemaintenance","msg":"WARNING: ignoring DaemonSet-managed Pods: default/local-volume-provisioner-5xft8, kubevirt/disks-images-provider-bxpc5, kubevirt/virt-handler-52kpr, openshift-monitoring/node-exporter-4c9jt, openshift-node/sync-8w5x8, openshift-sdn/ovs-kvz9w, openshift-sdn/sdn-qnjdz\n"}
{"level":"info","ts":1551783365.7471824,"logger":"controller_nodemaintenance","msg":"evicting pod \"virt-operator-5559b7d86f-2wsnz\"\n"}
{"level":"info","ts":1551783365.7472217,"logger":"controller_nodemaintenance","msg":"evicting pod \"cdi-operator-55b47b74b5-9v25c\"\n"}
{"level":"info","ts":1551783365.747241,"logger":"controller_nodemaintenance","msg":"evicting pod \"virt-api-7fcd86776d-652tv\"\n"}
{"level":"info","ts":1551783365.747243,"logger":"controller_nodemaintenance","msg":"evicting pod \"simple-deployment-1-m5qv9\"\n"}
{"level":"info","ts":1551783365.7472336,"logger":"controller_nodemaintenance","msg":"evicting pod \"virt-controller-8987cffb8-29w26\"\n"}
{"level":"info","ts":1551783365.7472618,"logger":"controller_nodemaintenance","msg":"evicting pod \"virt-controller-8987cffb8-jkgwx\"\n"}
{"level":"info","ts":1551783365.747213,"logger":"controller_nodemaintenance","msg":"evicting pod \"cdi-http-import-server-875858fbd-qh425\"\n"}
{"level":"info","ts":1551783365.7472284,"logger":"controller_nodemaintenance","msg":"evicting pod \"cdi-apiserver-76f5c57975-4d68g\"\n"}
{"level":"info","ts":1551783365.7472222,"logger":"controller_nodemaintenance","msg":"evicting pod \"cdi-uploadproxy-5c8cf6d885-gqvf6\"\n"}
{"level":"info","ts":1551783365.7471826,"logger":"controller_nodemaintenance","msg":"evicting pod \"cdi-deployment-7497dbd678-8ddf7\"\n"}
{"level":"info","ts":1551783372.3010166,"logger":"controller_nodemaintenance","msg":"evicted Pod: simple-deployment-1-m5qv9"}
{"level":"info","ts":1551783373.102215,"logger":"controller_nodemaintenance","msg":"evicted Pod: virt-operator-5559b7d86f-2wsnz"}
{"level":"info","ts":1551783373.7016778,"logger":"controller_nodemaintenance","msg":"evicted Pod: virt-controller-8987cffb8-jkgwx"}
{"level":"info","ts":1551783374.301048,"logger":"controller_nodemaintenance","msg":"evicted Pod: cdi-uploadproxy-5c8cf6d885-gqvf6"}
{"level":"info","ts":1551783374.7006464,"logger":"controller_nodemaintenance","msg":"evicted Pod: cdi-apiserver-76f5c57975-4d68g"}
{"level":"info","ts":1551783375.1010072,"logger":"controller_nodemaintenance","msg":"evicted Pod: virt-controller-8987cffb8-29w26"}
{"level":"info","ts":1551783375.3010762,"logger":"controller_nodemaintenance","msg":"evicted Pod: virt-api-7fcd86776d-652tv"}
{"level":"info","ts":1551783375.700461,"logger":"controller_nodemaintenance","msg":"evicted Pod: cdi-operator-55b47b74b5-9v25c"}
{"level":"info","ts":1551783376.8579638,"logger":"controller_nodemaintenance","msg":"evicted Pod: cdi-deployment-7497dbd678-8ddf7"}
{"level":"info","ts":1551783397.8587294,"logger":"controller_nodemaintenance","msg":"evicted Pod: cdi-http-import-server-875858fbd-qh425"}


```


### Set Maintenance off - Delete the NodeMaintenance CR
To remove maintenance from a node a `NodeMaintenance` CR with the node's name  should be deleted.

```sh
$ cat deploy/crds/kubevirt_v1alpha3_nodemaintenance_cr.yaml

apiVersion: kubevirt.io/v1alpha3
kind: NodeMaintenance
metadata:
  name: node02
spec:
  reason: "Test node maintenance"

$ kubectl delete -f deploy/crds/cache_v1alpha1_memcached_cr.yaml

{"level":"info","ts":1551794725.0018933,"logger":"controller_nodemaintenance","msg":"Reconciling NodeMaintenance","Request.Namespace":"default","Request.Name":"node02"}
{"level":"info","ts":1551794725.0021605,"logger":"controller_nodemaintenance","msg":"NodeMaintenance Object: default/node02 Deleted ","Request.Namespace":"default","Request.Name":"node02"}
{"level":"info","ts":1551794725.0022023,"logger":"controller_nodemaintenance","msg":"uncordon Node: node02"}

```


## Next Steps
- Trigger VM migration 
- Handle unremoved pods and daemonsets
- Check behavior for storage pods
- Verify that disruption budget is obeyed during pod eviction
- Fencing
- Versioning
- e2e tests
- Enhance error handling
- Operator integration and packaging

