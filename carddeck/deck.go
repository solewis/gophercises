package deck

import (
	"fmt"
	"math/rand"
	"time"
)

type Suit int

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

func (c Suit) String() string {
	names := []string{"Spade", "Diamond", "Club", "Heart"}
	return names[c]
}

type Rank int

const (
	Ace Rank = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

var ranks = [...]Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

func (r Rank) String() string {
	names := []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}
	return names[r]
}

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

type DeckOption func(cards []Card) []Card

func New(deckOptions ...DeckOption) []Card {
	var deck []Card
	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, Card{rank, suit})
		}
	}

	for _, opt := range deckOptions {
		deck = opt(deck)
	}
	return deck
}

func Expand(n int) DeckOption {
	return func(cards []Card) []Card {
		var expanded []Card
		for i := 0; i < n; i++ {
			expanded = append(expanded, cards...)
		}
		return expanded
	}
}

func Shuffle(cards []Card) []Card {
	shuffledCards := make([]Card, len(cards))
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i, r := range r.Perm(len(cards)) {
		shuffledCards[i] = cards[r]
	}
	return shuffledCards
}
