package main

//TODO
// Card counting AI
// allow user to pick rules
// calculate house edge based on rules
// strategy morph to rules (H17 vs S17, etc)
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
	fmt.Println("7. Late surrender allowed")
	fmt.Printf("8. %d decks\n", numDecks)

	rounds := determineRounds()
	ais := determineAIs()

	fmt.Println("------STARTING GAME------")

	balances := blackjack.Play(blackjack.Options{
		MinBet:                     minBet,
		MaxBet:                     maxBet,
		NaturalBlackjackMultiplier: naturalBlackjackMultiplier,
		NumDecks:                   numDecks,
		NumRounds:                  rounds,
		PercentDeckUsage:           percentDeckUsage,
		AIs:                        ais,
	})

	fmt.Println("------FINAL BALANCES------")
	for k, v := range balances {
		fmt.Printf("%s balance: %s\n", k, v.String())
	}
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

func determineAIs() map[string]blackjack.AI {
	fmt.Print("How many players? (1-8) ")
	var players int
	for {
		fmt.Scanln(&players)
		if players > 0 && players <= 8 {
			break
		} else {
			fmt.Println("Invalid choice. Choose a number from 1 to 8")
		}
	}

	ais := []blackjack.AI{ai.Human(), ai.Practice(), ai.Basic(), ai.Smart()}
	selectedAIs := make(map[string]blackjack.AI)
	for i := 0; i < players; i++ {
		fmt.Println("Which AI would you like to use for player", i + 1)
		for aiIdx, a := range ais {
			fmt.Printf("%d: %s\n", aiIdx+1, a.Name())
		}
		var aiChoice int
		for {
			fmt.Scanln(&aiChoice)
			if aiChoice > 0 && aiChoice <= len(ais) {
				a := ais[aiChoice-1]
				selectedAIs[a.Name()] = a
				break
			} else {
				fmt.Println("Invalid choice.")
			}
		}
	}
	return selectedAIs
}
