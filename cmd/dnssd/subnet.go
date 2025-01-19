package dnssd

import (
	"fmt"
	"net"

	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/esonhugh/k8spider/pkg/mutli"
	"github.com/esonhugh/k8spider/pkg/printer"
	"github.com/esonhugh/k8spider/pkg/scanner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	command.RootCmd.AddCommand(DNSSDCmd)
	DNSSDCmd.AddCommand(SubNetCmd, ServiceCmd)
	ServiceCmd.PersistentFlags().StringSliceVarP(&Opts.SvcDomains, "svc-domains", "s", []string{}, "service domains, like: kubernetes.default,etcd.default don't add zone like svc.cluster.local")
}

var DNSSDCmd = &cobra.Command{
	Use: "dnssd",
	Aliases: []string{
		"sd",
	},
	Short: "dnssd is a subcommand to discover k8s dns-sd technique",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var SubNetCmd = &cobra.Command{
	Use: "ptr",
	Aliases: []string{
		"sub",
		"s",
	},
	Short: "subnet is a tool to discover k8s available ip in subnet",
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
		var finalRecord define.Records
		finalRecord = RunMultiThread(ipNets, command.Opts.ThreadingNum)
		printer.PrintResult(finalRecord, command.Opts.OutputFile)
	},
}

func RunMultiThread(net *net.IPNet, num int) (finalRecord define.Records) {
	scan := mutli.NewSubnetScanner(num)
	for r := range scan.ScanSubnet(net) {
		finalRecord = append(finalRecord, r)
	}
	if len(finalRecord) == 0 {
		log.Warn("ScanSubnet Found Nothing")
		return
	}
	return
}

var Opts struct {
	SvcDomains []string
}

var ServiceCmd = &cobra.Command{
	Use: "srv",
	Aliases: []string{
		"service",
	},
	Short: "service is a tool to discover k8s services ports",
	Run: func(cmd *cobra.Command, args []string) {
		if command.Opts.Zone == "" || Opts.SvcDomains == nil || len(Opts.SvcDomains) == 0 {
			log.Warn("zone can't empty and svc-domains can't empty")
			return
		}
		var records define.Records
		for _, domain := range Opts.SvcDomains {
			records = append(records, define.Record{SvcDomain: fmt.Sprintf("%s.svc.%s", domain, command.Opts.Zone)})
		}
		records = scanner.ScanSvcForPorts(records)
		printer.PrintResult(records, command.Opts.OutputFile)
	},
}
