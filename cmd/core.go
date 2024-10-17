package morsebot

import (
	"os"
	"errors"
	"fmt"

	"github.com/wtbui/MorseBot/pkg/options"
	"github.com/wtbui/MorseBot/pkg/events"
	data "github.com/wtbui/MorseBot/pkg/data"
	discordgo "github.com/bwmarrin/discordgo"
)

// Initialize Morse Bot + Logger (TODO)
func Start(opts *options.Options) (int, error) {
	// Register any new govee keys to database
	if len(opts.Register) > 0 {
		fmt.Println("Registering new govee keys into DB")
		err := data.RegisterGKey(opts.Register)
		if err != nil {
			return ExitError, err
		}

		return ExitOk, nil
	}

	// Start up bot using discordgo package
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
