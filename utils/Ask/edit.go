package Ask

import (
	"encoding/json"

	"github.com/AlecAivazis/survey/v2"
	"github.com/esonhugh/go-cli-template-v2/utils/Error"
)

func Editor[T any](object T, msg string) (result T) {
	data, err := json.MarshalIndent(object, "", "  ")
	Error.HandleFatal(err, "error in json marshal.")
	prompt := &survey.Editor{
		Message:       msg,
		Default:       string(data),
		AppendDefault: true,
	}
	var Answer string
	Error.HandleFatal(survey.AskOne(prompt, &Answer), "error in user input.")
	Error.HandleFatal(json.Unmarshal([]byte(Answer), &result), "you could enter an invalid json.")
	return result
}
