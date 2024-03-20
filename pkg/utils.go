package pkg

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

var NetResolver *net.Resolver = net.DefaultResolver

func ParseStringToIPNet(s string) (ipnet *net.IPNet, err error) {
	_, ipnet, err = net.ParseCIDR(s)
	return
}

func ParseIPNetToIPs(ipv4Net *net.IPNet) (ips []net.IP) {
	// convert IPNet struct mask and address to uint32
	// network is BigEndian
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)

	// find the final address
	finish := (start & mask) | (mask ^ 0xffffffff)

	// loop through addresses as uint32
	for i := start; i <= finish; i++ {
		// convert back to net.IP
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip)
	}
	return
}

var RetryCount = 2

func retryCoolDown() {
	time.Sleep(10 * time.Millisecond)
}

func PTRRecord(ip net.IP) []string {
	var lastErr error
	for i := 0; i < RetryCount; i++ {
		names, err := NetResolver.LookupAddr(context.Background(), ip.String())
		if err == nil {
			return names
		}
		lastErr = err
		retryCoolDown()
	}
	log.Debugf("LookupAddr %v failed, no service addr found. reason: %v", ip.String(), lastErr)
	return nil
}

func SRVRecord(svcDomain string) (string, []*net.SRV, error) {
	var lastErr error
	for i := 0; i < RetryCount; i++ {
		cname, srvs, err := NetResolver.LookupSRV(context.Background(), "", "", svcDomain)
		if err == nil {
			return cname, srvs, nil
		}
		lastErr = err
		retryCoolDown()
	}
	log.Debugf("LookupSRV %v failed, no ports record found. reason: %v", svcDomain, lastErr)
	return "", nil, fmt.Errorf("LookupSRV %v failed", svcDomain)
}

func ARecord(domain string) (ips []net.IP, err error) {
	ips, err = net.LookupIP(domain)
	return
}
