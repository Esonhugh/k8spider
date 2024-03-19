package pkg

import (
	"github.com/esonhugh/k8spider/define"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
)

var _ StaticTaskRunner[net.IP] = (*SubnetTask)(nil)

type SubnetTask struct {
	Subnet *net.IPNet
	Res    []define.Record
	lock   sync.Mutex
}

func (s *SubnetTask) Generator(tasks chan net.IP) {
	for _, v := range ParseIPNetToIPs(s.Subnet) {
		tasks <- v
	}
}

func (s *SubnetTask) Solver(ip net.IP) {
	ptr := PTRRecord(ip)
	if len(ptr) <= 0 {
		return
	}

	for _, domain := range ptr {
		log.Infof("PTRrecord %v --> %v", ip, domain)
		s.AddRecord(&define.Record{Ip: ip, SvcDomain: domain})
	}
}

func (s *SubnetTask) AddRecord(r *define.Record) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Res = append(s.Res, *r)
}

func ScanSubnet(subnet *net.IPNet) []define.Record {
	tasker := &SubnetTask{
		Subnet: subnet,
	}
	RunStatic[net.IP](tasker, Threads)
	return tasker.Res
}
