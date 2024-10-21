package morsebot

import (
	"os"
	"errors"

	"github.com/wtbui/MorseBot/pkg/options"
	"github.com/wtbui/MorseBot/pkg/events"
	data "github.com/wtbui/MorseBot/pkg/data"
	discordgo "github.com/bwmarrin/discordgo"
	logger "github.com/wtbui/MorseBot/pkg/logging"
  "go.uber.org/zap"
)

// Initialize Morse Bot + Logger 
func Start(opts *options.Options) (int, error) {
	logger.InitLogger()

	// Register any new govee keys to database
	if len(opts.RegisterGKey) > 0 {
		zap.S().Info("Registering new govee keys into DB")
		err := data.RegisterUser(opts.RegisterGKey)
		if err != nil {
			return ExitError, err
		}

		return ExitOk, nil
	}

	if len(opts.DeleteGKey) > 0 {
		zap.S().Info("Deleting user from govee key DB")
		err := data.DeleteUser(opts.DeleteGKey)
		if err != nil {
			return ExitError, err
		}

		return ExitOk, nil
	}

	// Start up bot using discordgo package
	zap.S().Info("Starting MorseBot...")
	var token = ""

	if len(opts.APIKey) == 0 {
		token = os.Getenv("MORSEBOT")

		if len(token) == 0 {
			return ExitError, errors.New("Missing Discord API Key")
		}
	}

	zap.S().Info("Found API Key: " + token)
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
