package help

import (
	"fmt"

	discordgo "github.com/bwmarrin/discordgo"
	utils "github.com/wtbui/MorseBot/pkg/utils"
)


func RunEcho(s *discordgo.Session, botOpts *utils.BotOptions) error {
	fmt.Println("Running Echo")

	return nil
}

func RunHelp(s *discordgo.Session, botOpts *utils.BotOptions) error {
	fmt.Println("Running Help")

	return nil
}
