package events

import (
	"strings"

	discordgo "github.com/bwmarrin/discordgo"
)

var (
	CommandPrefix = "!"
)

func RegisterHandlers(s *discordgo.Session) (err error){
	// Todo: Use a map to register event handlers instead of static
	s.AddHandler(ready)
	s.AddHandler(messageCreate)

	return nil
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	s.UpdateGameStatus(0, "That's Golf!")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Avoid reading self created messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, CommandPrefix) {
		// Handle New Command
		return
	}
}
