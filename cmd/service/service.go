package service

import (
	"fmt"

	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg/printer"
	"github.com/esonhugh/k8spider/pkg/scanner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Opts struct {
	SvcDomains []string
}

func init() {
	ServiceCmd.PersistentFlags().StringSliceVarP(&Opts.SvcDomains, "svc-domains", "s", []string{}, "service domains, like: kubernetes.default,etcd.default don't add zone like svc.cluster.local")
	command.RootCmd.AddCommand(ServiceCmd)
}

var ServiceCmd = &cobra.Command{
	Use: "service",
	Aliases: []string{
		"srv",
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
