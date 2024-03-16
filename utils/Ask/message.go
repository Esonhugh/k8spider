package Ask

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
)

// ForString : Ask user for a string
// msg: The message to show to the user
// opts: Optional, opts[0] is default value and opts[1] is help message
func ForString(msg string, opts ...string) (str string) {
	prompt := &survey.Input{
		Message: msg,
	}
	if len(opts) > 0 {
		prompt.Default = opts[0]
	}
	if len(opts) > 1 {
		prompt.Help = opts[1]
	}
	err := survey.AskOne(prompt, &str)
	if err != nil {
		log.Debugf("Ask User \"%v\" input error: %v", msg, err)
		log.Fatal("User input error. Exiting...")
		// Error.HandleFatal(err)
	}
	return str
}

// ForPath will ask user for a valid path.
// If opts[0] is set, it will be used as the default beginning path
// Else it will use the CFP Home path as default beginning path.
func ForPath(msg string, opts ...string) (str string) {
	prompt := &survey.Input{
		Message: msg,
		Help: `Input the path of the file or directory and press Tab.
If the path is absolute, it will be used directly,
If the path is not absolute, it will be joined with the "" path, 
If it begins with ~/ it will be joined with the user's home path,
If it begins with ./ it will be joined with the current working directory.
`,
		Suggest: func(toComplete string) []string {
			ret := []string{toComplete}
			var initPath string
			if strings.HasPrefix(toComplete, "/") {
				// Start with /
				initPath = toComplete
			} else if strings.HasPrefix(toComplete, "~/") {
				// Start with ~/
				hp, err := os.UserHomeDir()
				if err != nil {
					return ret
				}
				initPath = filepath.Join(hp, toComplete)
			} else if strings.HasPrefix(toComplete, "./") {
				// Start with ./
				if path, err := os.Getwd(); err != nil {
					return ret
				} else {
					initPath = filepath.Dir(path)
				}
			} else if len(opts) > 0 {
				// Start with opts[0]
				// Module Developer Can use this opt to set a default path
				initPath = filepath.Join(opts[0], toComplete)
			} else {
				// Start with current
				initPath = filepath.Join("", toComplete)
			}

			// List Dir
			files, err := os.ReadDir(initPath)
			if err != nil {
				return ret
			}
			for _, f := range files {
				path := filepath.Join(initPath, f.Name())
				ret = append(ret, path)
			}
			return ret
		},
	}
	if len(opts) > 0 {
		prompt.Default = opts[0]
	}
	err := survey.AskOne(prompt, &str)
	if err != nil {
		log.Debugf("Ask User \"%v\" input error: %v", msg, err)
		log.Fatal("User input error. Exiting...")
		// Error.HandleFatal(err)
	}
	return str
}
