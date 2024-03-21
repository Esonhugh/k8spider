package wildcard

import (
	"os"

	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/define"
	"github.com/esonhugh/k8spider/pkg/scanner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	command.RootCmd.AddCommand(WildCardCmd)
}

var WildCardCmd = &cobra.Command{
	Use:   "wild",
	Short: "wild is a tool to abuse wildcard feature in kubernetes service discovery",
	Run: func(cmd *cobra.Command, args []string) {
		if command.Opts.Zone == "" {
			log.Warn("zone can't empty")
			return
		}
		printResult(scanner.DumpWildCard(command.Opts.Zone))
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
