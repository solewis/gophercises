package main

//TODO
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

	rounds := determineRounds()
	chosenAI := determineAI()

	finalBalance := blackjack.Play(blackjack.Options{
		MinBet:                     minBet,
		MaxBet:                     maxBet,
		NaturalBlackjackMultiplier: naturalBlackjackMultiplier,
		NumDecks:                   numDecks,
		NumRounds:                  rounds,
		PercentDeckUsage:           percentDeckUsage,
		AI:                         chosenAI,
	})

	fmt.Println("------FINAL BALANCE------")
	fmt.Println(finalBalance.String())
}

func determineRounds() int {
	fmt.Print("How many hands would you like to play? ")
	var rounds int
	for {
		fmt.Scanln(&rounds)
		switch {
		case rounds < 0:
			fmt.Println("Rounds must be more than 0")
		default:
			return rounds
		}
	}
}

func determineAI() blackjack.AI {
	fmt.Println("Which AI would you like to use?")
	ais := []blackjack.AI{ai.Human(), ai.Basic()}
	for i, ai := range ais {
		fmt.Printf("%d: %s\n", i+1, ai.Name())
	}
	var playerChoice int
	for {
		fmt.Scanln(&playerChoice)
		switch {
		case playerChoice > 0 && playerChoice <= len(ais):
			return ais[playerChoice-1]
		default:
			fmt.Println("Invalid choice.")
		}
	}
}
