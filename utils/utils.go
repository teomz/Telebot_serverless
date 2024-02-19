package utils

import(
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

)

func CreateButton(label, data string) tgbotapi.InlineKeyboardButton {
	button:=tgbotapi.NewInlineKeyboardButtonData(label, data)
	return button
}

func CreateButtons(label, data []string) []tgbotapi.InlineKeyboardButton {
	var buttons []tgbotapi.InlineKeyboardButton
	for i := 0; i < len(label); i++{
		button:=tgbotapi.NewInlineKeyboardButtonData(label[i], data[i])
		buttons = append(buttons, button)
	}
	return buttons
}

func DeleteButton(bot *tgbotapi.BotAPI, chatID int64, msgID int) {
	msg:=tgbotapi.NewDeleteMessage(chatID,msgID)
	bot.DeleteMessage(msg)
}

func SendMessage (bot *tgbotapi.BotAPI, chatID int64, text string){
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}

func SendMessageWithMarkup (bot *tgbotapi.BotAPI, chatID int64, text string,markup interface{}){
	msg := tgbotapi.NewMessage(chatID, text)

    switch m := markup.(type) {
    case tgbotapi.InlineKeyboardMarkup:
        msg.ReplyMarkup = m
    case tgbotapi.ReplyKeyboardMarkup:
        msg.ReplyMarkup = m
    default:
        // Handle unsupported keyboard type
        fmt.Println("Unsupported keyboard type")
        return
    }

    bot.Send(msg)
}

func CreateInlineMarkup (buttons []tgbotapi.InlineKeyboardButton) tgbotapi.InlineKeyboardMarkup{
	var row []tgbotapi.InlineKeyboardButton

	if len(buttons) == 1 {
        // If there's only one button, create a new row with that button
        row = append(row, buttons[0])
    } else if len(buttons) > 1 {
        // If there are multiple buttons, use the provided buttons as the row
        row = buttons
    }

	keyboard:=tgbotapi.NewInlineKeyboardMarkup(row)
	return keyboard
}

func EditMessageWithMarkup (bot *tgbotapi.BotAPI, chatID int64, text string, msgID int, markup *tgbotapi.InlineKeyboardMarkup){
	msg := tgbotapi.NewEditMessageText(chatID,msgID,text)
	msg.ReplyMarkup = markup
    bot.Send(msg)
}