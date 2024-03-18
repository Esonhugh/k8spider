package define

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Record struct {
	Ip         net.IP
	SvcDomain  string
	SrvRecords []SrvRecord
}

type SrvRecord struct {
	Cname string
	Srv   []*net.SRV
}

func (r *Record) SetSrvRecord(cname string, srv []*net.SRV) {
	if r.SrvRecords == nil {
		r.SrvRecords = make([]SrvRecord, 0)
	}
	r.SrvRecords = append(r.SrvRecords, SrvRecord{Cname: cname, Srv: srv})
}

func (r *Record) Print(writer ...io.Writer) {
	var W io.Writer
	if len(writer) == 0 {
		W = os.Stdout
	} else {
		W = io.MultiWriter(writer...)
	}
	data := fmt.Sprintf("Found svc: %v at IP %v\n", r.SvcDomain, r.Ip)
	if len(r.SrvRecords) != 0 {
		for _, srv := range r.SrvRecords {
			if srv.Srv == nil || len(srv.Srv) == 0 {
				data += "But No Service Endpoint Found\n"
				continue
			} else {
				data += fmt.Sprintf("Found %v's Endpoint:\n", srv.Cname)
				for _, s := range srv.Srv {
					data += fmt.Sprintf("- %v:%v \n", s.Target, s.Port)
				}
			}
		}
	}
	_, _ = fmt.Fprint(W, data)
}

type Records []Record

func (r Records) Print(writer ...io.Writer) {
	for _, record := range r {
		record.Print(writer...)
	}
}
