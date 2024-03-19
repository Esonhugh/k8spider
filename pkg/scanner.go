package pkg

import (
	"net"
	"strings"
	"sync"

	"github.com/esonhugh/k8spider/define"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

func ScanSubnet(subnet *net.IPNet, thread uint) (records []define.Record) {
	threadLimit := make(chan struct{}, thread)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, ip := range ParseIPNetToIPs(subnet) {
		wg.Add(1)
		threadLimit <- struct{}{}
		go func(ip net.IP) {
			defer wg.Done()
			defer func() { <-threadLimit }()
			ptr := PTRRecord(ip)
			if len(ptr) > 0 {
				for _, domain := range ptr {
					log.Infof("PTRrecord %v --> %v", ip, domain)
					r := define.Record{Ip: ip, SvcDomain: domain}
					mu.Lock()
					records = append(records, r)
					mu.Unlock()
				}
			} else {
				return
			}
		}(ip)
	}
	wg.Wait()
	close(threadLimit)
	return
}

func ScanSvcForPorts(records []define.Record, thread uint) []define.Record {
	threadLimit := make(chan struct{}, thread)
	var wg sync.WaitGroup
	for i, r := range records {
		wg.Add(1)
		threadLimit <- struct{}{}
		go func(i int, r define.Record) {
			defer wg.Done()
			defer func() { <-threadLimit }()
			cname, srv, err := SRVRecord(r.SvcDomain)
			if err != nil {
				log.Debugf("SRVRecord for %v,failed: %v", r.SvcDomain, err)
				return
			}
			for _, s := range srv {
				log.Infof("SRVRecord: %v --> %v:%v", r.SvcDomain, s.Target, s.Port)
			}
			records[i].SetSrvRecord(cname, srv)
		}(i, r)
	}
	wg.Wait()
	close(threadLimit)
	return records
}

// DumpAXFR default target should be zone
func DumpAXFR(target string, dnsServer string, thread uint) []define.Record {
	threadLimit := make(chan struct{}, thread)
	var wg sync.WaitGroup
	var mu sync.Mutex

	t := new(dns.Transfer)
	m := new(dns.Msg)
	m.SetAxfr(target)
	ch, err := t.In(m, dnsServer)
	if err != nil {
		log.Fatalf("Transfer failed: %v", err)
	}
	var records []define.Record
	for rr := range ch {
		wg.Add(1)
		threadLimit <- struct{}{}
		go func(rr *dns.Envelope) {
			defer wg.Done()
			defer func() { <-threadLimit }()
			for _, r := range rr.RR {
				mu.Lock()
				records = append(records, define.Record{
					SvcDomain: r.Header().Name,
					Extra:     strings.Join(strings.Split(r.String(), "\t"), " "),
				})
				mu.Unlock()
			}
			log.Debugf("Record: %v", rr.RR)
		}(rr)
	}
	wg.Wait()
	return records
}
