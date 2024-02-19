package entities

import (
	"errors"
	"math/rand"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Game struct{
	ID uint64
	Players []*tgbotapi.User
	deck Deck
	currentBid Bid
	hands []Hand
	InProgress bool
}

func NewGame () *Game{
	return &Game{
		ID: rand.Uint64(),
		Players: []*tgbotapi.User{},
		deck : Deck{},
		currentBid: Bid{},
		hands: []Hand{},
		InProgress: false,
	}
}

func (g *Game) FindPlayer (user *tgbotapi.User) bool{
	for _,player := range g.Players{
		if player.ID == user.ID{
			return true
		}
	}
	return false
}

func (g *Game) AddPlayer (user *tgbotapi.User) (error){
	if !g.FindPlayer(user) {
		g.Players = append(g.Players, user)
		return nil
	}
	return errors.New("player already in game")
}









