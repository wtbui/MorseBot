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
	Off bool
	Color int
	Temp int
	EffectId int
	EffectParamId int
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
		
		if lJob.Off {
			err = gclient.TurnOnOffAll(0)
			return err	
		} else {
			err = gclient.TurnOnOffAll(1)
		}

		if err != nil {
			return err
		}
		
		if lJob.EffectId > -1 {
			err = gclient.ChangeEffectAll(lJob.EffectParamId, lJob.EffectId)
			return err
		}

		if lJob.Temp > -1 {
			err = gclient.ChangeTempAll(lJob.Temp)
			return err
		}

		if lJob.Color > -1 {
			err = gclient.ChangeColorAll(lJob.Color)
			return err
		}
	}

	return nil
}

func parseOptions(botOpts *utils.BotOptions) (*LSyncJob, error) {
	newJob := &LSyncJob{[]string{}, false, -1, -1, -1, -1}

	// Fetch from Database
	goveeDb, err := data.RetrieveCurrentGDB()	
	if err != nil {
		return nil, err
	}

	zap.S().Debug("Length of recieved database")
	zap.S().Debug(len(goveeDb))

	if len(botOpts.Username) > 0 {
		if target, ok := goveeDb[botOpts.Username]; ok {
			newJob.Users = append(newJob.Users, target.GKey)
		}
	} else {
		if target, ok := goveeDb[botOpts.Sender]; ok {
			newJob.Users = append(newJob.Users, target.GKey)	
		}
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
			newJob.Off = true
		}

		if temp, ok := LTemps[strings.ToLower(opt)]; ok {
			newJob.Temp = temp
		}

		if effect, ok := LEffects[strings.ToLower(opt)]; ok {
			newJob.EffectId = effect.Id
			newJob.EffectParamId = effect.ParamId
		}
	}	

	return newJob, nil
}

