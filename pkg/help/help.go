package help

import (
	discordgo "github.com/bwmarrin/discordgo"	
)

var commandMap = map[string]CommandFunc {
	"echo": RunEcho,
	"help": RunHelp,
}

func RunEcho(s *discordgo.Session, m *discordgo.MessageCreate) {

}

func Run
