package entities
import(
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
type Hand struct{
	player *tgbotapi.User
	Cards []Card
}

func NewHand (player *tgbotapi.User, card []Card) *Hand{
	return &Hand{
		player:player,
		Cards:card,
	}
}

// func (h *Hand) SortSuit(hand Hand) map[string][]Card {
// 	sortedCards := make(map[string][]Card)

// 	for _, card := range hand.Cards {
// 		// Check if the suit is already a key in the map
// 		if _, ok := sortedCards[string(card.Suit)]; !ok {
// 			// If not, initialize a slice for the suit
// 			sortedCards[string(card.Suit)] = make([]Card, 0)
// 		}

// 		// Append the card to the corresponding suit in the map
// 		sortedCards[string(card.Suit)] = append(sortedCards[string(card.Suit)], card)
// 	}

// 	for _, key := range sortedCards
// }
