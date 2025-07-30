package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sush1sui/cico-bot-go/internal/bot"
	"github.com/Sush1sui/cico-bot-go/internal/config"
	"github.com/Sush1sui/cico-bot-go/internal/server"
)

func main() {
	err := config.New()
	if err != nil {
		fmt.Println("Error initializing configuration:", err)
	}

	mongoClient := config.MongoConnection()
	defer mongoClient.Disconnect(context.Background())
	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		panic(fmt.Errorf("failed to connect to MongoDB: %w", err))
	}

	clockChannelsCollection := mongoClient.Database(config.GlobalConfig.MongoDB_Name).Collection(config.GlobalConfig.MongoDB_Clock_Channels_Name)
	clockRecordsCollection := mongoClient.Database(config.GlobalConfig.MongoDB_Name).Collection(config.GlobalConfig.MongoDB_Clock_Records_Name)

	addr := fmt.Sprintf(":%s", config.GlobalConfig.PORT)
	router := server.NewRouter()
	fmt.Printf("Server is running on PORT: %s\n", addr)

	go func() {
		if err := http.ListenAndServe(addr, router); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	go bot.StartBot()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Shutting down server gracefully...")
}