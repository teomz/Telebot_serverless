package main

import (
	"log"
	"os"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"bridge/controllers"
)


func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the Telegram bot token from the environment variable
	botToken := os.Getenv("TELEGRAM_APITOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_APITOKEN not found in environment variables")
	}

	fmt.Println("Token extracted")

	// Create a new bot instance
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot instance created")

	MessageController := controllers.NewMessageController(bot)
	MessageController.StartListening()
}
