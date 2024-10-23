package lightsync

import (
	"strconv"
	"strings"

	discordgo "github.com/bwmarrin/discordgo"
	utils "github.com/wtbui/MorseBot/pkg/utils"
  goveego "github.com/wtbui/MorseBot/pkg/goveego"
	data "github.com/wtbui/MorseBot/pkg/data"
	"go.uber.org/zap"
)

type LSyncJob struct {
	Users []string
	Toggle string
	Color int
}
 
func RunLightsync(s *discordgo.Session, cid string, botOpts *utils.BotOptions) error {
	lJob, err := parseOptions(botOpts)
	if err != nil {
		return err
	}
	
	zap.S().Debug("Starting lights job")
	err = runLightsJob(lJob)
	jobReport := utils.JobReport{"Lights Change", true} 
	if err != nil {
		jobReport.Status = false
		utils.GenerateReportEmbed(s, cid, jobReport)
		
		return err
	}
	
	utils.GenerateReportEmbed(s, cid, jobReport)
	return nil
}

func runLightsJob(lJob *LSyncJob) error {
	for _, user := range lJob.Users {
		gclient, err := goveego.Init(user)
		if err != nil {
			return err
		}
		
		if lJob.Toggle == "on" {
			gclient.TurnOnOffAll(1)
		} else {
			gclient.TurnOnOffAll(0)
		}

		if lJob.Color > -1 {
			gclient.ChangeColorAll(lJob.Color)
		}
	}

	return nil
}

func parseOptions(botOpts *utils.BotOptions) (*LSyncJob, error) {
	newJob := &LSyncJob{[]string{}, "on", -1}

	// Fetch from Database
	goveeDb, err := data.RetrieveCurrentGDB()	
	if err != nil {
		return nil, err
	}

	zap.S().Debug("Length of recieved database")
	zap.S().Debug(len(goveeDb))

	if len(botOpts.Username) > 0 {
		newJob.Users = append(newJob.Users, goveeDb[botOpts.Username].GKey)
		zap.S().Debug("Appended " + goveeDb[botOpts.Username].GKey)
	} else {
		newJob.Users = append(newJob.Users, goveeDb[botOpts.Sender].GKey)	
	}

	for _, opt := range botOpts.Aux {
		if opt == "all" {
			newJob.Users = []string{}
			for _, reg := range goveeDb {
				newJob.Users = append(newJob.Users, reg.GKey)
			}
		}

		if color, ok := LColors[strings.ToLower(opt)]; ok {
			newJob.Color = color
		}

		if utils.IsNumeric(opt) {
			newJob.Color, _ = strconv.Atoi(opt)
		}

		if strings.ToLower(opt) == "off" {
			newJob.Toggle = strings.ToLower(opt)
		}
	}	

	return newJob, nil
}

