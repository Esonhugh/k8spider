package mutli

import (
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg/scanner"
)

func ScanServiceWithChan(rev <-chan []define.Record) <-chan []define.Record {
	out := make(chan []define.Record, 100)
	go func() {
		for records := range rev {
			out <- scanner.ScanSvcForPorts(records)
		}
		close(out)
	}()
	return out
}
