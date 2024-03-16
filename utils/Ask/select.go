package Ask

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/esonhugh/go-cli-template-v2/utils/Error"
)

type HasUniqueIDCanSelectedObject interface {
	GenUniqueID() string
}

// OneOf will ask user to select one of the options
func OneOf[T HasUniqueIDCanSelectedObject](options []T, msg string, opts ...string) (s T) {
	optionsStringArr := []string{}
	for _, k := range options {
		optionsStringArr = append(optionsStringArr, k.GenUniqueID())
	}
	prompt := &survey.Select{
		Message: msg,
		Options: optionsStringArr,
	}
	if len(opts) > 0 {
		prompt.Default = opts[0]
	}
	if len(opts) > 1 {
		prompt.Help = opts[1]
	}
	ResultString := ""
	Error.HandleFatal(survey.AskOne(prompt, &ResultString))
	for _, v := range options {
		if v.GenUniqueID() == ResultString {
			s = v
		}
	}
	return
}

// SomeOf will ask user to select some of the options
func SomeOf[T HasUniqueIDCanSelectedObject](options []T, msg string) (s []T) {
	optionsStringArr := []string{}
	for _, k := range options {
		optionsStringArr = append(optionsStringArr, k.GenUniqueID())
	}
	prompt := &survey.MultiSelect{
		Message: msg,
		Options: optionsStringArr,
	}
	ResultString := []string{}
	Error.HandleFatal(survey.AskOne(prompt, &ResultString))
	for _, v := range ResultString {
		for _, k := range options {
			if k.GenUniqueID() == v {
				s = append(s, k)
			}
		}
	}
	return
}
