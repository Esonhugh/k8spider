package all

import (
	"os"

	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	command.RootCmd.AddCommand(AllCmd)

}

var AllCmd = &cobra.Command{
	Use:   "all",
	Short: "all is a tool to discover k8s services and available ip in subnet",
	Run: func(cmd *cobra.Command, args []string) {
		if command.Opts.Cidr == "" {
			log.Warn("cidr is required")
			return
		}
		ipNets, err := pkg.ParseStringToIPNet(command.Opts.Cidr)
		if err != nil {
			log.Warnf("ParseStringToIPNet failed: %v", err)
			return
		}
		var records define.Records = pkg.ScanSubnet(ipNets)
		if records == nil || len(records) == 0 {
			log.Warnf("ScanSubnet Found Nothing: %v", err)
			return
		}
		records = pkg.ScanSvcForPorts(records)
		printResult(records)
		records = pkg.DumpAXFR(dns.Fqdn(command.Opts.Zone), "ns.dns."+command.Opts.Zone+":53")
		printResult(records)
	},
}

func printResult(records define.Records) {
	if command.Opts.OutputFile != "" {
		f, err := os.OpenFile(command.Opts.OutputFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Warnf("OpenFile failed: %v", err)
		}
		defer f.Close()
		records.Print(log.StandardLogger().Writer(), f)
	} else {
		records.Print(log.StandardLogger().Writer())
	}
}
