package scanner

import (
	"net"

	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	log "github.com/sirupsen/logrus"
)

func ScanSingleIP(subnet net.IP) (records []define.Record) {
	ptr := pkg.PTRRecord(subnet)
	if len(ptr) > 0 {
		for _, domain := range ptr {
			log.Infof("PTRrecord %v --> %v", subnet, domain)
			r := define.Record{Ip: subnet, SvcDomain: domain}
			records = append(records, r)
		}
	}
	return
}

func ScanSubnet(subnet *net.IPNet) (records []define.Record) {
	for _, ip := range pkg.ParseIPNetToIPs(subnet) {
		ptr := pkg.PTRRecord(ip)
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

func ScanSingleSvcForPorts(records define.Record) define.Record {
	cname, srv, err := pkg.SRVRecord(records.SvcDomain)
	if err != nil {
		log.Debugf("SRVRecord for %v,failed: %v", records.SvcDomain, err)
		return records
	}
	for _, s := range srv {
		log.Infof("SRVRecord: %v --> %v:%v", records.SvcDomain, s.Target, s.Port)
	}
	records.SetSrvRecord(cname, srv)
	return records
}

func ScanSvcForPorts(records []define.Record) []define.Record {
	for i, r := range records {
		cname, srv, err := pkg.SRVRecord(r.SvcDomain)
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
