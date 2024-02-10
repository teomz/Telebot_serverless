package main

import (
	"log"
	"math/rand"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Constants and type definitions for Card, Deck, etc.
type Suit string
type Rank int

const (
	Spades   Suit = "Spades"
	Hearts   Suit = "Hearts"
	Diamonds Suit = "Diamonds"
	Clubs    Suit = "Clubs"
)

const (
	Ace   Rank = 14
	King  Rank = 13
	Queen Rank = 12
	Jack  Rank = 11
	Ten   Rank = 10
	Nine  Rank = 9
	Eight Rank = 8
	Seven Rank = 7
	Six   Rank = 6
	Five  Rank = 5
	Four  Rank = 4
	Three Rank = 3
	Two   Rank = 2
)

type Card struct {
	Suit Suit
	Rank Rank
}

type Deck []Card
type Hand []Card

func NewDeck() Deck {
	ranks := []Rank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
	suits := []Suit{Spades, Hearts, Diamonds, Clubs}

	var deck Deck
	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Card{Suit: suit, Rank: rank})
		}
	}
	return deck
}

func Shuffle(deck Deck) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
}

func Deal(deck Deck, numPlayers int) []Hand {
	hands := make([]Hand, numPlayers)
	for i := 0; i < numPlayers; i++ {
		hands[i] = make(Hand, 0, len(deck)/numPlayers)
	}

	for i, card := range deck {
		hands[i%numPlayers] = append(hands[i%numPlayers], card)
	}

	return hands
}

func main() {
	bot, err := tgbotapi.NewBotAPI("6863492345:AAH-ak_depbfolBuCoI7PzfHu4ajJZ0L030") // Replace with your Bot Token
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Error getting updates: %v", err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			callbackData := update.CallbackQuery.Data
			if callbackData == "play_game" {
				// Game logic...
				log.Printf("Game started by user: %s", update.CallbackQuery.From.UserName)

				// Here you could shuffle and deal cards, then notify players
				deck := NewDeck()
				Shuffle(deck)
				hands := Deal(deck, 4) // Assuming 4 players
				for i, hand := range hands {
					log.Printf("Player %d's hand: %v\n", i+1, hand) // Simplified
				}

				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "The game has started. Check the console for your hand.")
				bot.Send(msg)

				callbackResp := tgbotapi.NewCallback(update.CallbackQuery.ID, "Game started!")
				bot.AnswerCallbackQuery(callbackResp)
			}
		} else if update.Message != nil {
			command := strings.Split(update.Message.Text, " ")[0]

			switch command {
			case "/start":
				msgText := "Welcome to the Bridge bot! Use /help to see available commands."
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				bot.Send(msg)
			case "/help":
				msgText := "Available commands:\n" +
					"/start - Start interacting with the bot\n" +
					"/help - Get help and see available commands\n" +
					"/play_game - Start a new game of Bridge"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				bot.Send(msg)
			case "/play_game":
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🎲 Play Bridge", "play_game"),
					),
				)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Press 'Play Bridge' to start the game.")
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
			default:
				msgText := "Sorry, I didn't recognize that command. Use /help to see available commands."
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				bot.Send(msg)
			}
		}
	}
}
