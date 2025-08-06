package commands

import (
	"github.com/Sush1sui/cico-bot-go/internal/common"
	"github.com/Sush1sui/cico-bot-go/internal/config"
	"github.com/bwmarrin/discordgo"
)

func ExportCSVWithResetDatabaseCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Exporting to CSV and resetting database...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if i.Member.User.ID != "982491279369830460" && i.Member.User.ID != "608646101712502825" {
		mess := "Only <@982491279369830460> and <@608646101712502825> can use this command."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		return
	}

	err := common.ExportToCSV_CLEAN_DATABASE(s)
	if err != nil {
		mess := "Failed to export to CSV and reset database: " + err.Error()
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		return
	}

	m := "CSV file exported successfully and database reset in <#" + config.GlobalConfig.AdminChannelID + ">!"
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &m})
}