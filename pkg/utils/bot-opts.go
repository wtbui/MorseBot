package utils

import (
	"strings"
)

type BotOptions struct {
	Command string	
	Sender string
	Username string
	Aux []string
}

func botOptsDefaults() BotOptions {
	BotOptionsDefaults := BotOptions{"help", "", "", []string{}}
	return BotOptionsDefaults
}

func ParseOptions(message string, commandPrefix string) (*BotOptions, error) {
	var botOption = botOptsDefaults()
	botOptions := &botOption

	content := strings.TrimPrefix(message, commandPrefix)
	parts := strings.Fields(content)
	
	command := parts[0]
	args := parts[1:]

	botOptions.Command = command
	
	for _, opt := range args {
		// Use switch statement if more options added
		if strings.HasPrefix(opt, "<") {
			botOptions.Username = strings.Trim(opt, "<")
			botOptions.Username = strings.Trim(botOptions.Username, "@")
			botOptions.Username = strings.Trim(botOptions.Username, ">")
		} else {
			botOptions.Aux = append(botOptions.Aux, opt)
		}
	}


	return botOptions, nil
}
