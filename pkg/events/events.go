package events

import (
	"strings"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	discordgo "github.com/bwmarrin/discordgo"
)

var (
	CommandPrefix = "!"
)

func InitBot(s *discordgo.Session) (err error) {
	s.Open()

	err = registerHandlers(s)
	if err != nil {
		return err
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Morse bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("Closing Morse Bot...")
	s.Close()

	return nil
}

func registerHandlers(s *discordgo.Session) (err error) {
	// TODO: Use a map to register event handlers instead of static
	s.AddHandler(ready)
	s.AddHandler(messageCreate)

	return nil
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	s.UpdateGameStatus(1, "That's Golf!")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Message Found in " + m.ChannelID + " from " + m.Author.ID)

	// Avoid reading self created messages
	if m.Author.ID == s.State.User.ID {
		return
	}
	
	if strings.HasPrefix(m.Content, CommandPrefix) {
		// Handle New Command
		opts, err := ParseOptions(m.Content)
		if err != nil {
			fmt.Println("Failure to parse command options")
			return
		}

		err = opts.Command(opts)
		if err != nil {
			fmt.Println("Failure to execute command")
		}
	}
}
