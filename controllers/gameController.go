package controllers

import (
	"bridge/entities"
	"bridge/utils"
	"errors"
	"fmt"
	"strings"
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

	Game := entities.NewGame(gc.bot,gc.chatID)
	gc.AddGame(Game)
	room := fmt.Sprintf("join_game:%d",Game.ID)
	btn := []tgbotapi.InlineKeyboardButton{utils.CreateButton("Join Game",room)}
	keyboard := utils.CreateInlineMarkup(btn)

	// Create a message with the inline keyboard
	utils.SendMessageWithMarkup(gc.bot,gc.chatID, "Starting a new game...\nNo of Players: 0",keyboard)
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
			var text []interface{}
			var tmp []string
			var result string
			text = append(text, fmt.Sprintf("Game Lobby %d\n\n", game.ID))
			for idx, player := range (game.Players){
				text = append(text,fmt.Sprintf("Player %d: %s\n", idx+1,player.UserName))
			}
			text = append(text, fmt.Sprintf("No of Players: %d, %d more players to start", len(game.Players),4-len(game.Players)))
			for _,item := range text{
				switch v:= item.(type){
				case string:
					tmp = append(tmp, v)
				}
			}
			result = strings.Join(tmp,"")
			room := fmt.Sprintf("join_game:%d",game.ID)
			btn := []tgbotapi.InlineKeyboardButton{utils.CreateButton("Join Game",room)}
			keyboard := utils.CreateInlineMarkup(btn)
			utils.EditMessageWithMarkup(gc.bot,gc.chatID,result,msgID,keyboard)
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




