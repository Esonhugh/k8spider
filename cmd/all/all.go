package all

import (
	"github.com/esonhugh/k8spider/cmd/wildcard"
	"net"
	"os"

	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/esonhugh/k8spider/pkg/mutli"
	"github.com/esonhugh/k8spider/pkg/scanner"
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
		wildcard.DumpWildCard(command.Opts.Zone)
		records, err := scanner.DumpAXFR(dns.Fqdn(command.Opts.Zone), "ns.dns."+command.Opts.Zone+":53")
		if err == nil {
			printResult(records)
		}
		log.Errorf("Transfer failed: %v", err)
		ipNets, err := pkg.ParseStringToIPNet(command.Opts.Cidr)
		if err != nil {
			log.Warnf("ParseStringToIPNet failed: %v", err)
			return
		}
		if command.Opts.BatchMode {
			RunBatch(ipNets)
		} else {
			Run(ipNets)
		}
	},
}

func Run(net *net.IPNet) {
	var records define.Records = scanner.ScanSubnet(net)
	if records == nil || len(records) == 0 {
		log.Warnf("ScanSubnet Found Nothing")
		return
	}
	records = scanner.ScanSvcForPorts(records)
	printResult(records)
}

func RunBatch(net *net.IPNet) {
	scan := mutli.ScanAll(net)
	var finalRecord []define.Record
	for r := range scan {
		finalRecord = append(finalRecord, r...)
	}
	printResult(finalRecord)
}

func printResult(records define.Records) {
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
}
