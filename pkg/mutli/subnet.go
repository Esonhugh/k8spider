package mutli

import (
	"net"
	"sync"

	"github.com/cheggaaa/pb/v3"
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
		var pblist []*pb.ProgressBar
		var pool *pb.Pool
		if subnets, err := pkg.SubnetShift(subnet, 4); err != nil {
			go s.scan(subnet, out, nil)
		} else {
			// start pool
			for _, sn := range subnets {
				pblist = append(pblist, pb.New(len(pkg.ParseIPNetToIPs(sn))))
			}
			pool, err = pb.StartPool(pblist...)
			if err != nil {
				panic(err)
			}
			for i, sn := range subnets {
				go s.scan(sn, out, pblist[i])
			}
		}
		s.wg.Wait()
		close(out)
		if len(pblist) > 0 {
			_ = pool.Stop()
		}
	}()
	return out
}

func (s *SubnetScanner) scan(subnet *net.IPNet, to chan []define.Record, pb *pb.ProgressBar) {
	s.wg.Add(1)
	to <- scanner.ScanSubnet(subnet, pb)
	s.wg.Done()
}
