package scanner

import (
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

func DumpWildCard(zone string) []define.Record {
	searchDNS := []string{
		dns.Fqdn("any.any.svc." + zone),
		dns.Fqdn("any.any.any.svc." + zone),
	}
	var records []define.Record
	for _, dns := range searchDNS {
		_, srv, err := pkg.SRVRecord(dns)
		if err != nil {
			log.Warnf("wildcard dns query to %v failed: %v", dns, err)
			continue
		}
		r := define.Record{}
		r.SetSrvRecord(dns, srv)
		records = append(records, r)
	}
	return records
}
