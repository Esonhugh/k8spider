package mutli

import (
	"net"

	"github.com/esonhugh/k8spider/define"
)

func ScanAll(subnet *net.IPNet, num int) (result <-chan define.Record) {
	subs := NewSubnetScanner(num)
	result = ScanServiceWithChan(subs.ScanSubnet(subnet))
	return result
}

func ScanNeighbor(namespace []string, subnet *net.IPNet, num int) <-chan []define.Record {
	subs := NewNeighborScanner(num)
	if len(namespace) == 1 {
		return subs.ScanSingleNeighbor(namespace[0], subnet)
	}
	return subs.ScanMultiNeighbor(namespace, subnet)
}

func ScanNeighborSvc(subnet *net.IPNet, num int) <-chan define.Record {
	subs := NewNeighborScanner(num)
	return ScanServiceWithChan(subs.ScanSvcNeighbor(subnet))
}
