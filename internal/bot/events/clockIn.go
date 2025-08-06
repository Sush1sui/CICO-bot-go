package events

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/Sush1sui/cico-bot-go/internal/config"
	"github.com/Sush1sui/cico-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func OnClockIn(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" { return }
	if i.Type != discordgo.InteractionMessageComponent { return }
	if i.MessageComponentData().CustomID != "clock_in" { return }

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Clocking in, please wait...",
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})

	if slices.Contains(i.Member.Roles, config.GlobalConfig.ClockInRoleID) {
		m := "Clock in failed, you are already clocked in."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &m})
		return
	}

	if !slices.Contains(i.Member.Roles, config.GlobalConfig.TL_ROLE_ID) && !slices.Contains(i.Member.Roles, config.GlobalConfig.CHATTER_ROLE_ID) {
		m := "You do not have Chatter or Team Leader role to clock in."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &m})
		return
	}

	// add clock in role to the user
	clockInTime := time.Now()
	err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, config.GlobalConfig.ClockInRoleID)
	if err != nil {
		if strings.Contains(err.Error(), "rate limit") || strings.Contains(err.Error(), "429") {
			os.WriteFile("rate_limited_marker", []byte("rate limited"), 0644)
		}
		m := "Error adding clock in role. Please try again later or message <@982491279369830460>." // its ya boi sush1sui
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &m})
		return
	}

	// Log the clock-in to the admin channel using the new format
    logMessage := fmt.Sprintf("🟢 <@%s> has clocked in at <t:%d:F>", i.Member.User.ID, clockInTime.Unix())
    s.ChannelMessageSend(config.GlobalConfig.AdminChannelID, logMessage)

	_, err = repository.ClockRecordService.ClockIn(i.Member.User.ID)
	if err != nil {
		m := "Error clocking in. Please try again later or message <@982491279369830460>." // its ya boi sush1sui
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &m})
		return
	}

	m := "You have successfully clocked in! Your clock in time has been recorded."
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &m})
}