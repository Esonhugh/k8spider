package defaultResource

import "embed"

//go:embed *.yaml
var fs embed.FS

var Resources struct {
	Pod        []byte
	Deployment []byte
	Ingress    []byte
	TLSIngress []byte
}

func init() {
	Resources.Pod, _ = fs.ReadFile("pod.yaml")
	Resources.Ingress, _ = fs.ReadFile("ingress.yaml")
	Resources.TLSIngress, _ = fs.ReadFile("tls-ingress.yaml")
	Resources.Deployment, _ = fs.ReadFile("deploy.yaml")
}
