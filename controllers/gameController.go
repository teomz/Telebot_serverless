package controllers

import (
	"bridge/entities"
	"bridge/utils"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var GlobalGameController GameController

type GameController struct{
	bot *tgbotapi.BotAPI
	chatID int64
	Games []*entities.Game
}

func NewGameController (bot *tgbotapi.BotAPI, chatID int64, games []*entities.Game) *GameController{
	return &GameController{
		bot:bot,
		chatID:chatID,
		Games:games,
	}
}

func (gc *GameController) StartNewGame() {
	fmt.Println("Start New Game")

	Game := entities.NewGame()
	gc.AddGame(Game)
	room := fmt.Sprintf("join_game:%d",Game.ID)
	btn := []tgbotapi.InlineKeyboardButton{utils.CreateButton("Join Game",room)}
	keyboard := utils.CreateInlineMarkup(btn)

	// Create a message with the inline keyboard
	msg := tgbotapi.NewMessage(gc.chatID, "Starting a new game...\nNo of Players: 0")
	msg.ReplyMarkup = keyboard

	// Send the message with the inline keyboard
	gc.bot.Send(msg)
}

func (gc *GameController) NotifyAddPlayer (user *tgbotapi.User, room uint32, msgID int){
	_,game,err := gc.GetGame(room)
	if err != nil {
		fmt.Println(err)
	} else {
		err:=game.AddPlayer(user)
		if err != nil{
			fmt.Println(err)
		}	else{
			newText := fmt.Sprintf("Starting a new game...\nNo of Players: %d", len(game.Players))
			room := fmt.Sprintf("join_game:%d",game.ID)
			btn := []tgbotapi.InlineKeyboardButton{utils.CreateButton("Join Game",room)}
			keyboard := utils.CreateInlineMarkup(btn)
			utils.EditMessageWithMarkup(gc.bot,gc.chatID,newText,msgID,&keyboard)
		}
	}
}

func (gc *GameController) GetGame (param interface{}) (int,*entities.Game,error){
	switch m := param.(type){
	case uint32:
		for idx,game := range gc.Games{
			if game.ID == m{
				return idx,game,nil
			}
		}
	case *tgbotapi.User:
		for idx,game := range gc.Games{
			for _,player := range game.Players{
				if player.ID == m.ID{
					return idx,game, nil
				}
			}
		}
	}
	return -1,nil,errors.New("cannot find gameID/ cannot find player")
}

func (gc *GameController) AddGame(game *entities.Game) {
	gc.Games = append(gc.Games, game)
}

func (gc *GameController) RemoveGame(game *entities.Game) {
	idx,_,err := gc.GetGame(game.ID)
	if err != nil{
		fmt.Println(err)
	}else{
		gc.Games = append(gc.Games[idx:], gc.Games[idx+1:]...)
		fmt.Printf("Deleted room %d\n", game.ID)
	}
}




