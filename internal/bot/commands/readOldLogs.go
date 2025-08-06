package commands

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Sush1sui/cico-bot-go/internal/config"
	"github.com/bwmarrin/discordgo"
)

// Updated regexes for the new log format with emojis and Unix timestamps
var (
    clockInRegex  = regexp.MustCompile(`🟢 <@(\d+)> has clocked in at <t:(\d+):F>`)
    clockOutRegex = regexp.MustCompile(`🔴 <@(\d+)> has clocked out at <t:(\d+):F>`)
    expiredRegex  = regexp.MustCompile(`⚠️ <@(\d+)> has exceeded the time limit`)
)

func ReadOldLogs(s *discordgo.Session, i *discordgo.InteractionCreate) {
    // Local maps to ensure data is fresh for each command run
    totalHours := make(map[string]time.Duration)
    openSessions := make(map[string]time.Time)

    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Content: "Reading old logs, this may take a moment...",
            Flags:   discordgo.MessageFlagsEphemeral,
        },
    })

    if i.Member.User.ID != "982491279369830460" && i.Member.User.ID != "608646101712502825" {
        mess := "You do not have permission to use this command."
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
        return
    }

    startOfMessageId := i.ApplicationCommandData().Options[0].StringValue()
    if startOfMessageId == "" {
        mess := "Please provide a valid start message ID."
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
        return
    }

    var allMessages []*discordgo.Message
    lastID := ""
    for {
        messages, err := s.ChannelMessages(config.GlobalConfig.AdminChannelID, 100, lastID, "", "")
        if err != nil {
            mess := "Failed to read old logs: " + err.Error()
            s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
            return
        }
        if len(messages) == 0 {
            break
        }
        allMessages = append(allMessages, messages...)
        lastID = messages[len(messages)-1].ID
        if lastID <= startOfMessageId {
            break
        }
    }

    for i, j := 0, len(allMessages)-1; i < j; i, j = i+1, j-1 {
        allMessages[i], allMessages[j] = allMessages[j], allMessages[i]
    }

    for _, message := range allMessages {
        if message.ID >= startOfMessageId {
            processLogLine(message.Content, totalHours, openSessions)
        }
    }

    var result strings.Builder
    result.WriteString("## Reconstructed Hours:\n")
    if len(totalHours) == 0 {
        mess := "No completed clock-in/out sessions found in the specified logs. No CSV generated."
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
        return
    }

    // if len(openSessions) > 0 {
    //     result.WriteString("\n### Still Clocked In (Unclosed Sessions):\n")
    //     for userID := range openSessions {
    //         result.WriteString(fmt.Sprintf("- <@%s>\n", userID))
    //     }
    // }

	// --- CSV Generation ---
    fileName := fmt.Sprintf("reconstructed_hours_%s.csv", time.Now().Format("2006-01-02_15-04-05"))
    filePath := "./csv/" + fileName
    file, err := os.Create(filePath)
    if err != nil {
        mess := "Failed to create CSV file: " + err.Error()
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
        return
    }

    writer := csv.NewWriter(file)

    // Write CSV header
    headers := []string{"User ID", "Username", "Total Hours"}
    if err := writer.Write(headers); err != nil {
        mess := "Failed to write CSV header: " + err.Error()
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
        file.Close() // Close before returning
        return
    }

    // Write data rows
    for userID, duration := range totalHours {
        user, err := s.User(userID)
        userName := userID // Default to ID if user fetch fails
        if err == nil {
            userName = user.Username
        }

        hours := fmt.Sprintf("%.2f", duration.Hours())
        row := []string{userID, userName, hours}
        if err := writer.Write(row); err != nil {
            // Log the error but continue trying to write other rows
            fmt.Printf("Error writing row to csv for user %s: %v\n", userID, err)
        }
    }

    writer.Flush()
    file.Close()

    // --- Send the CSV file ---
    csvFile, err := os.Open(filePath)
    if err != nil {
        mess := "Failed to open CSV file for sending: " + err.Error()
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
        return
    }
    defer csvFile.Close()
    defer os.Remove(filePath) // Clean up the file from the server afterwards

	 // --- DM the file to a specific user first ---
    dmChannel, err := s.UserChannelCreate("1387738107675410515")
    if err != nil {
        fmt.Printf("Failed to create DM channel with user 1387738107675410515: %v\n", err)
    } else {
        _, err = s.ChannelMessageSendComplex(dmChannel.ID, &discordgo.MessageSend{
            Content: "Here is the reconstructed hours CSV you requested.",
            Files: []*discordgo.File{
                {
                    Name:   fileName,
                    Reader: csvFile,
                },
            },
        })
        if err != nil {
            fmt.Printf("Failed to send DM to user 1387738107675410515: %v\n", err)
        }
        // IMPORTANT: Rewind the file reader to be sent again to the channel
        _, err = csvFile.Seek(0, 0)
        if err != nil {
            mess := "Failed to reset file reader after DM: " + err.Error()
            s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &mess})
            return
        }
    }

	// --- Send the CSV file to the admin channel ---
	_, err = s.ChannelMessageSendComplex(config.GlobalConfig.AdminChannelID, &discordgo.MessageSend{
		Content: "Here is the reconstructed hours CSV file.",
		Files: []*discordgo.File{
			{
				Name:   fileName,
				Reader: csvFile,
			},
		},
	})
	if err != nil {
		fmt.Printf("Failed to send CSV file to admin channel: %v\n", err)
	}

	// --- Cleanup ---
    // Explicitly close the file first, then remove it.
    csvFile.Close()
    if removeErr := os.Remove(filePath); removeErr != nil {
        fmt.Printf("Failed to remove CSV file: %v\n", removeErr)
    }

	messageLink := fmt.Sprintf("https://discord.com/channels/%s/%s/%s", i.GuildID, config.GlobalConfig.AdminChannelID, startOfMessageId)

    content := fmt.Sprintf("Here is the CSV file with the reconstructed hours since message: %s.", messageLink)
    s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
        Content: &content,
    })
}

func processLogLine(line string, totalHours map[string]time.Duration, openSessions map[string]time.Time) {
    if matches := clockInRegex.FindStringSubmatch(line); len(matches) > 2 {
        userID := matches[1]
        timestampStr := matches[2]
        unixTimestamp, err := strconv.ParseInt(timestampStr, 10, 64)
        if err == nil {
            openSessions[userID] = time.Unix(unixTimestamp, 0)
        }
    } else if matches := clockOutRegex.FindStringSubmatch(line); len(matches) > 2 {
        userID := matches[1]
        timestampStr := matches[2]
        if clockInTime, ok := openSessions[userID]; ok {
            unixTimestamp, err := strconv.ParseInt(timestampStr, 10, 64)
            if err == nil {
                clockOutTime := time.Unix(unixTimestamp, 0)
                duration := clockOutTime.Sub(clockInTime)
                totalHours[userID] += duration
                delete(openSessions, userID) // Close the session
            }
        }
    } else if matches := expiredRegex.FindStringSubmatch(line); len(matches) > 1 {
        userID := matches[1]
        delete(openSessions, userID) // Discard the session
    }
}