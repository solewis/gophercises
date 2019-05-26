package ai

import (
	blackjack "gophercises/blackjack/game"
	"gophercises/blackjack/money"
	deck "gophercises/carddeck"
)

func Smart() blackjack.AI {
	return smart{}
}

type smart struct{}

func (_ smart) Name() string {
	return "Smart AI"
}

func (_ smart) Play(hand []deck.Card, dealerShowing deck.Card, allowedMoves []blackjack.Move) blackjack.Move {
	pScore, soft := blackjack.Score(hand)
	pPair := len(hand) == 2 && hand[0].Rank == hand[1].Rank
	dScore, _ := blackjack.Score([]deck.Card{dealerShowing})
	return lookup[playerHand{pScore, soft, pPair}][dScore](allowedMoves)
}

func (_ smart) Bet(minBet money.USD, maxBet money.USD) money.USD {
	return money.ToUSD(5)
}

func (_ smart) Results(hand [][]deck.Card, dealer []deck.Card, winnings, balance money.USD) {}

func containsMove(moves []blackjack.Move, move blackjack.Move) bool {
	for _, m := range moves {
		if m == move {
			return true
		}
	}
	return false
}
