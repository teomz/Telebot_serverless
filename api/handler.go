package handler

import (
	"bridge/controllers"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	// Initialize the BotAPI instance with the bot token
	bot := &tgbotapi.BotAPI{
		Token:  os.Getenv("TELEGRAM_APITOKEN"),
		Client: &http.Client{},
		Buffer: 100,
	}

	// Set the API endpoint (optional)
	bot.SetAPIEndpoint(tgbotapi.APIEndpoint)

	// Create a new MessageController instance
	MessageController := controllers.NewMessageController(bot)

	// Listen for webhook updates using bot.ListenForWebhookRespReqFormat
	// updates := bot.ListenForWebhookRespReqFormat(w, r)

	// Process incoming updates
	var update tgbotapi.Update

	MessageController.StartListening(update)

}
