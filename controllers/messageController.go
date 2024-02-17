package controllers

import (
	"log"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Class
type MessageController struct{
	//Member variables
	bot *tgbotapi.BotAPI
	gameControllerLock bool
}

//Constructor
func NewMessageController(bot *tgbotapi.BotAPI) *MessageController{
	fmt.Println("Created Message Controller")
	return &MessageController{
		bot:     bot,
		gameControllerLock: false,
	}
}

//Listener
func (mc *MessageController) StartListening() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	fmt.Println("Start Listening...")

	updates, err := mc.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}


	for update := range updates {
		if update.Message == nil {
			continue
		}
		mc.HandleMessage(update.Message)
	}
}

//MessageHandler
func (mc *MessageController) HandleMessage(message *tgbotapi.Message) {
	if message.IsCommand() {
		command := message.Command()
		switch command {
		case "start":
			// Create a custom keyboard
			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("/help"),
					tgbotapi.NewKeyboardButton("/play_game"),
					tgbotapi.NewKeyboardButton("/leave"),
				),
			)
			// Hide the custom keyboard once a button is pressed
			keyboard.OneTimeKeyboard = true

			// Create a message with the keyboard markup
			msg := tgbotapi.NewMessage(message.Chat.ID, "Welcome to Bridge! Bridge is a four-player partnership trick-taking game with thirteen tricks per deal.")
			msg.ReplyMarkup = keyboard

			// Send the message
			mc.bot.Send(msg)
		case "help":
			msg := tgbotapi.NewMessage(message.Chat.ID, "Available commands:\n/start - Start the bot\n/help - Display help message")
			mc.bot.Send(msg)
		case "play_game":
			// To ensure only one instance of GameController is initialized
			if mc.gameControllerLock{
				fmt.Println("Game Controller exist")
				gameController.StartGame()
			} else {
				fmt.Println("Game Controller don't exist")
				gameController := GameController{mc.bot}
				gameController.StartGame()
				mc.gameControllerLock = true
			}
		case "leave":
			msg := tgbotapi.NewMessage(message.Chat.ID, "Leaving...")
			mc.bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Type /help for a list of available commands.")
			mc.bot.Send(msg)
		}
	} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Not a command. Try again")
			mc.bot.Send(msg)
	}
}