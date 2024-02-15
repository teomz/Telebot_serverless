package main

import (
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Assuming participants and mu are declared in main.go
extern var (
	participants map[int64]string
	mu           sync.Mutex
)

func handleJoin(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	mu.Lock()
	defer mu.Unlock()

	// Check if the user has already joined
	if _, exists := participants[message.From.ID]; exists {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You've already joined the game.")
		bot.Send(msg)
		return
	}

	// Add the user to participants if less than 4 have joined
	if len(participants) < 4 {
		participants[message.From.ID] = message.From.UserName
		msgText := fmt.Sprintf("You've joined the game. There are now %d participants.", len(participants))
		msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
		bot.Send(msg)

		// Notify when 4 participants have joined
		if len(participants) == 4 {
			notifyStartGame(bot)
		}
	} else {
		// Notify user that the game is full
		msg := tgbotapi.NewMessage(message.Chat.ID, "The game is currently full. Please wait for the next round.")
		bot.Send(msg)
	}
}

func notifyStartGame(bot *tgbotapi.BotAPI) {
	for userID := range participants {
		msg := tgbotapi.NewMessage(userID, "The game is starting now with 4 participants.")
		bot.Send(msg)
	}
}
