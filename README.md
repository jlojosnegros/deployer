# topology-aware-scheduling deployer

`deployer` is a set of go packages and a command line tool to setup all the components and settings needed to enable
the topology-aware-scheduling on a kubernetes cluster. Additionally, the tool can validate if the cluster configuration
is compatible to the topology-aware-scheduling requirements.

## requirements

* kubernetes >= 1.21
* a valid `kubeconfig`
* **validation only** `kubectl` >= 1.21 in your `PATH`

## how does it work?

### deploy on a kubernetes cluster

Considering a kind cluster configured like this:
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
kubeadmConfigPatches:
- |
  kind: KubeletConfiguration
  cpuManagerPolicy: "static"
  topologyManagerPolicy: "single-numa-node"
  reservedSystemCPUs: "0,16"
  featureGates:
    KubeletPodResourcesGetAllocatable: true
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
```

Deploy all the topology-aware scheduling components:
```
$ kubectl get nodes
NAME                 STATUS   ROLES                  AGE     VERSION
kind-control-plane   Ready    control-plane,master   3m48s   v1.21.1
kind-worker          Ready    worker                 3m15s   v1.21.1
kind-worker2         Ready    worker                 3m15s   v1.21.1
kind-worker3         Ready    worker                 3m15s   v1.21.1

$ ./deployer deploy -W
2021/07/20 06:16:34 deploying topology-aware-scheduling API...
2021/07/20 06:16:34 -  API> created CustomResourceDefinition "noderesourcetopologies.topology.node.k8s.io"
2021/07/20 06:16:34 ...deployed topology-aware-scheduling API!
2021/07/20 06:16:34 deploying topology-aware-scheduling topology updater...
2021/07/20 06:16:35 -  RTE> created Namespace "tas-topology-updater"
2021/07/20 06:16:35 -  RTE> created ServiceAccount "rte"
2021/07/20 06:16:35 -  RTE> created ClusterRole "rte"
2021/07/20 06:16:35 -  RTE> created ClusterRoleBinding "rte"
2021/07/20 06:16:35 -  RTE> created DaemonSet "resource-topology-exporter-ds"
2021/07/20 06:16:35 wait for all the pods in deployment tas-topology-updater resource-topology-exporter-ds to be running and ready
2021/07/20 06:16:35 no pods found for tas-topology-updater resource-topology-exporter-ds
2021/07/20 06:16:36 pod tas-topology-updater resource-topology-exporter-ds-xb88c not ready yet (Pending)
2021/07/20 06:16:37 all the pods in daemonset tas-topology-updater resource-topology-exporter-ds are running and ready!
2021/07/20 06:16:37 ...deployed topology-aware-scheduling topology updater!
2021/07/20 06:16:37 deploying topology-aware-scheduling scheduler plugin...
2021/07/20 06:16:38 -  SCD> created ServiceAccount "topo-aware-scheduler"
2021/07/20 06:16:38 -  SCD> created ClusterRole "noderesourcetoplogy-handler"
2021/07/20 06:16:38 -  SCD> created ClusterRoleBinding "topo-aware-scheduler-as-kube-scheduler"
2021/07/20 06:16:38 -  SCD> created ClusterRoleBinding "noderesourcetoplogy"
2021/07/20 06:16:38 -  SCD> created ClusterRoleBinding "topo-aware-scheduler-as-volume-scheduler"
2021/07/20 06:16:38 -  SCD> created RoleBinding "topo-aware-scheduler-as-kube-scheduler"
2021/07/20 06:16:38 -  SCD> created ConfigMap "topo-aware-scheduler-config"
2021/07/20 06:16:38 -  SCD> created Deployment "topo-aware-scheduler"
2021/07/20 06:16:38 wait for all the pods in deployment kube-system topo-aware-scheduler to be running and ready
2021/07/20 06:16:38 no pods found for kube-system topo-aware-scheduler
2021/07/20 06:16:48 all the pods in deployment kube-system topo-aware-scheduler are running and ready!
2021/07/20 06:16:48 ...deployed topology-aware-scheduling scheduler plugin!
```

#### cleaning up (removing):

```
$ ./deployer remove -W
2021/07/20 06:18:12 removing topology-aware-scheduling scheduler plugin...
2021/07/20 06:18:13 -  SCD> deleted Deployment "topo-aware-scheduler"
2021/07/20 06:18:13 wait for all the pods in deployment kube-system topo-aware-scheduler to be gone
2021/07/20 06:18:13 error removing: still 1 pods found for kube-system topo-aware-scheduler
2021/07/20 06:18:13 removing topology-aware-scheduling topology updater...
2021/07/20 06:18:13 -  RTE> deleted Namespace "tas-topology-updater"
2021/07/20 06:18:13 wait for the deployment namespace "tas-topology-updater" to be gone
2021/07/20 06:18:40 namespace gone for daemonset "tas-topology-updater"!
2021/07/20 06:18:40 -  RTE> deleted ClusterRoleBinding "rte"
2021/07/20 06:18:40 -  RTE> deleted ClusterRole "rte"
2021/07/20 06:18:40 ...removed topology-aware-scheduling topology updater!
2021/07/20 06:18:40 removing topology-aware-scheduling API...
2021/07/20 06:18:41 -  API> deleted CustomResourceDefinition "noderesourcetopologies.topology.node.k8s.io"
2021/07/20 06:18:41 ...removed topology-aware-scheduling API!
```

### validate the cluster configuration:

A kind cluster with the correct configuration:
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
kubeadmConfigPatches:
- |
  kind: KubeletConfiguration
  cpuManagerPolicy: "static"
  topologyManagerPolicy: "single-numa-node"
  reservedSystemCPUs: "0,16"
  featureGates:
    KubeletPodResourcesGetAllocatable: true
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
```

Does pass the validation:
```
$ ./deployer validate
PASSED>>: the cluster configuration looks ok!
```

A kind cluster configured like this
```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
kubeadmConfigPatches:
- |
  kind: KubeletConfiguration
  cpuManagerPolicy: "static"
  reservedSystemCPUs: "0,16"
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
```

Does **not** pass the validation:
```
$ ./deployer validate
ERROR#000: Incorrect configuration of node "kind-worker" area "kubelet" component "feature gates" setting "": expected "present" detected "missing data"
ERROR#001: Incorrect configuration of node "kind-worker" area "kubelet" component "topology manager" setting "policy": expected "single-numa-node" detected "none"
ERROR#002: Incorrect configuration of node "kind-worker2" area "kubelet" component "feature gates" setting "": expected "present" detected "missing data"
ERROR#003: Incorrect configuration of node "kind-worker2" area "kubelet" component "topology manager" setting "policy": expected "single-numa-node" detected "none"
ERROR#004: Incorrect configuration of node "kind-worker3" area "kubelet" component "feature gates" setting "": expected "present" detected "missing data"
ERROR#005: Incorrect configuration of node "kind-worker3" area "kubelet" component "topology manager" setting "policy": expected "single-numa-node" detected "none"
```

## license
(C) 2021 Red Hat Inc and licensed under the Apache License v2

## build
just run
```bash
make
```

## releases

TBD
