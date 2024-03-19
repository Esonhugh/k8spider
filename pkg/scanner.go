package pkg

import (
	"github.com/esonhugh/k8spider/define"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"strings"
)

var Threads = 4

// default target should be zone
func DumpAXFR(target string, dnsServer string) []define.Record {
	t := new(dns.Transfer)
	m := new(dns.Msg)
	m.SetAxfr(target)
	ch, err := t.In(m, dnsServer)
	if err != nil {
		log.Fatalf("Transfer failed: %v", err)
	}
	var records []define.Record
	for rr := range ch {
		if rr.Error != nil {
			log.Errorf("Error: %v", rr.Error)
			continue
		}
		for _, r := range rr.RR {
			records = append(records, define.Record{
				SvcDomain: r.Header().Name,
				Extra:     strings.Join(strings.Split(r.String(), "\t"), " "),
			})
		}
		log.Debugf("Record: %v", rr.RR)
	}
	return records
}
