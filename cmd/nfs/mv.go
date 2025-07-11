package nfs

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	NFSCommand.AddCommand(RenameCommand)
}

var RenameCommand = &cobra.Command{
	Use:     "rename [f1] [f2]",
	Short:   "rename nfs file",
	Long:    "rename nfs file",
	Aliases: []string{"mv"},
	PreRunE: NotInit,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Errorf("rename: missing arguments")
			return
		}
		err := Client.Rename(args[0], args[1])
		if err != nil {
			log.Error(err)
			return
		}
		log.Infof("rename %s to %s success\n", args[0], args[1])
	},
}
