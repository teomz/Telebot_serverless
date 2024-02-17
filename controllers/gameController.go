package controllers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var gameController GameController

type GameController struct{
	bot *tgbotapi.BotAPI
}

func NewGameController (bot *tgbotapi.BotAPI) * GameController{
	return &GameController{
		bot:bot,
	}
}

func (gc *GameController) StartGame() {
	fmt.Println("Start Game")
}
