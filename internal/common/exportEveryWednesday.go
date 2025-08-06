package common

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Sush1sui/cico-bot-go/internal/config"
	"github.com/Sush1sui/cico-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func ExportEveryWednesday(s *discordgo.Session) {
	go func() {
		for {
			now := time.Now().UTC().Add(8 * time.Hour) // convert to UTC+8
		
			// calculate the next Wednesday 6AM UTC+8
			daysUntilWednesday := (3 - int(now.Weekday()) + 7) % 7
			nextWednesday := now
			if daysUntilWednesday == 0 && now.Hour() >= 6 {
				nextWednesday = now.AddDate(0,0,7)
			} else {
				nextWednesday = now.AddDate(0, 0, daysUntilWednesday)
			}
			nextWednesday = time.Date(
				nextWednesday.Year(), nextWednesday.Month(), nextWednesday.Day(),
				6, 0, 0, 0, nextWednesday.Location(),
			)

			// convert back to UTC for timer
			nextWednesdayUTC := nextWednesday.Add(-8 * time.Hour)
			timeUntilExport := time.Until(nextWednesdayUTC)

			fmt.Printf("Next CSV export scheduled for: %s UTC+8\n", nextWednesday.Format(time.RFC3339))
			fmt.Printf("Time until export: %.2f hours\n", timeUntilExport.Hours())

			time.Sleep(timeUntilExport)

			err := ExportToCSV_CLEAN_DATABASE(s)
			if err != nil {
				fmt.Printf("Error during CSV export: %v\n", err)
				continue
			}
		}
	}()
	fmt.Println("Weekly CSV export started initialized.")
}

func ExportToCSV_CLEAN_DATABASE(s *discordgo.Session) (error) {
	// export to CSV
	fmt.Println("Starting CSV export...")
	filePath, err := ExportToCSV(s)
	if err != nil || filePath == "" {
		return fmt.Errorf("failed to export to CSV: %v", err)
	}
	fmt.Printf("CSV exported successfully to: %s\n", filePath)

	// send the CSV file to the admin channel
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	attachment := &discordgo.File{
		Name:        filepath.Base(filePath),
		ContentType: "text/csv",
		Reader:      file,
	}

	member, err := s.State.Member(config.GlobalConfig.GuildID, "608646101712502825")
	if err != nil {
		fmt.Printf("Failed to fetch member: %v\n", err)
	}
	if member != nil {
		// 1. Create a DM channel with the user
        dmChannel, err := s.UserChannelCreate("608646101712502825")
        if err != nil {
            fmt.Printf("Failed to create DM channel: %v\n", err)
        } else {
            // 2. Send the message to the created DM channel ID
            _, err = s.ChannelMessageSendComplex(dmChannel.ID, &discordgo.MessageSend{
                Content: "📊 Weekly clock records exported successfully! Here is the file:",
                Files:   []*discordgo.File{attachment},
            })
            if err != nil {
                fmt.Printf("Failed to send DM: %v\n", err)
            }
        }
		file.Seek(0,0) // Reset file pointer to the beginning before sending
	} else {
		fmt.Println("Member not found, sending without mention.")
	}
	
	_, err = s.ChannelMessageSendComplex(config.GlobalConfig.AdminChannelID, &discordgo.MessageSend{
		Content: "📊 Weekly clock records exported successfully! Here is the file:",
		Files:   []*discordgo.File{attachment},
	})
	if err != nil {
		fmt.Printf("Failed to send CSV file: %v\n", err)
	}

	file.Close()
	err = os.Remove(filePath)
	if err != nil {
		fmt.Println("Failed to remove CSV file:", err)
	}

	err = repository.ClockRecordService.RemoveClockRecordOfThoseNotClockedIn()
	if err != nil {
		fmt.Println("Failed to remove clock records:", err)
	}

	fmt.Println("Clock records cleaned up successfully.")
	return nil
}