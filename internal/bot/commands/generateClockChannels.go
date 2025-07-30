package commands

import (
	"fmt"

	"github.com/Sush1sui/cico-bot-go/internal/config"
	"github.com/Sush1sui/cico-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func GenerateClockChannels(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" { return }

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Setting up clock channels, please wait...",
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})

	res, err := repository.ClockChannelService.GetAllClockChannelInterface()
	if err != nil {
		mess := "Error retrieving clock channels."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		fmt.Println(err)
		return
	}
	if len(res) >= 1 {
		mess := "Clock channels already exist."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		fmt.Println(err)
		return
	}

	category, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name: "⏰ Time Tracking",
		Type: discordgo.ChannelTypeGuildCategory,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   i.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel,
			},
		},
	})
	if err != nil {
		mess := "Error creating category."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		fmt.Println(err)
		return
	}

	clockInChannel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name: "🟢-clock-in",
		Type: discordgo.ChannelTypeGuildText,
		ParentID: category.ID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   i.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel,
			},
			{
				ID:   config.GlobalConfig.TL_ROLE_ID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory,
				Deny: discordgo.PermissionSendMessages | discordgo.PermissionAddReactions | discordgo.PermissionUseExternalEmojis | discordgo.PermissionAttachFiles | discordgo.PermissionEmbedLinks | discordgo.PermissionManageMessages | discordgo.PermissionManageThreads | discordgo.PermissionCreatePrivateThreads | discordgo.PermissionCreatePublicThreads | discordgo.PermissionSendMessagesInThreads,
			},
			{
				ID:   config.GlobalConfig.CHATTER_ROLE_ID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory,
				Deny: discordgo.PermissionSendMessages | discordgo.PermissionAddReactions | discordgo.PermissionUseExternalEmojis | discordgo.PermissionAttachFiles | discordgo.PermissionEmbedLinks | discordgo.PermissionManageMessages | discordgo.PermissionManageThreads | discordgo.PermissionCreatePrivateThreads | discordgo.PermissionCreatePublicThreads | discordgo.PermissionSendMessagesInThreads,
			},
		},
	})
	if err != nil {
		mess := "Error creating clock-in channel."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		fmt.Println(err)
		return
	}


	clockOutChannel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name: "🔴-clock-out",
		Type: discordgo.ChannelTypeGuildText,
		ParentID: category.ID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   i.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel,
			},
			{
				ID:   config.GlobalConfig.TL_ROLE_ID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory,
				Deny: discordgo.PermissionSendMessages | discordgo.PermissionAddReactions | discordgo.PermissionUseExternalEmojis | discordgo.PermissionAttachFiles | discordgo.PermissionEmbedLinks | discordgo.PermissionManageMessages | discordgo.PermissionManageThreads | discordgo.PermissionCreatePrivateThreads | discordgo.PermissionCreatePublicThreads | discordgo.PermissionSendMessagesInThreads,
			},
			{
				ID:   config.GlobalConfig.CHATTER_ROLE_ID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory,
				Deny: discordgo.PermissionSendMessages | discordgo.PermissionAddReactions | discordgo.PermissionUseExternalEmojis | discordgo.PermissionAttachFiles | discordgo.PermissionEmbedLinks | discordgo.PermissionManageMessages | discordgo.PermissionManageThreads | discordgo.PermissionCreatePrivateThreads | discordgo.PermissionCreatePublicThreads | discordgo.PermissionSendMessagesInThreads,
			},
		},
	})
	if err != nil {
		mess := "Error creating clock-out channel."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		fmt.Println(err)
		return
	}

	adminChannel, err := s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
		Name: "⚙️-time-admin",
		Type: discordgo.ChannelTypeGuildText,
		ParentID: category.ID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:   i.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionViewChannel,
			},
		},
	})
	if err != nil {
		mess := "Error creating admin channel."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		fmt.Println(err)
		return
	}

	greenColor := 0x00FF00 // Green color for clock in role
	// create clock in role
	clockInRole, err := s.GuildRoleCreate(i.GuildID, &discordgo.RoleParams{
		Name: "Clocked In",
		Color: &greenColor,
	})
	if err != nil {
		mess := "Error creating clock-in role."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		fmt.Println(err)
		return
	}

	clockInInterface, err := s.ChannelMessageSendComplex(clockInChannel.ID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color: greenColor,
			Title: "🟢 Clock In System",
			Description: "**Ready to start your shift?**\n\nClick the button below to clock in and begin tracking your work time.",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name: "📋 What happens when you clock in:",
					Value: "• You'll receive the **Clocked In** role\n• Your start time will be recorded\n• You'll get a confirmation message",
					Inline: false,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Time Tracking Bot • Click the button to get started",
				IconURL: "https://cdn.discordapp.com/emojis/1234567890.png",
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.Button{
						Label: "🟢 Clock In",
						Style: discordgo.SuccessButton,
						CustomID: "clock_in",
						Emoji: &discordgo.ComponentEmoji{Name: "⏰"},
					},
				},
			},
		},
	})
	if err != nil {
		mess := "Error creating clock-in interface."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		fmt.Println(err)
		return
	}

	clockOutInterface, err := s.ChannelMessageSendComplex(clockOutChannel.ID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color: 0xFF0000, // Red color for clock out role
			Title: "🔴 Clock Out System",
			Description: "**Finishing your shift?**\n\nClick the button below to clock out and stop tracking your work time.",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name: "📋 What happens when you clock out:",
					Value: "• The **Clocked In** role will be removed\n• Your end time will be recorded\n• Total work time will be calculated",
					Inline: false,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Time Tracking Bot • Click the button to finish your shift",
				IconURL: "https://cdn.discordapp.com/emojis/1234567890.png",
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					&discordgo.Button{
						Label: "🔴 Clock Out",
						Style: discordgo.DangerButton,
						CustomID: "clock_out",
						Emoji: &discordgo.ComponentEmoji{Name: "⏹️"},
					},
				},
			},
		},
	})
	if err != nil {
		mess := "Error creating clock-out interface."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		fmt.Println(err)
		return
	}

	clockChannelSetup, err := repository.ClockChannelService.CreateClockChannelInterface(category.ID, clockInChannel.ID,clockInInterface.ID, clockOutChannel.ID, clockOutInterface.ID, adminChannel.ID, clockInRole.ID)
	if err != nil || clockChannelSetup == nil {
		mess := "Error saving clock channel setup in the database."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		fmt.Println(err)
		return
	}

	config.GlobalConfig.ClockInRoleID = clockChannelSetup.ClockInRoleID
	config.GlobalConfig.AdminChannelID = clockChannelSetup.AdminChannelID

	mess := "Clock channels setup completed successfully."
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
}