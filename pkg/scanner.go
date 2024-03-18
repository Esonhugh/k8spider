package pkg

import (
	"net"
	"strings"

	"github.com/esonhugh/k8spider/define"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

func ScanSubnet(subnet *net.IPNet) (records []define.Record) {
	for _, ip := range ParseIPNetToIPs(subnet) {
		ptr := PTRRecord(ip)
		if len(ptr) > 0 {
			for _, domain := range ptr {
				log.Infof("PTRrecord %v --> %v", ip, domain)
				r := define.Record{Ip: ip, SvcDomain: domain}
				records = append(records, r)
			}
		} else {
			continue
		}
	}
	return
}

func ScanSvcForPorts(records []define.Record) []define.Record {
	for i, r := range records {
		cname, srv, err := SRVRecord(r.SvcDomain)
		if err != nil {
			log.Debugf("SRVRecord for %v,failed: %v", r.SvcDomain, err)
			continue
		}
		for _, s := range srv {
			log.Infof("SRVRecord: %v --> %v:%v", r.SvcDomain, s.Target, s.Port)
		}
		records[i].SetSrvRecord(cname, srv)
	}
	return records
}

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
