package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sush1sui/cico-bot-go/internal/bot/deploy"
	"github.com/Sush1sui/cico-bot-go/internal/config"
	"github.com/bwmarrin/discordgo"
)

func StartBot() {
	s, e := discordgo.New("Bot "+config.GlobalConfig.BotToken)
	if e != nil {
		log.Fatalf("error creating Discord session: %v", e)
	}

	s.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildPresences | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages

	s.AddHandler(func(sess *discordgo.Session, ready *discordgo.Ready) {
		sess.UpdateStatusComplex(discordgo.UpdateStatusData{
			Status: "idle",
			Activities: []*discordgo.Activity{
				{
					Name: "to Clock in Clock outs",
					Type: discordgo.ActivityTypeListening,
				},
			},
		})
	})

	e = s.Open()
	if e != nil {
		log.Fatalf("error opening connection: %v", e)
	}
	defer s.Close()

	deploy.DeployCommands(s)
	deploy.DeployEvents(s)
	
	fmt.Println("Bot is now running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("Shutting down bot gracefully...")
}