package nfs

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rmOpts struct {
	Force bool
}

func init() {
	RemoveCommand.Flags().BoolVar(&rmOpts.Force, "force", false, "force remove")
	NFSCommand.AddCommand(RemoveCommand, RemoveDirCommand)
}

var RemoveCommand = &cobra.Command{
	Use:     "rm [f1] [f2] ...",
	Short:   "remove nfs file",
	Long:    "remove nfs file",
	Aliases: []string{"del", "delete"},
	PreRunE: NotInit,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Errorf("rm: missing arguments")
			return
		}
		for _, arg := range args {
			var err error
			if rmOpts.Force {
				err = Client.RemoveAll(arg)
			} else {
				err = Client.Remove(arg)
			}
			if err != nil {
				log.Error(err)
				continue
			}
			log.Infof("remove %s success\n", arg)
		}
	},
}

var RemoveDirCommand = &cobra.Command{
	Use:     "rmdir [dir1] [dir2] ...",
	Short:   "remove nfs directory",
	Long:    "remove nfs directory",
	Aliases: []string{"rmdir"},
	PreRunE: NotInit,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Errorf("rmdir: missing arguments")
			return
		}
		for _, arg := range args {
			err := Client.RemoveDir(arg)
			if err != nil {
				log.Error(err)
				continue
			}
			log.Infof("rmdir %s success\n", arg)
		}
	},
}
