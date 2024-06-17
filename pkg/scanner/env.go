package scanner

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"regexp"
	"strings"
)

func extractValue(env string) (val string) {
	val = strings.Split(env, "=")[1]
	return strings.TrimSpace(val)
}

// detect service subnet from env
func FindServiceSubnet(env []string) (ServiceIPList []string) {
	pattern := regexp.MustCompile(`(.*)_SERVICE_HOST=`)
	if len(env) != 0 {
		for _, k := range env {
			matches := pattern.FindStringSubmatch(k)
			if matches == nil {
				continue
			}
			ServiceIPList = append(ServiceIPList, extractValue(k))
		}
	}
	return
}

func UniqueSubnet(serviceIPList ...string) (UniqueSubnetList []string) {
	if len(serviceIPList) != 0 {
		uniqueSubnets := make(map[string]bool)
		for _, ip := range serviceIPList {
			_, subnet, err := net.ParseCIDR(ip + "/16")
			if err != nil {
				log.Fatalf("Invalid IP/CIDR: %v", err)
			}

			// 获取子网的起始IP和网络地址
			subnetIP := subnet.IP.Mask(subnet.Mask)

			// 将网络地址作为唯一标识添加到map中
			key := subnetIP.String() + "/" + fmt.Sprint(subnet.Mask.Size())
			if !uniqueSubnets[key] {
				uniqueSubnets[key] = true
				UniqueSubnetList = append(UniqueSubnetList, subnet.String())
			}
		}
	}
	return
}
