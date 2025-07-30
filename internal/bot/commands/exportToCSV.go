package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Sush1sui/cico-bot-go/internal/common"
	"github.com/Sush1sui/cico-bot-go/internal/config"
	"github.com/bwmarrin/discordgo"
)

func ExportCSVCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" { return }

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Exporting to CSV...",
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})

	filePath, err := common.ExportToCSV(s)
	if err != nil || filePath == "" {
		mess := "Failed to export to CSV: " + err.Error()
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		mess := "Failed to open CSV file: " + err.Error()
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		return
	}

	// create attachment from the csv file
	attachment := &discordgo.File{
		Name:        filepath.Base(filePath),
		ContentType: "text/csv",
		Reader: file,
	}

	s.ChannelMessageSendComplex(config.GlobalConfig.AdminChannelID, &discordgo.MessageSend{
		Files: []*discordgo.File{attachment},
	})

	file.Close()

	err = os.Remove(filePath)
	if err != nil {
		fmt.Printf("Failed to remove CSV file: %v\n", err)
	}

	m := "CSV file exported successfully in <#" + config.GlobalConfig.AdminChannelID + ">!"
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &m})
}