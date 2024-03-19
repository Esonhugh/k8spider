package pkg

import (
	"github.com/esonhugh/k8spider/define"
	log "github.com/sirupsen/logrus"
)

var _ StaticTaskRunner[*define.Record] = (*SvcTask)(nil)

type SvcTask struct {
	Records []define.Record
}

func (s *SvcTask) Generator(tasks chan *define.Record) {
	for i := range s.Records {
		tasks <- &s.Records[i]
	}
}

func (s *SvcTask) Solver(r *define.Record) {
	cname, srv, err := SRVRecord(r.SvcDomain)
	if err != nil {
		log.Debugf("SRVRecord for %v,failed: %v", r.SvcDomain, err)
		return
	}
	for _, s := range srv {
		log.Infof("SRVRecord: %v --> %v:%v", r.SvcDomain, s.Target, s.Port)
	}
	r.SetSrvRecord(cname, srv)
}

func ScanSvcForPorts(records []define.Record) []define.Record {
	RunStatic[*define.Record](&SvcTask{Records: records}, Threads)
	return records
}
