package nfs

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	NFSCommand.AddCommand(CatCommand)
}

var CatCommand = &cobra.Command{
	Use:     "cat",
	Short:   "cat nfs file",
	Long:    "cat nfs file",
	PreRunE: NotInit,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}
		for _, arg := range args {
			content, err := Client.Cat(arg)
			if err != nil {
				log.Error(err)
				continue
			}
			fmt.Println(string(content))
		}
	},
}
