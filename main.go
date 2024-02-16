package main

import (
	"log"
	"os"
	"math/rand"
	"time"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

// Variable initialization
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

type GameSession struct {
	Players []tgbotapi.User
	Deck    Deck
	Hands   []Hand
}

var gameInProgress bool
var session GameSession

//Shuffle cards
func Shuffle (deck Deck){
	rand.Seed(time.Now().UnixNano()) //Create a RNG seed using the current timestamp 
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] }) //Fisher-Yates algo to shuffle
} 

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the Telegram bot token from the environment variable
	botToken := os.Getenv("TELEGRAM_APITOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_APITOKEN not found in environment variables")
	}

	// Create a new bot instance
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// Set up an update configuration
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Get updates from Telegram
	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Process incoming messages
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		switch update.Message.Text {
		case "/start":
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
			bot.Send(msg)

		case "/help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This is the help message.")
			bot.Send(msg)

		case "/play_game":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Starting the game...")
			bot.Send(msg)


		case "/leave":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Leaving...")
			bot.Send(msg)

		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Please use /start to see available options.")
			bot.Send(msg)
		}
	}
}
