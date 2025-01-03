package mutli

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/esonhugh/k8spider/pkg/post"
	"github.com/esonhugh/k8spider/pkg/scanner"
	log "github.com/sirupsen/logrus"
)

type NeighborScanner struct {
	wg    *sync.WaitGroup
	count int
	m     ScanMode
}

type ScanMode string

const (
	ScanPods ScanMode = "pods"
	ScanSvc  ScanMode = "svc"
)

func NewNeighborScanner(mode ScanMode, threading ...int) *NeighborScanner {
	if len(threading) == 0 {
		return &NeighborScanner{
			wg: new(sync.WaitGroup),
		}
	} else {
		return &NeighborScanner{
			wg:    new(sync.WaitGroup),
			count: threading[0],
		}
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
			go s.scan(ns, subnet, out)
		} else {
			log.Debugf("Subnet split into %v success", len(subnets))
			for _, sn := range subnets {
				go s.scan(ns, sn, out)
			}
		}
		time.Sleep(10 * time.Millisecond) // wait for all goroutines to start
		s.wg.Wait()
		close(out)
	}()
	return out
}

func (s *NeighborScanner) ScanMultiNeighbor(nss []string, subnet *net.IPNet) <-chan []define.Record {
	out := make(chan []define.Record, 100)
	go func() {
		for _, ns := range nss {
			go s.scan(ns, subnet, out)
		}
		time.Sleep(10 * time.Millisecond) // wait for all goroutines to start
		s.wg.Wait()
		close(out)
	}()
	return out
}

func (s *NeighborScanner) scan(ns string, subnet *net.IPNet, to chan []define.Record) {
	s.wg.Add(1)
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
	s.wg.Done()
}

func (s *NeighborScanner) ScanSvcNeighbor(subnet *net.IPNet) <-chan []define.Record {
	if subnet == nil {
		log.Debugf("subnet is nil")
		return nil
	}
	out := make(chan []define.Record, 100)
	go func() {
		// if subnets, err := pkg.SubnetShift(subnet, 4); err != nil {
		if subnets, err := pkg.SubnetInto(subnet, s.count); err != nil {
			log.Errorf("Subnet split into %v failed, fallback to single mode, reason: %v", s.count, err)
			go s.scanSvc(subnet, out)
		} else {
			log.Debugf("Subnet split into %v success", len(subnets))
			for _, sn := range subnets {
				go s.scanSvc(sn, out)
			}
		}
		time.Sleep(10 * time.Millisecond) // wait for all goroutines to start
		s.wg.Wait()
		close(out)
	}()
	return out
}

func (s *NeighborScanner) scanSvc(subnet *net.IPNet, to chan []define.Record) {
	s.wg.Add(1)
	for _, ip := range pkg.ParseIPNetToIPs(subnet) {
		hostList := pkg.PTRRecord(ip)
		for _, host := range hostList {
			if post.IsPodServiceFormat(host) {
				newRecord := define.Record{
					Ip:        ip,
					SvcDomain: strings.SplitN(host, ".", 1)[1],
				}
				to <- []define.Record{newRecord}
			} else {
				continue
			}
		}
	}
	s.wg.Done()
}
