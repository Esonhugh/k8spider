package wildcard

import (
	"encoding/json"
	command "github.com/esonhugh/k8spider/cmd"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"os"
	"sort"
	"strings"
)

func init() {
	command.RootCmd.AddCommand(WildCardCmd)
}

var WildCardCmd = &cobra.Command{
	Use:   "wildcard",
	Short: "wildcard is a command to dump every record from CoreDNS",
	Run: func(cmd *cobra.Command, args []string) {
		if command.Opts.Zone == "" {
			log.Warn("zone can't empty")
			return
		}

		zone := dns.Fqdn(command.Opts.Zone)
		DumpWildCard(zone)
	},
}

func DumpWildCard(zone string) {
	dnsNames := []string{
		"any.any.svc." + zone,
		"any.any.any.svc." + zone,
	}

	log.Debugf("same command: dig +short  %s %s", dnsNames[0], dnsNames[1])
	var results []*net.SRV
	for _, name := range dnsNames {
		_, srvs, err := net.LookupSRV("", "", name)
		if err != nil {
			log.Warnln(err.Error())
			continue
		}

		results = append(results, srvs...)
	}

	if len(results) == 0 {
		log.Warn("CoreDNS Wildcards Found Nothing")
		return
	}

	sort.Slice(results, func(i, j int) bool {
		switch strings.Compare(results[i].Target, results[j].Target) {
		case -1:
			return true
		case 0:
			return results[i].Port < results[j].Port
		case 1:
			return false
		}
		return false
	})

	if command.Opts.OutputFile != "" {
		f, err := os.OpenFile(command.Opts.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Warnf("OpenFile failed: %v", err)
		}
		data, err := json.Marshal(results)
		if err != nil {
			log.Error(err)
			return
		}
		data = append(data, '\n')
		_, err = f.Write(data)
		if err != nil {
			log.Error(err)
			return
		}
		defer f.Close()
	}

	for _, srv := range results {
		log.Println(srv.Target, srv.Port)
	}
}
