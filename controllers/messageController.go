package controllers

import (
	"bridge/entities"
	"bridge/utils"
	"errors"
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
	GameControllers []*GameController
}

//Constructor
func NewMessageController(bot *tgbotapi.BotAPI) *MessageController{
	fmt.Println("Created Message Controller")
	return &MessageController{
		bot:     bot,
		GameControllers: []*GameController{},
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

func (mc *MessageController) CheckGameController (gc *GameController) bool{
	if len(mc.GameControllers) == 0{
		return false //dont exist
	}
	for _,c := range mc.GameControllers{
		if c == gc{
			return true
		}
	}
	return false
}

func (mc *MessageController) AddGameController (gc *GameController){
	if !mc.CheckGameController(gc){
		mc.GameControllers = append(mc.GameControllers, gc)
		fmt.Println("Added game controller to list!")
		return
	} else{
		fmt.Printf("From existing game controller: chat %d\n", gc.chatID)
		return
	}
}

func (mc *MessageController)CheckOngoingController(chatID int64) bool{
	_,err:=mc.FindGameController(chatID)
	if err!=nil{
		fmt.Println(err)
		return false
	}else{
		return true
	}
}
//MessageHandler
func (mc *MessageController) HandleMessage(update tgbotapi.Update) {
	if update.Message.IsCommand(){
		command := update.Message.Command()
		switch command {
		case "start":
			utils.SendMessage(mc.bot,update.Message.Chat.ID, "Welcome to Bridge! Bridge is a four-player partnership trick-taking game with thirteen tricks per deal.\n\n/play_game : to play\n/leave: to leave\n/help: for more commands")
		case "help":
			utils.SendMessage(mc.bot,update.Message.Chat.ID, "Available commands:\n/start - Start the bot\n/help - Display help message")
		case "play_game":
			// Check game controller
			if !mc.CheckOngoingController(update.Message.Chat.ID){
				gameController:=GameController{mc.bot,update.Message.Chat.ID,entities.NewGame(mc.bot,update.Message.Chat.ID)}
				mc.AddGameController(&gameController)
				mc.PrintAllControllers()
				gameController.StartNewGame()
			}else{
				utils.SendMessage(mc.bot,update.Message.Chat.ID,fmt.Sprintf("%s, a game is already ongoing...",update.Message.From.UserName))
			}
		case "leave":
			gc,err := mc.FindGameController(update.Message.Chat.ID)
			if err != nil{
				fmt.Println(err)
			}else{
				gc.Game.RemovePlayer(update.Message.From)
				msg := fmt.Sprintf("%s has left room %d\n\nShutting down game...", update.Message.From, gc.Game.ID)
				utils.SendMessage(mc.bot,update.Message.Chat.ID,msg)
				gc.RemoveGame()
				mc.RemoveGameController(update.Message.Chat.ID)
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
			gc,err := mc.FindGameController(query.Message.Chat.ID)
			if err != nil{
				fmt.Println(err)
			}else{
				game := gc.Game
				if len(game.Players) < 4{
						gc.NotifyAddPlayer(query.From,roomID,msgID)
						game:= gc.Game
						game.CheckPlayers(mc.bot,query.Message.Chat.ID,roomID,msgID) //Check if room is full, else start game
				}
			}
		}
	default:
		// Handle other callback query scenarios
	}
}

func (mc *MessageController) FindGameController (chatID int64) (*GameController,error){
	for _,controller := range mc.GameControllers{
		if controller.chatID == chatID{
			return controller,nil
		}
	}
	return nil,errors.New("No controller found")
}

func (mc *MessageController) RemoveGameController (chatID int64){
	var index int
	for idx,controller := range mc.GameControllers{
		if controller.chatID == chatID{
			index = idx
			break
		}
	}
	mc.GameControllers = append(mc.GameControllers[:index],mc.GameControllers[index+1:]...)
}

func (mc *MessageController) PrintAllControllers (){
	for _,controller := range mc.GameControllers{
		fmt.Printf("ChatID: %d\n", controller.chatID)
	}
}