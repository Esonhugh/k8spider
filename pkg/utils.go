package pkg

import (
	"context"
	"encoding/binary"
	"net"

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

func PTRRecord(ip net.IP) []string {
	names, err := NetResolver.LookupAddr(context.Background(), ip.String())
	if err != nil {
		log.Errorf("LookupAddr failed: %v", err)
		return nil
	}
	return names
}

func SRVRecord(svcDomain string) (string, []*net.SRV, error) {
	cname, srvs, err := NetResolver.LookupSRV(context.Background(), "", "", svcDomain)
	return cname, srvs, err
}
