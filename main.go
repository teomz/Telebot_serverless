package main

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("6863492345:AAH-ak_depbfolBuCoI7PzfHu4ajJZ0L030"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true
}
