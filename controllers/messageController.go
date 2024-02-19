package controllers

import (
	"bridge/utils"
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
			utils.SendMessageWithMarkup(mc.bot,update.Message.Chat.ID, "Welcome to Bridge! Bridge is a four-player partnership trick-taking game with thirteen tricks per deal.",keyboard)
		case "help":
			utils.SendMessage(mc.bot,update.Message.Chat.ID, "Available commands:\n/start - Start the bot\n/help - Display help message")
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
			_,game,err:=GlobalGameController.GetGame(update.Message.From)
			if err != nil{
				fmt.Println(err)
			}else{
				game.RemovePlayer(update.Message.From)
				msg := fmt.Sprintf("%s has left room %d\n\n\nShutting down game...", update.Message.From, game.ID)
				utils.SendMessage(mc.bot,update.Message.Chat.ID,msg)
				GlobalGameController.RemoveGame(game)
			}
		default:
			utils.SendMessage(mc.bot,update.Message.Chat.ID, "Unknown command. Type /help for a list of available commands.")
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
		roomID,err:=strconv.ParseUint(data,10,32)
		fmt.Printf("Room ID: %d, user: %s pressed the button.\n",roomID, user.UserName)
		if err != nil {
			fmt.Println(err)
		} else {
			roomID:=uint32(roomID)
			_,game,err := GlobalGameController.GetGame(roomID)
			if err != nil {
				fmt.Println(err)
			} else{
				if len(game.Players) < 4{
						GlobalGameController.NotifyAddPlayer(query.From,roomID,msgID)
				} else{
					fmt.Printf("%d players in room. Starting Game...\n", len(game.Players))
					//Deletes button to join game
					utils.DeleteButton(mc.bot,query.Message.Chat.ID,msgID)
					game.InProgress = true
					msg:=fmt.Sprintf("Room %d\n\nPlayer 1: %s\nPlayer 2: %s\nPlayer 3: %s\nPlayer 4: %s\n\nStarting game now...",roomID,game.Players[0].UserName,game.Players[1].UserName,game.Players[2].UserName,game.Players[3].UserName)
					// msg:=fmt.Sprintf("Room %d\n\nPlayer 1: %s\n\nStarting game now...",roomID,game.Players[0].UserName)
					utils.SendMessage(mc.bot,query.Message.Chat.ID,msg)
					//Start Game Sequence
					game.StartGame()
				}
			}
		}
	default:
		// Handle other callback query scenarios
	}
}