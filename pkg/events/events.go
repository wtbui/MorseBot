package events

import (
	"strings"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	discordgo "github.com/bwmarrin/discordgo"
	help "github.com/wtbui/MorseBot/pkg/help"
	utils "github.com/wtbui/MorseBot/pkg/utils"
)

var (
	CommandPrefix = "!"
)

type CommandFunc func(s *discordgo.Session, botOpts *utils.BotOptions) error
var commandMap = map[string]CommandFunc {
	"echo": help.RunEcho, "help": help.RunHelp,
}

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
	s.UpdateGameStatus(1, "That's Golf!")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Message Found in " + m.ChannelID + " from " + m.Author.ID)

	if m.Author.ID == s.State.User.ID {
		return
	}
	
	if strings.HasPrefix(m.Content, CommandPrefix) {
		opts, err := utils.ParseOptions(m.Content, CommandPrefix)
		if err != nil {
			fmt.Println("Failure to parse command options")
			return
		}

		if cmdFunc, ok := commandMap[opts.Command]; ok {
			err = cmdFunc(s, opts)
		} else {
			fmt.Println("Unrecognized Command")
			return
		} 
	
		if err != nil {
			fmt.Println("Failure to execute command")
		}
	}
}
