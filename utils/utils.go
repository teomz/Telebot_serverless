package utils

import(
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

)

func CreateButton(label, data string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(label, data)
}