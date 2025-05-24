package echo

import (
	discordgo "github.com/bwmarrin/discordgo"
	utils "github.com/wtbui/MorseBot/pkg/utils"
)


func RunEcho(s *discordgo.Session, cid string, botOpts *utils.BotOptions) utils.JobReport {
	s.ChannelMessageSend(cid, "echo")

	return utils.JobReport{"test", true, false, nil}
}


