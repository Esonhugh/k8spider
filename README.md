# K8Spider 

<img src="./K8spider.webp" width="200px">

> work like a spider inside your Kubernetes and hunting other service.

K8Spider is a simple tools for Kubernetes Service Discovery. 

It inspired from k8slanparty.com. That dnscan subnet is useful in challenges.

And I extended it ability on Kubernetes Service Discovery.

Now it supports to scan all services installed in Kubernetes cluster and all exposed ports in service. 

## Build

```bash
make 
```

## Download 

Checkout the release page. 

## Usage

```bash
# in kubernetes pods
echo $KUBERNETES_SERVICE_HOST
# if KUBERNETES_SERVICE_HOST is empty, you can use the following command to set it.
# export KUBERNETES_SERVICE_HOST=x.x.x.x
# or ./k8spider -c x.x.x.x/16 all
./k8spider all
```

Use in the kubernetes

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: spider
  name: spider
spec:
  containers:
  - image: k8spider/k8spider
    name: spider
    resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Always
# kubectl apply -f spider.yaml
```

or just using kubectl run

```bash
## just run it! 
kubectl run spider --image k8spider/k8spider 
```

and watch result with

```bash
kubectl logs spider
```

## Example

### Normal Attack - all command - ALL IN ONE

```bash
root@pod:/var/www/html/tools# env |grep KUBERNETES
KUBERNETES_SERVICE_PORT_HTTPS=443
KUBERNETES_SERVICE_PORT=443
KUBERNETES_PORT_443_TCP=tcp://10.43.0.1:443
KUBERNETES_PORT_443_TCP_PROTO=tcp
KUBERNETES_PORT_443_TCP_ADDR=10.43.0.1
KUBERNETES_SERVICE_HOST=10.43.0.1
KUBERNETES_PORT=tcp://10.43.0.1:443
KUBERNETES_PORT_443_TCP_PORT=443

root@pod:/var/www/html/tools# ./k8spider all # or  try ./k8spider all -c 10.43.0.1/16  
INFO[0000] PTRrecord 10.43.43.87 --> kube-state-metrics.lens-metrics.svc.cluster.local. 
INFO[0000] PTRrecord 10.43.43.93 --> metrics-server.kube-system.svc.cluster.local. 
INFO[0000] SRVRecord: kube-state-metrics.lens-metrics.svc.cluster.local. --> kube-state-metrics.lens-metrics.svc.cluster.local.:8080 
INFO[0000] SRVRecord: metrics-server.kube-system.svc.cluster.local. --> metrics-server.kube-system.svc.cluster.local.:443 
INFO[0000] {"Ip":"10.43.43.87","SvcDomain":"kube-state-metrics.lens-metrics.svc.cluster.local.","SrvRecords":[{"Cname":"kube-state-metrics.lens-metrics.svc.cluster.local.","Srv":[{"Target":"kube-state-metrics.lens-metrics.svc.cluster.local.","Port":8080,"Priority":0,"Weight":100}]}]} 
```

This command will try wildcard (any.any.svc.cluster.local) / Axfr dumping at first and brute force all services in the cluster.

#### Advanced 1: threading mode

```bash
./k8spider all -t  
# if you want to higher threads, you can use 
./k8spider all -t -n 16
```

#### Advanced 2: no default Zone (cluster.local) and specific DNS server

```bash
./k8spider all -z myzone.com -d 10.43.0.10:53
```

> remember if kubernetes DNS is reachable at remote, you can use it to scan all services under the cluster COMPLETELY REMOTELY.
> 

### Normal Attack - wildcard and axfr command

```bash
./k8spider axfr 
./k8spider axfr -z myzone.com -d 10.10.0.10:53
./k8spider wild
```

### Advanced Conditional Attack - neighbor command

```bash
./k8spider neighbor -p <pod-cidr check your ifconfig eth0> -n <current-ns>
```

If your kubernetes dns sets verified pod mode, it will give your pod ip a DNS name under this namespace, and non allocated
IP never have.

But it's non-default option for dns settings. 

Default is insecure pod, and it will respond your any (include invalid/non-exists) pod DNS with given IP.

### Customized Attack - service 

```bash
./k8spider srv -s kubernetes.default 
```

This command will respond you with registered service ports.

### Customized Attack - subnet

```bash
./k8spider subnet <-c cidr-srv> 
```

This command will only scan PTR service in the given subnet.

### helpers - whereisdns 

This command will help you to find out where is the kubernetes DNS server. It uses some specific DNS query to find it in given 
cidr

### helpers - metrics

This command will help you to parse the kube-state-metrics information and extract all useful information in metrics.

like 

```text
# HELP kube_service_info [STABLE] Information about service.
# TYPE kube_service_info gauge
kube_service_info{namespace="default",service="fastgpt-sandbox-service",uid="61b0674c-33c3-4e6d-a7a1-51157491a35a",cluster_ip="10.43.81.90",external_name="",load_balancer_ip=""} 1
```

to 

```text
{"namespace":"default","type":"service","name":"fastgpt-sandbox-service","spec":{"cluster_ip":["10.43.81.90"]}}
```

