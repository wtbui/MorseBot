package utils

import (
	"strings"
)

type BotOptions struct {
	Command string	
	Username string
	Aux []string
}

var BotOptionsDefaults = BotOptions{"help", "", []string{}}

func ParseOptions(message string, commandPrefix string) (*BotOptions, error) {
	var botOptions = &BotOptionsDefaults
	
	content := strings.TrimPrefix(message, commandPrefix)
	parts := strings.Fields(content)
	
	command := parts[0]
	args := parts[1:]

	botOptions.Command = command
	
	for _, opt := range args {
		// Use switch statement if more options added
		if strings.HasPrefix(opt, "@") {
			botOptions.Username = strings.TrimPrefix(opt, "@")
		} else {
			botOptions.Aux = append(botOptions.Aux, opt)
		}
	}

	return botOptions, nil
}
