package lightsync

import (
	"strconv"
	"strings"
	"context"
	"golang.org/x/sync/errgroup"
	"time"

	discordgo "github.com/bwmarrin/discordgo"
	utils "github.com/wtbui/MorseBot/pkg/utils"
 	goveego "github.com/wtbui/MorseBot/pkg/goveego"
	data "github.com/wtbui/MorseBot/pkg/data"
	"go.uber.org/zap"
)

type LSyncJob struct {
	Off           bool
	Color         int
	Temp          int
	Brightness    int 
	EffectId      int
	EffectParamId int
}
 
func RunLightsync(s *discordgo.Session, cid string, botOpts *utils.BotOptions) utils.JobReport {
	start := time.Now()
	zap.S().Debug("Starting lights job")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobReport := utils.JobReport{"Lights Change", false, true, nil}

	users, lJob, err := parseOptions(botOpts)
	if err != nil {
		jobReport.E = err
		return jobReport
	}
	
	group, ctx := errgroup.WithContext(ctx)
	for _, user := range users {
		user := user
		group.Go(func() error {
			return runLightsJob(ctx, user, lJob)
		})
	}

	err = group.Wait()

	if err != nil {
		jobReport.E = err
	} else {
		jobReport.Status = true
	}
	
	elapsed := time.Since(start)
	zap.S().Debug("Execution Time: ", elapsed)
	return jobReport
}

func parseOptions(botOpts *utils.BotOptions) ([]string, *LSyncJob, error) {
	newJob := &LSyncJob{false, -1, -1, -1, -1, -1}
	users := []string{}

	// Fetch from Database
	goveeDb, err := data.RetrieveCurrentGDB()	
	if err != nil {
		return users, nil, err
	}

	zap.S().Debug("Length of recieved database")
	zap.S().Debug(len(goveeDb))

	if len(botOpts.Username) > 0 {
		if target, ok := goveeDb[botOpts.Username]; ok {
			users = append(users, target.GKey)
		}
	} else {
		if target, ok := goveeDb[botOpts.Sender]; ok {
			users = append(users, target.GKey)	
		}
	}
	
	// Parse light job arguments 
	for _, opt := range botOpts.Aux {
		if opt == "all" {
			users = []string{}
			for _, reg := range goveeDb {
				users = append(users, reg.GKey)
			}
		}

		if color, ok := LColors[strings.ToLower(opt)]; ok {
			newJob.Color = color
		}

		if utils.IsNumeric(opt) {
			newJob.Brightness, _ = strconv.Atoi(opt)
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

	return users, newJob, nil
}

func runLightsJob(ctx context.Context, user string, lJob *LSyncJob) error {
	gclient, err := goveego.NewClient(user)
	if err != nil { 
		return err 
	}
	
	if lJob.Off {
		err = gclient.ChangeLightAll(ctx, goveego.OFF, []int{0}) 
		return err
	} else {
		err = gclient.ChangeLightAll(ctx, goveego.ON, []int{1}) 
		if err != nil { 
			return err 
		}
	}

	if lJob.Brightness > -1 {
		err = gclient.ChangeLightAll(ctx, goveego.BRIGHT, []int{lJob.Brightness})
		if err != nil {
			return err
		} 
	}
	
	if lJob.EffectId > -1 {
		err = gclient.ChangeLightAll(ctx, goveego.EFFECT, []int{lJob.EffectId, lJob.EffectParamId})
		return err
	}

	if lJob.Temp > -1 {
		err = gclient.ChangeLightAll(ctx, goveego.TEMP, []int{lJob.Temp}) 
		return err
	}

	if lJob.Color > -1 {
		err = gclient.ChangeLightAll(ctx, goveego.COLOR, []int{lJob.Color}) 
		return err
	}

	return nil
}


