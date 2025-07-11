package nfs

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	NFSCommand.AddCommand(MkdirCommand)
}

var MkdirCommand = &cobra.Command{
	Use:     "mkdir",
	Aliases: []string{"md"},
	Short:   "create nfs directory",
	Long:    "create nfs directory",
	PreRunE: NotInit,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Errorf("mkdir: missing argument")
			return
		}
		for _, arg := range args {
			err := Client.Mkdir(arg)
			if err != nil {
				log.Error(err)
				continue
			}
			log.Infof("mkdir %s success", arg)
		}
	},
}
