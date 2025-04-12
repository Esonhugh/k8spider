package nfs

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var lsOpts struct {
	ListAll bool
}

func init() {
	ListCommand.Flags().BoolVarP(&lsOpts.ListAll, "all", "a", false, "List all attrs")
	NFSCommand.AddCommand(ListCommand)
}

var ListCommand = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list nfs directory",
	Long:    "list nfs directory",
	PreRunE: NotInit,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			args = []string{"."}
		}
		for _, arg := range args {
			entries, err := Client.ListDir(arg)
			if err != nil {
				log.Error(err)
				continue
			}

			for _, entry := range entries {
				if lsOpts.ListAll {
					modeStr := entry.Mode().Perm().String()
					if entry.IsDir() {
						modeStr = strings.Replace(modeStr, "-", "d", 1)
					}
					if entry.Attr.IsSet {
						fmt.Printf("%s\t%d\t%d\t%d\t%s\t%s\n",
							modeStr, entry.Attr.Attr.UID, entry.Attr.Attr.GID,
							entry.Size(), entry.ModTime().Format("2006-01-02 15:04:05"), entry.Name())
					} else {

						fmt.Printf("%s\t%d\t%s\t%s\n",
							modeStr, entry.Size(), entry.ModTime().Format("2006-01-02 15:04:05"), entry.Name())
					}
				} else {
					fmt.Printf("%s\n", entry.Name())
				}
			}
		}
	},
}
