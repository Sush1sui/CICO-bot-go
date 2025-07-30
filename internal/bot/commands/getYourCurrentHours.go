package commands

import (
	"slices"
	"strconv"
	"time"

	"github.com/Sush1sui/cico-bot-go/internal/config"
	"github.com/Sush1sui/cico-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func GetYourCurrentHours(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" { return }

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Fetching your current hours...",
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})

	if !slices.Contains(i.Member.Roles, config.GlobalConfig.TL_ROLE_ID) && !slices.Contains(i.Member.Roles, config.GlobalConfig.CHATTER_ROLE_ID) {
		m := "You are not a chatter or team leader to use this command."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &m})
		return
	}

	clockRecord, err := repository.ClockRecordService.GetUserClockRecord(i.Member.User.ID)
	if err != nil || clockRecord == nil {
		mess := "Failed to fetch your clock record: " + err.Error()
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		return
	}

	if clockRecord.TotalHours == nil || *clockRecord.TotalHours == 0 {
		mess := "Your total hours have not been set yet, try to clock in and clock out later, then try checking again."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
		return
	}

	timeSince := time.Since(*clockRecord.ClockInTime).Hours()
	mess := "**Hours count since your last clock in:** *" + strconv.Itoa(int(timeSince)) + " hours.*\n" +
		"**Total hours (not including today):** *" + strconv.Itoa(int(*clockRecord.TotalHours)) + " hours.*\n" +
		"*Clock out if you want to include the total hours for today, just clock in again if you want to continue counting hours for today.*"
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})

}