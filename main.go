package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6863492345:AAH-ak_depbfolBuCoI7PzfHu4ajJZ0L030") // Replace with your Bot Token
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Error getting updates: %v", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		command := strings.Split(update.Message.Text, " ")[0]

		switch command {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to the Bridge bot! Use /help to see available commands.")
			bot.Send(msg)
		case "/help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Available commands:\n/start - start interacting with the bot\n/join - join a Bridge game\n/leave - leave the current game\n/fold - fold your hand\n/check - check during your turn\n/play_game - play the game")
			bot.Send(msg)
		case "/play_game":
			gameShortName := "your_game_short_name" // Replace with your game short name
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Let's play a game!")
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonSwitch("🎲 Play Game", gameShortName),
				),
			)
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		case "/join", "/leave", "/fold", "/check":
			// Placeholder for future command implementation
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This command is not implemented yet.")
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, I didn't recognize that command.")
			bot.Send(msg)
		}
	}
}
