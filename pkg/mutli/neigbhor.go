package mutli

import (
	"fmt"
	"net"
	"sync"

	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/esonhugh/k8spider/pkg/post"
	"github.com/esonhugh/k8spider/pkg/scanner"
	log "github.com/sirupsen/logrus"
)

type NeighborScanner struct {
	wg    *sync.WaitGroup
	count int
}

func NewNeighborScanner(threading int) *NeighborScanner {
	return &NeighborScanner{
		wg:    new(sync.WaitGroup),
		count: threading,
	}
}

func (s *NeighborScanner) ScanSingleNeighbor(ns string, subnet *net.IPNet) <-chan []define.Record {
	if subnet == nil {
		log.Debugf("subnet is nil")
		return nil
	}
	out := make(chan []define.Record, 100)
	go func() {
		// if subnets, err := pkg.SubnetShift(subnet, 4); err != nil {
		if subnets, err := pkg.SubnetInto(subnet, s.count); err != nil {
			log.Errorf("Subnet split into %v failed, fallback to single mode, reason: %v", s.count, err)
			s.wg.Add(1)
			go func() {
				s.scan(ns, subnet, out)
				defer s.wg.Done()
			}()
		} else {
			log.Debugf("Subnet split into %v success", len(subnets))
			s.wg.Add(len(subnets))
			for _, sn := range subnets {
				go func(sn *net.IPNet) {
					s.scan(ns, sn, out)
					defer s.wg.Done()
				}(sn)
			}
		}
		s.wg.Wait()
		close(out)
	}()
	return out
}

func (s *NeighborScanner) ScanMultiNeighbor(nss []string, subnet *net.IPNet) <-chan []define.Record {
	out := make(chan []define.Record, 100)
	go func() {
		s.wg.Add(len(nss))
		for _, ns := range nss {
			go func(ns string) {
				defer s.wg.Done()
				s.scan(ns, subnet, out)
			}(ns)
		}
		s.wg.Wait()
		close(out)
	}()
	return out
}

func (s *NeighborScanner) scan(ns string, subnet *net.IPNet, to chan []define.Record) {
	// to <- scanner.ScanSubnet(subnet)
	for _, ip := range pkg.ParseIPNetToIPs(subnet) {
		if scanner.ScanPodExist(ip, ns) {
			newRecord := define.Record{
				Ip:    ip,
				Extra: fmt.Sprintf("%v. 0 IN A %v", pkg.IPtoPodHostName(ip.String(), ns), ip.String()),
			}
			to <- []define.Record{newRecord}
		} else {
			continue
		}
	}
}

func (s *NeighborScanner) ScanSvcNeighbor(subnet *net.IPNet) <-chan define.Record {
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
				s.scanSvc(subnet, out)
				defer s.wg.Done()
			}()
		} else {
			log.Debugf("Subnet split into %v success", len(subnets))
			s.wg.Add(len(subnets))
			for _, sn := range subnets {
				go func(sn *net.IPNet) {
					defer s.wg.Done()
					s.scanSvc(sn, out)
				}(sn)
			}
		}
		s.wg.Wait()
		close(out)
	}()
	return out
}

func (s *NeighborScanner) scanSvc(subnet *net.IPNet, to chan define.Record) {
	for _, ip := range pkg.ParseIPNetToIPs(subnet) {
		hostList := pkg.PTRRecord(ip)
		for _, host := range hostList {
			if post.IsPodServiceFormat(host) {
				newRecord := define.Record{
					Ip:        ip,
					SvcDomain: host,
				}
				to <- newRecord
			} else {
				log.Tracef("Pod Service: %v(%v) is not a pod service", host, ip.String())
				continue
			}
		}
	}
}
