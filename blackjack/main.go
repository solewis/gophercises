package main

//TODO
// basic AI
// Smart AI
// Card counting AI
// allow user to pick rules
// calculate house edge based on rules
// multiple players
import (
	"fmt"
	"gophercises/blackjack/ai"
	blackjack "gophercises/blackjack/game"
	"gophercises/blackjack/money"
)

const (
	minBet                     money.USD = 200
	maxBet                     money.USD = 50000
	naturalBlackjackMultiplier float64   = 1.5
	numDecks                   int       = 4
	numRounds                  int       = 10
	percentDeckUsage           float64   = .75
)

func main() {
	fmt.Println("Welcome to blackjack")
	fmt.Println("Table rules:")
	fmt.Println("1. Dealer hits on soft 17")
	fmt.Printf("2. Minimum bet %s\n", minBet.String())
	fmt.Printf("3. Maximum bet %s\n", maxBet.String())
	fmt.Println("4. Dealer does not automatically win on blackjack")
	fmt.Printf("5. Payout on natural blackjack %.2f\n", naturalBlackjackMultiplier)
	fmt.Println("6. No double down restrictions (can double down on any hand")
	fmt.Printf("7. %d decks\n", numDecks)
	//fmt.Println("8. No split restrictions")
	//fmt.Println("9. Early and late surrender allowed")
	//fmt.Println("10. Continuous shuffle?")

	blackjack.Play(blackjack.Options{
		MinBet:                     minBet,
		MaxBet:                     maxBet,
		NaturalBlackjackMultiplier: naturalBlackjackMultiplier,
		NumDecks:                   numDecks,
		NumRounds:                  numRounds,
		PercentDeckUsage:           percentDeckUsage,
		AI:                         ai.Human(),
	})
}
