package axfr

import (
	"os"
	"strings"

	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg/scanner"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	command.RootCmd.AddCommand(AxfrCmd)
}

var AxfrCmd = &cobra.Command{
	Use:   "axfr",
	Short: "axfr is a command to dump every record from dns server",
	Run: func(cmd *cobra.Command, args []string) {

		if command.Opts.Zone == "" {
			log.Warn("zone can't empty")
			return
		}
		zone := dns.Fqdn(command.Opts.Zone)

		dnsServer := command.Opts.DnsServer
		if command.Opts.DnsServer == "" {
			dnsServer = "ns.dns." + command.Opts.Zone + ":53"
		} else if len(strings.Split(dnsServer, ":")) < 2 {
			dnsServer = dnsServer + ":53"
		}

		log.Debugf("same command: dig axfr %v @%v", zone, dnsServer)
		var records define.Records
		records, err := scanner.DumpAXFR(zone, dnsServer)
		if err != nil {
			log.Errorf("Transfer failed: %v", err)
			return
		}
		if command.Opts.OutputFile != "" {
			f, err := os.OpenFile(command.Opts.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Warnf("OpenFile failed: %v", err)
			}
			defer f.Close()
			records.Print(log.StandardLogger().Writer(), f)
		} else {
			records.Print(log.StandardLogger().Writer())
		}

	},
}
