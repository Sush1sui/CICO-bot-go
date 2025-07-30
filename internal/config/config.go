package config

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Config struct {
	PORT string
	BotToken string
	MongoDB_URI string
	MongoDB_Name string
	MongoDB_Clock_Channels_Name string
	MongoDB_Clock_Records_Name string
	TL_ROLE_ID string
	CHATTER_ROLE_ID string
	TimeLimit map[string]float64
	GuildID string
	ServerURL string
	AdminChannelID string
	ClockInRoleID string
}

var GlobalConfig *Config

func New() (error) {
	if err := godotenv.Load(); err != nil { fmt.Println("Error loading .env file") }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8169"
	}
	serverUrl := os.Getenv("SERVER_URL")
	if serverUrl == "" {
		fmt.Println("SERVER_URL is not set")
	}
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return fmt.Errorf("BOT_TOKEN is required")
	}
	mongoDB_URI := os.Getenv("MONGODB_URI")
	if mongoDB_URI == "" {
		return fmt.Errorf("MONGODB_URI is required")
	}
	mongoDB_Name := os.Getenv("MONGODB_NAME")
	if mongoDB_Name == "" {
		return fmt.Errorf("MONGODB_NAME is required")
	}
	mongoDB_Clock_Channels_Name := os.Getenv("MONGODB_CLOCK_CHANNELS_NAME")
	if mongoDB_Clock_Channels_Name == "" {
		return fmt.Errorf("MONGODB_CLOCK_CHANNELS_NAME is required")
	}
	mongoDB_Clock_Records_Name := os.Getenv("MONGODB_CLOCK_RECORDS_NAME")
	if mongoDB_Clock_Records_Name == "" {
		return fmt.Errorf("MONGODB_CLOCK_RECORDS_NAME is required")
	}
	tlRoleID := os.Getenv("TL_ROLE_ID")
	if tlRoleID == "" {
		return fmt.Errorf("TL_ROLE_ID is required")
	}
	chatterRoleID := os.Getenv("CHATTER_ROLE_ID")
	if chatterRoleID == "" {
		return fmt.Errorf("CHATTER_ROLE_ID is required")
	}
	guildID := os.Getenv("GUILD_ID")
	if guildID == "" {
		return fmt.Errorf("GUILD_ID is required")
	}

	GlobalConfig = &Config{
		PORT: port,
		BotToken: botToken,
		MongoDB_URI: mongoDB_URI,
		MongoDB_Name: mongoDB_Name,
		MongoDB_Clock_Channels_Name: mongoDB_Clock_Channels_Name,
		MongoDB_Clock_Records_Name: mongoDB_Clock_Records_Name,
		TL_ROLE_ID: tlRoleID,
		CHATTER_ROLE_ID: chatterRoleID,
		TimeLimit: map[string]float64{
			tlRoleID:      12.25, // 12 hours 15 minutes
			chatterRoleID: 16.25, // 16 hours 15 minutes
		},
		GuildID: guildID,
		ServerURL: serverUrl,
	}
	fmt.Println("Config initialized successfully")
	return nil
}

func MongoConnection() *mongo.Client {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
  client, err := mongo.Connect(opts)
  if err != nil {
    panic(err)
  }

  // Send a ping to confirm a successful connection
  if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
    panic(err)
  }
  fmt.Println("DB Connected!")

	return client
}