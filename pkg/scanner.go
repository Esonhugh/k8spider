package pkg

import (
	"net"

	"github.com/esonhugh/k8spider/define"
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
		log.Infof("SRVRecord: %v --> %v", r.SvcDomain, srv)
		records[i].SetSrvRecord(cname, srv)
	}
	return records
}
