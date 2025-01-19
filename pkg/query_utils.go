package pkg

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	DnsTimeout  = 2
	NetResolver *SpiderResolver
	Zone        string // Zone is the domain name of the cluster

	Latency = 0
)

type SpiderResolver struct {
	dns      string
	ctx      context.Context
	r        *net.Resolver
	filter   []*regexp.Regexp
	contains []string

	lock *sync.Mutex
}

func DefaultResolver() *SpiderResolver {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(DnsTimeout)*time.Second) // I don't think if a inside cluster dns query has more than 2s latency.
	return &SpiderResolver{
		dns:      "default-dns",
		r:        net.DefaultResolver,
		ctx:      ctx,
		filter:   []*regexp.Regexp{},
		contains: []string{},
	}
}

func (r *SpiderResolver) SetFilter(filters ...string) {
	for _, filter := range filters {
		r.filter = append(r.filter, regexp.MustCompile(filter))
	}
}

func (r *SpiderResolver) SetContainsFilter(name ...string) {
	r.contains = append(r.contains, name...)
}

func (r *SpiderResolver) SetSuffixFilter(filter string) {
	r.SetFilter(filter + "$")
}

func (r *SpiderResolver) filterString(target string) bool {
	log.Tracef("filtering %s", target)
	for _, re := range r.filter {
		if re.MatchString(target) {
			log.Tracef("target %s matched regexp rule %s", target, re.String())
			return true
		}
	}
	for _, re := range r.contains {
		if strings.Contains(target, re) {
			log.Tracef("target %s matched contains rule %s", target, re)
			return true
		}
	}
	return false
}

func (r *SpiderResolver) filterStringArray(target []string) []string {
	var filtered []string
	for _, re := range target {
		if r.filterString(re) {
			continue
		}
		filtered = append(filtered, re)
	}
	log.Tracef("filtering %s \nresult: %s", strings.Join(target, " "), strings.Join(filtered, " "))
	return filtered
}

func WarpDnsServer(dnsServer string) *SpiderResolver {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(DnsTimeout)*time.Second)
	return &SpiderResolver{
		dns: dnsServer,
		r: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{}
				return d.DialContext(ctx, network, dnsServer)
			},
		},
		ctx:      ctx,
		filter:   []*regexp.Regexp{},
		contains: []string{},
	}
}

func (s *SpiderResolver) CurrentDNS() string {
	return s.dns
}

func (s *SpiderResolver) PTRRecord(ip net.IP) []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	names, err := s.r.LookupAddr(s.ctx, ip.String())
	if err != nil {
		log.Debugf("LookupAddr failed: %v", err)
		return nil
	}
	time.Sleep(time.Duration(Latency) * time.Millisecond)
	return s.filterStringArray(names)
}

func (s *SpiderResolver) SRVRecord(svcDomain string) (string, []*net.SRV, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	cname, srvs, err := s.r.LookupSRV(s.ctx, "", "", svcDomain)
	var finalsrv []*net.SRV
	for _, srv := range srvs {
		if s.filterString(srv.Target) {
			continue
		}
		finalsrv = append(finalsrv, srv)
	}
	time.Sleep(time.Duration(Latency) * time.Millisecond)
	return cname, srvs, err
}

func (s *SpiderResolver) CustomSRVRecord(svcDomain string, service, proto string) (string, []*net.SRV, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	cname, srvs, err := s.r.LookupSRV(s.ctx, service, proto, svcDomain)
	time.Sleep(time.Duration(Latency) * time.Millisecond)
	return cname, srvs, err
}

func (s *SpiderResolver) ARecord(domain string) ([]net.IP, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	time.Sleep(time.Duration(Latency) * time.Millisecond)
	return s.r.LookupIP(s.ctx, "ip", domain)
}

func (s *SpiderResolver) TXTRecord(domain string) ([]string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	time.Sleep(time.Duration(Latency) * time.Millisecond)
	return s.r.LookupTXT(s.ctx, domain)
}

func PTRRecord(ip net.IP) []string {
	return NetResolver.PTRRecord(ip)
}

func SRVRecord(svcDomain string) (string, []*net.SRV, error) {
	return NetResolver.SRVRecord(svcDomain)
}

func ARecord(domain string) (ips []net.IP, err error) {
	return NetResolver.ARecord(domain)
}

func TXTRecord(domain string) (txts []string, err error) {
	return NetResolver.TXTRecord(domain)
}

type DnsQuery func(domain string) ([]string, error)

var (
	QueryPTR DnsQuery = func(domain string) ([]string, error) {
		return PTRRecord(net.ParseIP(domain)), nil
	}
	QueryA DnsQuery = func(domain string) ([]string, error) {
		res, err := ARecord(domain)
		var ret []string
		for _, r := range res {
			ret = append(ret, r.String())
		}
		return ret, err
	}
	QueryTXT DnsQuery = TXTRecord
	QuerySRV DnsQuery = func(domain string) ([]string, error) {
		_, res, err := SRVRecord(domain)
		var ret []string
		for _, r := range res {
			ret = append(ret, fmt.Sprintf("%s:%d", r.Target, r.Port))
		}
		return ret, err
	}
)
