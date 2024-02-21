package handler

import (
	"bridge/controllers"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	updates := ListenForWebhookRespReqFormat(w,r)

	bot := &tgbotapi.BotAPI{

		Token: os.Getenv("TELEGRAM_APITOKEN"),

		Client: &http.Client{},

		Buffer: 100,
	}

	bot.SetAPIEndpoint(tgbotapi.APIEndpoint)


	MessageController := controllers.NewMessageController(bot)

	for update := range updates:
		MessageController.StartListening()
}