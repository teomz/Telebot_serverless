package controllers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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
		if update.Message == nil{
			if update.CallbackQuery !=nil{
				mc.HandleCallbackQuery(update.CallbackQuery)
			} else{
				continue
			}
		} else{
			mc.HandleMessage(update)
		}
	}
}

//MessageHandler
func (mc *MessageController) HandleMessage(update tgbotapi.Update) {
	if update.Message.IsCommand(){
		command := update.Message.Command()
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
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to Bridge! Bridge is a four-player partnership trick-taking game with thirteen tricks per deal.")
			msg.ReplyMarkup = keyboard

			// Send the message
			mc.bot.Send(msg)
		case "help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Available commands:\n/start - Start the bot\n/help - Display help message")
			mc.bot.Send(msg)
		case "play_game":
			// To ensure only one instance of GameController is initialized
			if mc.gameControllerLock{
				fmt.Println("Game Controller exist")
				GlobalGameController.StartNewGame()
			} else {
				fmt.Println("Game Controller don't exist")
				GlobalGameController = GameController{mc.bot,update.Message.Chat.ID,nil}
				GlobalGameController.StartNewGame()
				mc.gameControllerLock = true
			}
		case "leave":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Leaving...")
			mc.bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Type /help for a list of available commands.")
			mc.bot.Send(msg)
		}
	}
}

//Callback Query Handler
func (mc *MessageController) HandleCallbackQuery (query *tgbotapi.CallbackQuery) {
	// Extract relevant information from the callback query
	user := query.From
	msgID := query.Message.MessageID
	parts := strings.Split(query.Data,":")
	command := parts[0]
	data := parts[1]

	// Handle the callback query logic based on the data
	switch command {
	case "join_game":
		// Respond to the button click
		roomID,err:=strconv.ParseUint(data,10,64)
		fmt.Printf("%d\n",roomID)
		if err != nil {
			fmt.Println(err)
		} else {
			GlobalGameController.AddPlayer(user,roomID,msgID)
		}
	default:
		// Handle other callback query scenarios
	}
}