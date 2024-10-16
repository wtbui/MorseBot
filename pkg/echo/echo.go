package echo

import (
	discordgo "github.com/bwmarrin/discordgo"
	utils "github.com/wtbui/MorseBot/pkg/utils"
)


func RunEcho(s *discordgo.Session, cid string, botOpts *utils.BotOptions) error {
	s.ChannelMessageSend(cid, "echo")

	return nil
}


