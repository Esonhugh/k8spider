package subnet

import (
	"net"
	"os"

	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg"
	"github.com/esonhugh/k8spider/pkg/mutli"
	"github.com/esonhugh/k8spider/pkg/scanner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	command.RootCmd.AddCommand(SubNetCmd)
}

var SubNetCmd = &cobra.Command{
	Use:   "subnet",
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
		if command.Opts.BatchMode {
			BatchRun(ipNets)
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
	printResult(records)
}

func BatchRun(net *net.IPNet) {
	scan := mutli.NewSubnetScanner()
	var finalRecord []define.Record
	for r := range scan.ScanSubnet(net) {
		finalRecord = append(finalRecord, r...)
	}
	if len(finalRecord) == 0 {
		log.Warn("ScanSubnet Found Nothing")
		return
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
