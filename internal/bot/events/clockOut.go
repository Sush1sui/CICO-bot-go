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

func OnClockOut(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" { return }
	if i.Type != discordgo.InteractionMessageComponent { return }
	if i.MessageComponentData().CustomID != "clock_out" { return }

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Clocking out, please wait...",
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})

	if !slices.Contains(i.Member.Roles, config.GlobalConfig.ClockInRoleID) {
		m := "Clock out failed, you are already clocked out."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &m})
		return
	}

	if !slices.Contains(i.Member.Roles, config.GlobalConfig.TL_ROLE_ID) && !slices.Contains(i.Member.Roles, config.GlobalConfig.CHATTER_ROLE_ID) {
		m := "You do not have Chatter or Team Leader role to clock in."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &m})
		return
	}

	// remove clock in role from the user
	clockOutTime := time.Now()
	err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, config.GlobalConfig.ClockInRoleID)
	if err != nil {
		if strings.Contains(err.Error(), "rate limit") || strings.Contains(err.Error(), "429") {
			os.WriteFile("rate_limited_marker", []byte("rate limited"), 0644)
		}
		m := "Error removing clock in role. Please try again later or message <@982491279369830460>." // its ya boi sush1sui
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &m})
		return
	}

	// Log the clock-out to the admin channel using the new format
    logMessage := fmt.Sprintf("🔴 <@%s> has clocked out at <t:%d:F>", i.Member.User.ID, clockOutTime.Unix())
    s.ChannelMessageSend(config.GlobalConfig.AdminChannelID, logMessage)

	_, err = repository.ClockRecordService.ClockOut(i.Member.User.ID)
	if err != nil {
		fmt.Println("Error recording clock out time:", err)
		m := "Error recording clock out time. Please try again later or message <@982491279369830460>." // its ya boi sush1sui
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &m})
		return
	}

	m := "You have successfully clocked out."
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &m})
}