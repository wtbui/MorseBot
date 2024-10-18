package events

import (
	"strings"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	discordgo "github.com/bwmarrin/discordgo"
	echo "github.com/wtbui/MorseBot/pkg/echo"
	utils "github.com/wtbui/MorseBot/pkg/utils"
	lsyn "github.com/wtbui/MorseBot/pkg/lightsync"
)

var (
	CommandPrefix = "#"
)

// Define entry points into commands here
type CommandFunc func(s *discordgo.Session, cid string, botOpts *utils.BotOptions) error
var commandMap = map[string]Command {
	"help": Command{"help", echo.RunEcho, "Displays this help message"},
	"echo": Command{"echo", echo.RunEcho, "Echoes back a message in the same channel"},
	"lights": Command{"lights", lsyn.RunLightsync, "Adjusts lights"},
}

type Command struct {
	Name string
	CmdFunc CommandFunc
	Descrip string
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

// Ready Event Handler
func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(1, "That's Golf!")
}

// Message Create Event Handler
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("Message Found in " + m.ChannelID + " from " + m.Author.ID)

	if m.Author.ID == s.State.User.ID {
		return
	}
	
	if strings.HasPrefix(m.Content, CommandPrefix) {
		opts, err := utils.ParseOptions(m.Content, CommandPrefix)
		opts.Sender = m.Author.ID

		if err != nil {
			fmt.Println("Failure to parse command options")
			return
		}

		if opts.Command == "help" {
			runHelp(s, m.ChannelID)
			return
		}

		if cmdFunc, ok := commandMap[opts.Command]; ok {
			err = cmdFunc.CmdFunc(s, m.ChannelID, opts)
		} else {
			fmt.Println("Unrecognized Command")
			return
		} 
	
		if err != nil {
			fmt.Println("Failure to execute command")
		}
	}
}

func runHelp(s *discordgo.Session, cid string) {
	eb := utils.NewEmbed()
	eb.SetTitle("Morse Bot Help Menu")
	eb.SetColor(0x61E294)
	eb.SetThumbnail(s.State.User.AvatarURL(""))
	for _, cmd := range commandMap {
		eb.AddField(cmd.Descrip, "```Usage: " + CommandPrefix + cmd.Name + "```", false)
	}

	s.ChannelMessageSendEmbed(cid, eb.MessageEmbed)
}
