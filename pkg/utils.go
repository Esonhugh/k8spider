package pkg

import (
	log "github.com/sirupsen/logrus"
)

// CheckPodVerified is utils to check if current Kubernetes has set pod verified
func CheckPodVerified() bool {
	iplist := []string{
		"8.8.8.8",
		"1.1.1.1",
		"114.114.114.114",
	}
	for _, ip := range iplist {
		targetHostName := IPtoPodHostName(ip, "kube-system")
		log.Tracef("test if record %v is ip: %v", targetHostName, ip)
		ips, err := ARecord(targetHostName)
		if err != nil {
			continue
		}
		// all of this ip should not return correct ip if verified is set.
		for _, i := range ips {
			if i.String() == ip {
				return false
			}
		}
	}
	return true
}

// https://github.com/kubernetes/dns/blob/master/docs/specification.md
// CheckKubernetes checks if the current environment is a kubernetes cluster
func CheckKubeDNS(dns ...*SpiderResolver) bool {
	if len(dns) > 1 {
		for _, d := range dns {
			if CheckKubeDNS(d) {
				log.Infof("kubernetes cluster found in dns(%v)", d.CurrentDNS())
				return true
			}
		}
	} else { // less or equal to 1
		rs := NetResolver
		if len(dns) > 0 { // use custom dns
			rs = dns[0]
		}
		if CheckKubeDNS_DefaultAPIServer(rs) || CheckKubeDNS_DNSVersion(rs) || CheckKubeDNS_NS_DNS_DOMAIN(rs) {
			return true
		}
	}
	return false
}

func CheckKubeDNS_DefaultAPIServer(dns *SpiderResolver) bool {
	info, err := dns.ARecord("kubernetes.default.svc." + Zone)
	if err == nil {
		log.Debugf("kubernetes.default.svc.%v found in dns(%v)! response: %v", Zone, dns.CurrentDNS(), info)
		return true
	}
	log.Tracef("kubernetes.default.svc.%v not found in dns(%v)", Zone, dns.CurrentDNS())
	info, err = dns.ARecord("kubernetes.default.svc") // try shorted if Zone is incorrect
	if err == nil {
		log.Warnf("kubernetes.default.svc found in dns(%v)! response: %v, maybe %v is incorrect", dns.CurrentDNS(), info, Zone)
		return true
	}
	log.Tracef("kubernetes.default.svc not found in dns(%v)", dns.CurrentDNS())
	return false
}

func CheckKubeDNS_DNSVersion(dns *SpiderResolver) bool {
	info, err := dns.TXTRecord("dns-version." + Zone)
	if err == nil {
		log.Debugf("dns-version.%v found in dns(%v)! response: %v", Zone, dns.CurrentDNS(), info)
		return true
	}
	log.Tracef("dns-version.%v not found in dns(%v)", Zone, dns.CurrentDNS())
	info, err = dns.TXTRecord("dns-version") // try shorted if Zone is incorrect
	if err == nil {
		log.Warnf("dns-version found in dns(%v)! response: %v, maybe %v is incorrect", dns.CurrentDNS(), info, Zone)
		return true
	}
	log.Tracef("dns-version not found in dns(%v)", dns.CurrentDNS())
	return false
}

func CheckKubeDNS_NS_DNS_DOMAIN(dns *SpiderResolver) bool {
	info, err := dns.ARecord("ns.dns." + Zone)
	if err == nil {
		log.Debugf("ns.dns.%v found in dns(%v)! response: %v", Zone, dns.CurrentDNS(), info)
		return true
	}
	log.Tracef("ns.dns.%v not found in dns(%v)", Zone, dns.CurrentDNS())
	info, err = dns.ARecord("ns.dns") // try shorted if Zone is incorrect
	if err == nil {
		log.Warnf("ns.dns found in dns(%v)! response: %v, maybe %v is incorrect", dns.CurrentDNS(), info, Zone)
		return true
	}
	log.Tracef("ns.dns not found in dns(%v)", dns.CurrentDNS())
	return false
}
