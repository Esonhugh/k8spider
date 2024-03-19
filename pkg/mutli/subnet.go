package mutli

import (
	"net"
	"sync"

	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/esonhugh/k8spider/pkg/scanner"
	log "github.com/sirupsen/logrus"
)

type SubnetScanner struct {
	wg *sync.WaitGroup
}

func NewSubnetScanner() *SubnetScanner {
	return &SubnetScanner{
		wg: new(sync.WaitGroup),
	}
}

func (s *SubnetScanner) ScanSubnet(subnet *net.IPNet) <-chan []define.Record {
	if subnet == nil {
		log.Tracef("subnet is nil")
		return nil
	}
	out := make(chan []define.Record, 100)
	go func() {
		if subnets, err := pkg.SubnetShift(subnet, 4); err != nil {
			go s.scan(subnet, out)
		} else {
			for _, sn := range subnets {
				go s.scan(sn, out)
			}
		}
		s.wg.Wait()
		close(out)
	}()
	return out
}

func (s *SubnetScanner) scan(subnet *net.IPNet, to chan []define.Record) {
	s.wg.Add(1)
	to <- scanner.ScanSubnet(subnet)
	s.wg.Done()
}
