package nfs

import (
	"errors"

	cmd "github.com/esonhugh/k8spider/cmd"
	"github.com/esonhugh/k8spider/pkg/nfs"
	"github.com/spf13/cobra"
)

var Opts struct {
	Server   string
	BasePath string

	UID      int
	GID      int
	Hostname string
}

var Client *nfs.NFSClient

func init() {
	NFSCommand.PersistentFlags().StringVarP(&Opts.Server, "server", "s", "", "nfs server")
	NFSCommand.PersistentFlags().StringVarP(&Opts.BasePath, "base-path", "b", "/", "base mount path")
	NFSCommand.PersistentFlags().IntVarP(&Opts.UID, "uid", "u", 0, "uid")
	NFSCommand.PersistentFlags().IntVarP(&Opts.GID, "gid", "g", 0, "gid")
	NFSCommand.PersistentFlags().StringVarP(&Opts.Hostname, "hostname", "H", "localhost", "hostname")
	cmd.RootCmd.AddCommand(NFSCommand)
}

var NFSCommand = &cobra.Command{
	Use:   "nfs",
	Short: "nfs is a nfs client command",
	Long:  "nfs is a nfs client command",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		Client, err = nfs.NewNFSClient(Opts.Server, Opts.BasePath, nfs.CustomAuth(Opts.Hostname, Opts.UID, Opts.GID))
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func NotInit(cmd *cobra.Command, args []string) error {
	if Client == nil {
		return errors.New("nfs client not initialized")
	}
	return nil
}
