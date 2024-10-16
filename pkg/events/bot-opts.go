package events

import (
	"strings"

	discordgo "github.com/bwmarrin/discordgo"
	echo "github.com/wtbui/MorseBot/pkg/echo"
)

type CommandFunc func(botOptions BotOptions)

var commandMap = map[string]CommandFunc {
	"echo": RunEcho,
	"help": RunHelp,
}

type BotOptions struct {
	Command CommandFunc
	Username string
}

var BotOptionsDefaults = BotOptions{commandMap["help"], "", true, false}

func ParseOptions(args string) (*BotOptions, error) {
	var botOptions = &BotOptionsDefaults
	
	content := strings.TrimPrefix(args, CommandPrefix)
	parts := strings.Fields(content)
	
	command = parts[0]
	args = parts[1:]

	if cmdFunc, ok := commandMap[command]; ok {
		botOptions.Command = cmdFunc
	} else {
		return nil, errors.New("No command specified")
	} 

	for index, opt := range args {
		if strings.HasPrefix(opt, "@") {
			botOptions.Username = strings.TrimPrefix(opt, "@")
		}
	}

	return botOptions, nil
}


