package handler

import (
	"bridge/controllers"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)

	var update tgbotapi.Update

	err := json.Unmarshal(body, &update)

	if err != nil {

		log.Println(err)

		return

	}

	bot := &tgbotapi.BotAPI{

		Token: os.Getenv("TELEGRAM_APITOKEN"),

		Client: &http.Client{},

		Buffer: 100,
	}

	bot.SetAPIEndpoint(tgbotapi.APIEndpoint)


	MessageController := controllers.NewMessageController(bot)
	MessageController.StartListening()

}