# K8Spider 

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

## Example

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
root@pod:/var/www/html/tools# ./k8spider all -c 10.43.43.1/24
INFO[0000] PTRrecord 10.43.43.87 --> kube-state-metrics.lens-metrics.svc.cluster.local. 
INFO[0000] PTRrecord 10.43.43.93 --> metrics-server.kube-system.svc.cluster.local. 
INFO[0000] SRVRecord: kube-state-metrics.lens-metrics.svc.cluster.local. --> kube-state-metrics.lens-metrics.svc.cluster.local.:8080 
INFO[0000] SRVRecord: metrics-server.kube-system.svc.cluster.local. --> metrics-server.kube-system.svc.cluster.local.:443 
INFO[0000] {"Ip":"10.43.43.87","SvcDomain":"kube-state-metrics.lens-metrics.svc.cluster.local.","SrvRecords":[{"Cname":"kube-state-metrics.lens-metrics.svc.cluster.local.","Srv":[{"Target":"kube-state-metrics.lens-metrics.svc.cluster.local.","Port":8080,"Priority":0,"Weight":100}]}]} 
```



