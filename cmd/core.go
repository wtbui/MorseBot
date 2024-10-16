package morsebot

import (
	"os"
	"errors"
	"fmt"

	"github.com/wtbui/MorseBot/pkg/options"
	"github.com/wtbui/MorseBot/pkg/events"
	discordgo "github.com/bwmarrin/discordgo"
)

// Initialize Morse Bot + Logger (TODO)
func Start(opts *options.Options) (int, error) {
	fmt.Println("Starting MorseBot...")

	var token = ""

	if len(opts.APIKey) == 0 {
		token = os.Getenv("MORSEBOT")

		if len(token) == 0 {
			return ExitError, errors.New("Missing Discord API Key")
		}

		if token == "" {
			// TODO setup logger
		}
	}

	fmt.Println("Found API Key: " + token)
	mb_sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return ExitError, err
	}
	
	err = events.InitBot(mb_sess)
	if err != nil {
		return ExitError, err
	}
	
	return ExitOk, nil
}
