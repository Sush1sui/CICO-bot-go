package common

import (
	"fmt"
	"time"

	"github.com/Sush1sui/cico-bot-go/internal/repository"
	"github.com/bwmarrin/discordgo"
)

func CheckForExpiredClock(s *discordgo.Session) error {
	for {
		time.Sleep(20 * time.Minute)
		go func() {
			err := repository.ClockRecordService.CheckForExpiredClock(s)
			if err != nil {
				fmt.Printf("Error checking for expired clock: %v\n", err)
			}
		}()
	}
}