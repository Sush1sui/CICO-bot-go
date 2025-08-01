package common

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sush1sui/cico-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func InitializeClockInIfUnexpected(s *discordgo.Session) {
	const markerFile = "shutdown_marker"
	const rateLimitedMarkerFile = "rate_limited_marker"
	// if merker file does not exist, its an unexpected shutdown
	if _, err := os.Stat(markerFile); os.IsNotExist(err) {
		initializeClockIn(s)
	} else if _, err := os.Stat(rateLimitedMarkerFile); os.IsNotExist(err) {
		initializeClockIn(s)
	}
	
	os.Remove(markerFile)
	os.Remove(rateLimitedMarkerFile)
}

func MarkGracefulShutdown() {
	const markerFile = "shutdown_marker"
	os.WriteFile(markerFile, []byte("shutdown"), 0644)
}

func initializeClockIn(s *discordgo.Session) {
	allRecords, err := repository.ClockRecordService.GetAllClockRecords()
	if err != nil {
		s.ChannelMessageSend("error-channel-id", "Error fetching clock records: "+err.Error())
		return
	}

	var userMentions []string
	for _, record := range allRecords {
		if record.ClockInTime == nil { continue }

		_, err := repository.ClockRecordService.ReClockUser(record.UserID)
		if err != nil {
			fmt.Println("Error re-clocking user:", err)
			continue
		}

		userMentions = append(userMentions, "<@"+record.UserID+">")
	}
	if len(userMentions) > 0 {
		_, err := s.ChannelMessageSend("1355806810778636419", fmt.Sprintf("# Bot restarted\n\n%s\n\n**If you can't clock out while the bot is offline, you can clock out now.**", strings.Join(userMentions, ", ")))
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
	fmt.Println("Clock in initialized for clocked in users.")
}