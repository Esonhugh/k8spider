package post

import (
	"net"
	"strings"

	"github.com/esonhugh/k8spider/define"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

func RecordsDumpFullService(r []define.Record, zone string) []string {
	var result []string
	for _, record := range r {
		// ops svc domain
		if record.SvcDomain != "" && IsPodServiceFormat(record.SvcDomain, zone) {
			result = append(result, GetPodServiceRawService(record.SvcDomain, zone))
		} else if record.SvcDomain != "" && IsServiceFormat(record.SvcDomain, zone) {
			result = append(result, record.SvcDomain)
		}
		// ops srv records
		for _, srv := range record.SrvRecords {
			for _, s := range srv.Srv {
				if IsPodServiceFormat(s.Target, zone) {
					result = append(result, GetPodServiceRawService(s.Target, zone))
				} else if IsServiceFormat(s.Target, zone) {
					result = append(result, dns.Fqdn(s.Target))
				} else {
					log.Debugf("Unhandled service type: %v", s.Target)
				}
			}
		}
		R := ExtraParser(record.Extra)
		if R != "" {
			if IsServiceFormat(R, zone) {
				result = append(result, R)
			}
		}
	}
	return UniqueSlice(result)
}

func IsServiceFormat(domain string, zone string) bool {
	zonelen := len(strings.Split(dns.Fqdn(zone), "."))
	dn := ReverseSlice(strings.Split(dns.Fqdn(domain), "."))
	if len(dn) > 4 {
		if dn[zonelen] == "svc" { // Check if it is a service domain
			return true
		}
	}
	return false
}

func RecordsDumpNameSpace(r []define.Record, zone string) []string {
	result := RecordsDumpFullService(r, zone)
	for i, record := range result {
		result[i] = GetNamespaceFromDomain(record, zone)
	}
	return UniqueSlice(result)
}

func ExtraParser(record string) string {
	return strings.Split(record, " ")[0]
}

func GetNamespaceFromDomain(domain string, zone string) string {
	zonelen := len(strings.Split(dns.Fqdn(zone), "."))
	dn := ReverseSlice(strings.Split(dns.Fqdn(domain), "."))
	if len(dn) > 4 {
		if dn[zonelen] == "svc" { // Check if it is a service domain
			return dn[1+zonelen]
		}
	}
	return ""
}

func IsPodServiceFormat(domain, zone string) bool {
	if domain == "" {
		return false
	}
	str := strings.Split(strings.ReplaceAll(dns.Fqdn(domain), dns.Fqdn(zone), ""), ".")
	if len(str) > 3 {
		// 0:IP 1:ServiceName 2:Namespace 3:svc
		if str[3] == "svc" {
			return true
		}
	}
	return false
}

func GetPodServiceRawService(domain, zone string) string {
	str := strings.Split(strings.ReplaceAll(dns.Fqdn(domain), dns.Fqdn(zone), ""), ".")
	return dns.Fqdn(strings.Join(str[1:], ".") + zone)
}

func GetPodServiceRawIP(domain, zone string) net.IP {
	str := strings.Split(strings.ReplaceAll(dns.Fqdn(domain), dns.Fqdn(zone), ""), ".")
	return net.ParseIP(strings.ReplaceAll(str[0], "-", "."))
}

func PodServiceMap(BaseService define.Records, zone string) map[string][]string {
	result := make(map[string][]string)
	for _, r := range BaseService {
		log.Tracef("Processing Record: %v", r)
		if r.SvcDomain != "" && IsPodServiceFormat(r.SvcDomain, zone) {
			svcDomain := GetPodServiceRawService(r.SvcDomain, zone)
			result[svcDomain] = append(result[svcDomain], GetPodServiceRawIP(r.SvcDomain, zone).String())
		} else if r.SvcDomain != "" && IsServiceFormat(r.SvcDomain, zone) {
			if r.Ip != nil {
				result[dns.Fqdn(r.SvcDomain)] = append(result[dns.Fqdn(r.SvcDomain)], r.Ip.String())
			} else {
				log.Debugf("Lost service ip addr %v", r.SvcDomain)
			}
		} else {
			log.Debugf("Unhandled service type: %v", r.SvcDomain)
		}

		for _, srv := range r.SrvRecords {
			for _, s := range srv.Srv {
				if IsPodServiceFormat(s.Target, zone) {
					domain := GetPodServiceRawService(s.Target, zone)
					result[domain] = append(result[domain], GetPodServiceRawIP(s.Target, zone).String())
				} else if IsServiceFormat(s.Target, zone) {
					// result[dns.Fqdn(s.Target)]
					log.Debugf("can't put ip address in service %v in record %v", dns.Fqdn(s.Target), r)
				} else {
					log.Debugf("Unhandled service type: %v", s.Target)
				}
			}
		}
	}
	for k, v := range result {
		result[k] = UniqueSlice(v)
	}
	return result
}
