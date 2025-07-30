package commands

import (
	"fmt"

	"github.com/Sush1sui/cico-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func DeleteGeneratedChannels(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" { return }

	clockChannelSetups, err := repository.ClockChannelService.GetAllClockChannelInterface()
	if err != nil || len(clockChannelSetups) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error retrieving clock channel setup from the database.",
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	
	for _, setup := range clockChannelSetups {
		_, err = s.ChannelDelete(setup.ClockInChannelID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error deleting clock in channel.",
					Flags:  discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		_, err = s.ChannelDelete(setup.ClockOutChannelID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error deleting clock out channel.",
					Flags:  discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		_, err = s.ChannelDelete(setup.AdminChannelID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error deleting admin channel.",
					Flags:  discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		_, err = s.ChannelDelete(setup.CategoryID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error deleting category channel.",
					Flags:  discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		err = s.GuildRoleDelete(i.GuildID, setup.ClockInRoleID)
		if err != nil { fmt.Println("Error deleting clock in role:", err) }
	}

	err = repository.ClockChannelService.DeleteAllClockChannelInterface()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error deleting clock channels setup in the database.",
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Successfully deleted all generated channels.",
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}