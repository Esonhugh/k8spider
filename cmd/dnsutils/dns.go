package dnsutils

import (
	"strings"

	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var queryType string

func init() {
	DNSCmd.PersistentFlags().StringVarP(&queryType, "type", "x", "A", "query type")
	command.RootCmd.AddCommand(DNSCmd)
}

var DNSCmd = &cobra.Command{
	Use:     "dns",
	Aliases: []string{"dig"},
	Short:   "dns is a command to query dns server",
	Run: func(cmd *cobra.Command, args []string) {
		var querier pkg.DnsQuery
		switch strings.ToLower(queryType) {
		case "a", "aaaa":
			querier = pkg.QueryA
		case "ptr":
			querier = pkg.QueryPTR
		case "srv":
			querier = pkg.QuerySRV
		case "txt":
			querier = pkg.QueryTXT
		default:
			querier = pkg.QueryA
		}
		for _, query := range args {
			res, err := querier(query)
			if err != nil {
				log.Warnf("Query %s failed: %v", query, err)
				continue
			}
			log.Infof("Query [%s] %s: %v", queryType, query, res)
		}
	},
}
