package common

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Sush1sui/cico-bot-go/internal/config"
	"github.com/Sush1sui/cico-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func ExportToCSV(s *discordgo.Session) (string, error) {
	records, err := repository.ClockRecordService.GetAllClockRecords()
	if err != nil {
		fmt.Println("Error fetching clock records:", err)
		return "", err
	}

	// prep csv rows
	csvRows := [][]string{
		{"User ID", "Username", "Total Hours"},
	}

	guild, err := s.State.Guild(config.GlobalConfig.GuildID)
	if err != nil {
		fmt.Println("Error fetching guild:", err)
		return "", err
	}
	guildMembers := make(map[string]string)
	for _, member := range guild.Members {
		guildMembers[member.User.ID] = member.User.Username
	}

	for _, record := range records {
		var totalHours int
		if record.TotalHours != nil {
			decimal := *record.TotalHours - float64(int(*record.TotalHours))
			roundOff := decimal >= 0.5
			if roundOff {
				totalHours = int(*record.TotalHours) + 1
			} else {
				totalHours = int(*record.TotalHours)
			}
		} else {
			totalHours = 0
		}

		username := guildMembers[record.UserID]
		if username == "" {
			username = "Unknown"
		}

		row := []string{
			record.UserID,
			username,
			strconv.Itoa(totalHours),
		}
		csvRows = append(csvRows, row)
	}

	// prep file path
	today := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("clock_records_%s.csv", today)
	csvFolderPath := filepath.Join(".", "csv")
	filePath := filepath.Join(csvFolderPath, fileName)

	if err := os.MkdirAll(csvFolderPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", csvFolderPath, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range csvRows {
		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("failed to write record to CSV: %w", err)
		}
	}
	
	fmt.Println("CSV file created successfully:", filePath)
	return filePath, nil
}