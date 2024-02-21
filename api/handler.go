package handler

import (
	"bridge/controllers"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize the bot API with the bot token
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		http.Error(w, "Failed to initialize bot", http.StatusInternalServerError)
		return
	}

	// Set the API endpoint (optional)
	bot.SetAPIEndpoint(tgbotapi.APIEndpoint)

	// Create a new MessageController instance
	messageController := controllers.NewMessageController(bot)

	// Listen for webhook updates using bot.ListenForWebhook
	updates := bot.ListenForWebhook("/" + bot.Token)

	// Process incoming updates
	for update := range updates {
		messageController.StartListening(update)
	}
}
