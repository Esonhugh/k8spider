package mutli

import (
	"net"
	"sync"

	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	log "github.com/sirupsen/logrus"
)

type SubnetScanner struct {
	wg    *sync.WaitGroup
	count int
}

func NewSubnetScanner(threading int) *SubnetScanner {
	return &SubnetScanner{
		wg:    new(sync.WaitGroup),
		count: threading,
	}
}

func (s *SubnetScanner) ScanSubnet(subnet *net.IPNet) <-chan define.Record {
	if subnet == nil {
		log.Debugf("subnet is nil")
		return nil
	}
	out := make(chan define.Record, 100)
	go func() {
		// if subnets, err := pkg.SubnetShift(subnet, 4); err != nil {
		if subnets, err := pkg.SubnetInto(subnet, s.count); err != nil {
			log.Errorf("Subnet split into %v failed, fallback to single mode, reason: %v", s.count, err)
			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				s.scan(subnet, out)
			}()
		} else {
			log.Debugf("Subnet split into %v success", len(subnets))
			s.wg.Add(len(subnets))
			for _, sn := range subnets {
				go func(sn *net.IPNet) {
					defer s.wg.Done()
					s.scan(sn, out)
				}(sn)
			}
		}
		s.wg.Wait()
		close(out)
	}()
	return out
}

func (s *SubnetScanner) scan(subnet *net.IPNet, to chan define.Record) {
	for _, ip := range pkg.ParseIPNetToIPs(subnet) {
		ptr := pkg.PTRRecord(ip)
		if len(ptr) > 0 {
			for _, domain := range ptr {
				log.Infof("PTRrecord %v --> %v", subnet, domain)
				r := define.Record{Ip: ip, SvcDomain: domain}
				to <- r
			}
		}
	}
	return
}
