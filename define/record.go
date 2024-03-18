package define

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
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
	data, err := json.Marshal(r)
	if err != nil {
		log.Error(err)
		return
	}
	_, _ = fmt.Fprintf(W, "%v\n", string(data))
}

type Records []Record

func (r Records) Print(writer ...io.Writer) {
	for _, record := range r {
		record.Print(writer...)
	}
}
