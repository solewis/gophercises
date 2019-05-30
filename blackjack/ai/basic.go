package ai

import (
	blackjack "gophercises/blackjack/game"
	"gophercises/blackjack/money"
	deck "gophercises/carddeck"
)

func Basic() blackjack.AI {
	return basic{}
}

type basic struct{}

func (_ basic) Name() string {
	return "Basic AI"
}

func (_ basic) Play(hand []deck.Card, dealerShowing deck.Card, allowedMoves []blackjack.Move) blackjack.Move {
	eights := hand[0].Rank == deck.Eight && hand[1].Rank == deck.Eight && len(hand) == 2
	aces := hand[0].Rank == deck.Ace && hand[1].Rank == deck.Ace && len(hand) == 2
	if eights || aces {
		return blackjack.Split
	}
	handScore, _ := blackjack.Score(hand)
	dealerScore, _ := blackjack.Score([]deck.Card{dealerShowing})
	if (handScore == 10 || handScore == 11) && dealerScore < 10 && len(hand) == 2 {
		return blackjack.Double
	}
	if dealerScore >= 2 && dealerScore <= 3 && handScore >= 13 {
		return blackjack.Stand
	}
	if dealerScore >= 4 && dealerScore <= 6 && handScore >= 12 {
		return blackjack.Stand
	}
	if dealerScore >= 7 && handScore >= 17 {
		return blackjack.Stand
	}
	return blackjack.Hit
}

func (_ basic) Bet(minBet money.USD, maxBet money.USD, _ bool) money.USD {
	return money.ToUSD(5)
}

func (_ basic) HandResults(hand, dealer []deck.Card, winnings, balance money.USD) {}

func (_ basic) RoundRecap(allHands [][]deck.Card) {}
