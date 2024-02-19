package entities

import (
	"errors"
	"math/rand"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Game struct{
	ID uint32
	Players []*tgbotapi.User
	deck Deck
	currentBid Bid
	hands []Hand
	InProgress bool
}

func NewGame () *Game{
	return &Game{
		ID: rand.Uint32(),
		Players: []*tgbotapi.User{},
		deck : Deck{},
		currentBid: Bid{},
		hands: []Hand{},
		InProgress: false,
	}
}

// TODO
func (g *Game) StartGame (){
	//Give out cards
	//Start Bid
}

func (g *Game) FindPlayer (user *tgbotapi.User) (bool,int){
	for idx,player := range g.Players{
		if player.ID == user.ID{
			return true,idx
		}
	}
	return false,-1
}

func (g *Game) AddPlayer (user *tgbotapi.User) (error){
	present,_ := g.FindPlayer(user)
	if !present {
		g.Players = append(g.Players, user)
		return nil
	}
	return errors.New("player already in game")
}


func (g *Game) RemovePlayer (user *tgbotapi.User) (error){
	present,idx := g.FindPlayer(user)
	currentPlayers := g.Players

	if present  {
		currentPlayers = append(currentPlayers[:idx], currentPlayers[idx+1:]...)
		return nil
	}
	return errors.New("player not in game")
}









