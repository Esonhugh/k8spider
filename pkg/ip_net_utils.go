package pkg

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

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

func IPtoPodHostName(ip, namespace string) string {
	return fmt.Sprintf("%s.%s.pod.%s", strings.ReplaceAll(ip, ".", "-"), namespace, Zone)
}
