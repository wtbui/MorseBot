package lightsync

import (
	discordgo "github.com/bwmarrin/discordgo"
	utils "github.com/wtbui/MorseBot/pkg/utils"
//	goveego "github.com/wtbui/MorseBot/pkg/goveego"
)

type LSyncJob struct {
	Users []string
	Toggle string
	Color int
}
 
func RunLightsync(s *discordgo.Session, cid string, botOpts *utils.BotOptions) error {
	s.ChannelMessageSend(cid, "Lights")

	return nil
}

func parseOptions(botOpts *utils.BotOptions) (*LSyncJob, error) {
	newJob := &LSyncJob{[]string{}, "", -1}
	
	if len(botOpts.Username) == 0 {
		newJob.Users = append(newJob.Users, botOpts.Username)
	} else {
		
	}

	for _, opt := range botOpts.Aux {
		if opt == "all" {
			
		}

		if opt == "" {

		}
	}	

	return newJob, nil
}


