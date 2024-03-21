package scanner

import (
	"net"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

var BarTemplate = `{{ string . "subnet" }} {{ bar . "[" "-" (cycle . ">" ) "." "]"}} {{speed .  }} {{percent .}}`

func ScanSubnet(subnet *net.IPNet, optPb *pb.ProgressBar) (records []define.Record) {
	subnetIPlist := pkg.ParseIPNetToIPs(subnet)
	var bar *pb.ProgressBar
	if optPb != nil {
		bar = optPb.SetTemplateString(BarTemplate).Start()
	} else {
		bar = pb.ProgressBarTemplate(BarTemplate).Start(len(subnetIPlist))
	}
	bar.Set("subnet", subnet.String())
	for _, ip := range subnetIPlist {
		ptr := pkg.PTRRecord(ip)
		if len(ptr) > 0 {
			for _, domain := range ptr {
				log.Infof("PTRrecord %v --> %v", ip, domain)
				r := define.Record{Ip: ip, SvcDomain: domain}
				records = append(records, r)
			}
		}
		bar.Increment()
	}
	bar.Finish()
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

// default target should be zone
func DumpAXFR(target string, dnsServer string) ([]define.Record, error) {
	t := new(dns.Transfer)
	m := new(dns.Msg)
	m.SetAxfr(target)
	ch, err := t.In(m, dnsServer)
	if err != nil {
		return nil, err
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
	return records, nil
}
