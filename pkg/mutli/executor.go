package mutli

import (
	"net"

	"github.com/esonhugh/k8spider/define"
)

func ScanAll(subnet *net.IPNet) (result <-chan []define.Record) {
	subs := NewSubnetScanner()
	result = ScanServiceWithChan(subs.ScanSubnet(subnet))
	return result
}
