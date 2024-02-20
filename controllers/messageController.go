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
	Chats []tgbotapi.Chat
}

//Constructor
func NewMessageController(bot *tgbotapi.BotAPI) *MessageController{
	fmt.Println("Created Message Controller")
	return &MessageController{
		bot:     bot,
		gameControllerLock: false,
		Chats: []tgbotapi.Chat{},
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

func (mc *MessageController) CheckChat (chat *tgbotapi.Chat) bool{
	if len(mc.Chats) == 0{
		return false //dont exist
	}
	for _,c := range mc.Chats{
		if c == *chat{
			return true
		}
	}
	return false
}

func (mc *MessageController) AddChat (chat *tgbotapi.Chat){
	if !mc.CheckChat(chat){
		mc.Chats = append(mc.Chats, *chat)
		fmt.Println("Added chat to list!")
		return
	} else{
		fmt.Printf("From existing chat: chat %d\n", chat.ID)
		return
	}
}

func (mc *MessageController)CheckOngoingGame(chat *tgbotapi.Chat) bool{
	if len(GlobalGameController.Games)==0{
		return false
	}else{
		for _,g := range (GlobalGameController.Games){
			if g.ChatID == chat.ID{
				return true
			}
		}
	}
	return false
}
//MessageHandler
func (mc *MessageController) HandleMessage(update tgbotapi.Update) {
	if update.Message.IsCommand(){
		mc.AddChat(update.Message.Chat)
		command := update.Message.Command()
		switch command {
		case "start":
			utils.SendMessage(mc.bot,update.Message.Chat.ID, "Welcome to Bridge! Bridge is a four-player partnership trick-taking game with thirteen tricks per deal.\n\n/play_game : to play\n/leave: to leave\n/help: for more commands")
		case "help":
			utils.SendMessage(mc.bot,update.Message.Chat.ID, "Available commands:\n/start - Start the bot\n/help - Display help message")
		case "play_game":
			// To ensure only one instance of GameController is initialized
			if mc.gameControllerLock{
				fmt.Println("Game Controller exist")
				if !mc.CheckOngoingGame(update.Message.Chat){
					GlobalGameController.StartNewGame()
				}else{
					utils.SendMessage(mc.bot,update.Message.Chat.ID,fmt.Sprintf("%s, a game is already ongoing in this chat",update.Message.From.UserName))
				}
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
				msg := fmt.Sprintf("%s has left room %d\n\nShutting down game...", update.Message.From, game.ID)
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
	mc.AddChat(query.Message.Chat)

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
						_,game,err := GlobalGameController.GetGame(roomID)
						if err != nil{
							fmt.Println(err)
						} else{
							game.CheckPlayers(mc.bot,query.Message.Chat.ID,roomID,msgID) //Check if room is full, else start game
						}
				}
			}
		}
	default:
		// Handle other callback query scenarios
	}
}