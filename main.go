package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// Replace "YOUR_TELEGRAM_BOT_TOKEN" with the API token provided by the BotFather.
	bot, err := tgbotapi.NewBotAPI("6863492345:AAH-ak_depbfolBuCoI7PzfHu4ajJZ0L030")
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	// Set the bot to use debug mode (verbose logging).
	bot.Debug = true
	log.Printf("Authorized as @%s", bot.Self.UserName)
	// Set up updates configuration.
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// Get updates from the Telegram API.
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Error getting updates: %v", err)
	}
	// Process incoming messages.
	for update := range updates {
		if update.Message == nil { // Ignore any non-Message updates.
			continue
		}
		// Log the received message text and sender username.
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Extract the command from the message text.
		command := strings.Split(update.Message.Text, " ")[0]

		// Respond to the user based on the command.
		var responseText string
		switch command {
		case "/start":
			responseText = "Welcome to the Bridge bot! Use /help to see available commands."
		case "/join":
			responseText = "You've joined the game. Waiting for more players..."
		case "/leave":
			responseText = "You've left the game."
		case "/help":
			responseText = "Available commands:\n/start - start interacting with the bot\n/join - join a poker game\n/leave - leave the current game\n/fold - fold your hand\n/check - check during your turn"
		case "/fold":
			responseText = "You've folded."
		case "/check":
			responseText = "You've checked."
		default:
			responseText = "Sorry, I didn't recognize that command."
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}
