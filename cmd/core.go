package morsebot

import (
	"os"

	"github.com/wtbui/MorseBot/pkg/options"
	"github.com/wtbui/MorseBot/pkg/events"
	discordgo "github.com/bwmarrin/discordgo"
)

// Initialize Morse Bot + Logger (TODO)
func Start(opts *options.Options) (int, error) {
	var token = ""

	if len(opts.APIKey) == 0 {
		token = os.Getenv("MORSEBOT")

		if token == "" {
			// TODO setup logger
		}
	}

	mb_sess, err := discordgo.New("MorseBot " + token)
	if err != nil {
		return ExitError, err
	}
	
	err = events.RegisterHandlers(mb_sess)
	if err != nil {
		return ExitError, err
	}

	return ExitOk, nil
}
