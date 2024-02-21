package handler

import (
	"bridge/controllers"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(w *http.ResponseWriter, r *http.Request) {

	bot := &tgbotapi.BotAPI{

		Token: os.Getenv("TELEGRAM_APITOKEN"),

		Client: &http.Client{},

		Buffer: 100,
	}

	bot.SetAPIEndpoint(tgbotapi.APIEndpoint)
	MessageController := controllers.NewMessageController(bot)

	MessageController.StartListening(w, r)
	select {}

}
