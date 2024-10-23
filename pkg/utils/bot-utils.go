package utils

import (
	discordgo "github.com/bwmarrin/discordgo"
)
type JobReport struct {
	Job string
	Status bool
}

func GenerateReportEmbed(s *discordgo.Session, cid string, report JobReport) {
	eb := NewEmbed()
	eb.SetColor(0x7BCDBA)
	eb.SetThumbnail(s.State.User.AvatarURL(""))

	if report.Status {
		eb.AddField("MorseBot Report", report.Job + " successful", false)
	} else {
		eb.AddField("MorseBot Report", report.Job + " unsuccessful", false)
	}

	s.ChannelMessageSendEmbed(cid, eb.MessageEmbed)
}
