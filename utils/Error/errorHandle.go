package Error

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

const ERRORMSG = "Cli Runtime Error:"

// HandleError is a function to handle error using Errorln, do nothing if error is nil
func HandleError(e error, attachedMsg ...string) {
	if e != nil {
		if len(attachedMsg) > 0 {
			log.Errorln(ERRORMSG, strings.Join(attachedMsg, " "), e)
		} else {
			log.Errorln(ERRORMSG, e)
		}
	}
}

// HandlePanic is a function to handle error using Panicln, do nothing if error is nil
func HandlePanic(e error, attachedMsg ...string) {
	if e != nil {
		if len(attachedMsg) > 0 {
			log.Panicln(ERRORMSG, strings.Join(attachedMsg, " "), e)
		} else {
			log.Panicln(ERRORMSG, e)
		}
	}
}

// HandleFatal is a function to handle error using Fatalln, do nothing if error is nil
func HandleFatal(e error, attachedMsg ...string) {
	if e != nil {
		if len(attachedMsg) > 0 {
			log.Fatalln(ERRORMSG, strings.Join(attachedMsg, " "), e)
		} else {
			log.Fatalln(ERRORMSG, e)
		}
	}
}

func HandleWarn(e error, attachedMsg ...string) {
	if e != nil {
		if len(attachedMsg) > 0 {
			log.Warnln(ERRORMSG, strings.Join(attachedMsg, " "), e)
		} else {
			log.Warnln(ERRORMSG, e)
		}
	}
}
