package deploy

import (
	"fmt"

	"github.com/Sush1sui/cico-bot-go/internal/bot/events"
	"github.com/bwmarrin/discordgo"
)

var eventHandlers = []any{
	events.OnClockIn,
	events.OnClockOut,
}

func DeployEvents(s *discordgo.Session) {
	if len(eventHandlers) == 0 { return }

	for _, handler := range eventHandlers {
		s.AddHandler(handler)
	}

	fmt.Println("Event handlers deployed successfully")
}