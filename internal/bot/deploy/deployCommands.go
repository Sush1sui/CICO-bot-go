package deploy

import (
	"fmt"
	"log"

	"github.com/Sush1sui/cico-bot-go/internal/bot/commands"
	"github.com/Sush1sui/cico-bot-go/internal/config"
	"github.com/bwmarrin/discordgo"
)

var slashCommands = []*discordgo.ApplicationCommand{
	{
		Name: "generate-clock-channels",
		Description: "Generates clock-in, clock-out, and admin channels for the server",
		Type: discordgo.ChatApplicationCommand,
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(),
	},
	{
		Name: "delete-generated-channels",
		Description: "Deletes all generated clock channels in the server",
		Type: discordgo.ChatApplicationCommand,
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(),
	},
	{
		Name: "export-current-clock-records",
		Description: "Exports the current clock records to a CSV file",
		Type: discordgo.ChatApplicationCommand,
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(),
	},
	{
		Name: "get-your-current-total-hours",
		Description: "Shows your total clocked hours for this week and today (if you clocked in today)",
		Type: discordgo.ChatApplicationCommand,
	},
	{
		Name: "export-csv-reset-clocks",
		Description: "Exports the current clock records to a CSV file and resets the clocks",
		Type: discordgo.ChatApplicationCommand,
		DefaultMemberPermissions: func() *int64 { p := int64(discordgo.PermissionAdministrator); return &p }(),
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"generate-clock-channels": commands.GenerateClockChannels,
	"delete-generated-channels": commands.DeleteGeneratedChannels,
	"export-current-clock-records": commands.ExportCSVCommand,
	"get-your-current-total-hours": commands.GetYourCurrentHours,
	"export-csv-reset-clocks": commands.ExportCSVWithResetDatabaseCommand,
}

func DeployCommands(s *discordgo.Session) {
	globalCmds, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		for _, cmd := range globalCmds {
			err := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
			if err != nil {
				log.Printf("Failed to delete command %s: %v", cmd.Name, err)
			} else {
				log.Printf("Deleted command: %s", cmd.Name)
			}
		}
	}

	if len(slashCommands) == 0 { return }

	guilds := s.State.Guilds
	for _, guild := range guilds {
		if guild.ID == config.GlobalConfig.GuildID {
			_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, guild.ID, slashCommands)
			if err != nil {
				log.Printf("Failed to deploy commands to guild %s: %v", guild.Name, err)
			} else {
				log.Printf("Successfully deployed commands to guild: %s", guild.Name)
			}
			break
		}
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand { return }

		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		} else {
			fmt.Println("Unknown command:", i.ApplicationCommandData().Name)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Unknown command. Please try again.",
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
		}
	})
}