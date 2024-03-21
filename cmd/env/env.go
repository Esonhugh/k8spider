package env

import (
	command "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/pkg/scanner"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	command.RootCmd.AddCommand(EnvCmd)
}

var EnvCmd = &cobra.Command{
	Use:   "env",
	Short: "env is a command to detect service subnet from environ",
	Run: func(cmd *cobra.Command, args []string) {
		env := os.Environ()
		serviceSubnetList := scanner.FindServiceSubnet(env)
		serviceSubnetList = scanner.UniqueSubnet(serviceSubnetList...)
		for _, k := range serviceSubnetList {
			log.Infoln("Subnet from env:", k)
		}
	},
}
