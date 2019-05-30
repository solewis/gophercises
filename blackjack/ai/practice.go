package ai

import (
	"fmt"
	blackjack "gophercises/blackjack/game"
	"gophercises/blackjack/money"
	"gophercises/carddeck"
)

func Practice() blackjack.AI {
	return practice{}
}

type practice struct{}

func (_ practice) Name() string {
	return "Practice Human AI"
}

func (_ practice) Play(hand []deck.Card, dealerShowing deck.Card, allowedMoves []blackjack.Move) blackjack.Move {
	fmt.Println("Dealer:", dealerShowing.String())
	handScore, soft := blackjack.Score(hand)
	softString := ""
	if soft {
		softString = "soft "
	}
	fmt.Printf("You: %s (%s%d)\n", blackjack.HandString(hand), softString, handScore)
	for {
		fmt.Println("Would you like to:")
		for i, m := range allowedMoves {
			fmt.Printf("%d: %v\n", i+1, m)
		}

		var playerChoice int
		for {
			fmt.Scanln(&playerChoice)
			switch {
			case playerChoice > 0 && playerChoice <= len(allowedMoves):
				move := allowedMoves[playerChoice-1]
				pPair := len(hand) == 2 && hand[0].Rank == hand[1].Rank
				dScore, _ := blackjack.Score([]deck.Card{dealerShowing})
				correctMove := lookup[playerHand{handScore, soft, pPair}][dScore](allowedMoves)
				if move == correctMove {
					fmt.Println("Correct!")
				} else {
					fmt.Println("Incorrect... correct move was", correctMove.String())
				}
				return move
			default:
				fmt.Println("Invalid chooseF.")
			}
		}
	}
}

func (_ practice) Bet(minBet money.USD, maxBet money.USD, _ bool) money.USD {
	var betInput float64
	fmt.Printf("How much would you like to bet? %s (min) %s (max)\n", minBet.String(), maxBet.String())
	for {
		_, e := fmt.Scanf("%f\n", &betInput)
		switch {
		case e != nil:
			fmt.Println("Invalid bet... Must be a float with a max of 2 decimal places")
		case betInput < minBet.Float64():
			fmt.Println("Bet must be greater than or equal to", minBet.String())
		case betInput > maxBet.Float64():
			fmt.Println("Bet must be less than or equal to", maxBet.String())
		default:
			bet := money.ToUSD(betInput)
			fmt.Printf("Your bet is %s\n", bet.String())
			return bet
		}
	}
}

func (_ practice) HandResults(hand, dealer []deck.Card, winnings, balance money.USD) {
	fmt.Println("---Hand Results---")
	dealerScore, _ := blackjack.Score(dealer)
	fmt.Printf("Dealer: %s (%d)\n", blackjack.HandString(dealer), dealerScore)
	handScore, _ := blackjack.Score(hand)
	fmt.Printf("You: %s (%d)\n", blackjack.HandString(hand), handScore)

	switch {
	case winnings.Float64() == 0:
		fmt.Println("No money gained or lost. Balance:", balance.String())
	case winnings.Float64() < 0:
		fmt.Printf("You lost %s. Balance: %s\n", money.ToUSD(winnings.Float64() * -1).String(), balance.String())
	default:
		fmt.Printf("You won %s. Balance: %s\n", winnings.String(), balance.String())
	}
}

func (_ practice) RoundRecap(allHands [][]deck.Card) {}
