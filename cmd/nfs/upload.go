package nfs

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	NFSCommand.AddCommand(UploadCommand)
	NFSCommand.AddCommand(DownloadCommand)
}

var UploadCommand = &cobra.Command{
	Use:     "upload localfile remotepath",
	Aliases: []string{"put"},
	Short:   "upload file to nfs",
	Long:    "upload file to nfs",
	PreRunE: NotInit,
	Run: func(cmd *cobra.Command, args []string) {
		var localfile, remotepath string
		switch len(args) {
		case 2:
			remotepath = args[1]
			fallthrough
		case 1:
			localfile = filepath.Clean(args[0])
			if remotepath == "" {
				remotepath = localfile
			}
		default:
			log.Errorf("usage: nfs upload localfile [remotepath]")
		}
		if err := Client.Upload(localfile, remotepath); err != nil {
			log.Errorf("upload %s to %s failed: %v", localfile, remotepath, err)
			return
		}
		log.Infof("upload %s to %s success", localfile, remotepath)
	},
}

var DownloadCommand = &cobra.Command{
	Use:     "download remotefile localpath",
	Aliases: []string{"get"},
	Short:   "download file from nfs",
	Long:    "download file from nfs",
	PreRunE: NotInit,
	Run: func(cmd *cobra.Command, args []string) {
		var remotefile, localpath string
		switch len(args) {
		case 2:
			localpath = args[1]
			fallthrough
		case 1:
			remotefile = filepath.Clean(args[0])
			if localpath == "" {
				localpath = filepath.Base(remotefile)
			}
		default:
			log.Errorf("usage: nfs download remotefile [localpath]")
		}
		if err := Client.Download(remotefile, localpath); err != nil {
			log.Errorf("download %s to %s failed: %v", remotefile, localpath, err)
			return
		}
		log.Infof("download %s to %s success", remotefile, localpath)
	},
}
