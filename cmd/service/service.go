package service

import (
	"fmt"
	"os"

	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg/scanner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	command.RootCmd.AddCommand(ServiceCmd)
}

var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "service is a tool to discover k8s services ports",
	Run: func(cmd *cobra.Command, args []string) {
		if command.Opts.Zone == "" || command.Opts.SvcDomains == nil || len(command.Opts.SvcDomains) == 0 {
			log.Warn("zone can't empty and svc-domains can't empty")
			return
		}
		var records define.Records
		for _, domain := range command.Opts.SvcDomains {
			records = append(records, define.Record{SvcDomain: fmt.Sprintf("%s.svc.%s", domain, command.Opts.Zone)})
		}
		records = scanner.ScanSvcForPorts(records)
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
	},
}
