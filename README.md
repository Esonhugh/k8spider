# K8Spider 

K8Spider is a simple tools for Kubernetes Service Discovery. 

It inspired from k8slanparty.com. That dnscan subnet is useful in challenges.

And I extended it ability on Kubernetes Service Discovery.

Now it supports to scan all services installed in Kubernetes cluster and all exposed ports in service.

## build

```bash
make 
```

## Usage

```bash
# in kubernetes pods
echo $KUBERNETES_SERVICE_HOST
# if KUBERNETES_SERVICE_HOST is empty, you can use the following command to set it.
# export KUBERNETES_SERVICE_HOST=x.x.x.x
# or ./k8spider -c x.x.x.x/16 all
./k8spider all
```



