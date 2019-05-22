package blackjack

import (
	"gophercises/blackjack/money"
	deck "gophercises/carddeck"
)

type TestAI struct {
	returnedMove              Move
	sentAllowedMoves          []Move
	sentPlayHand              []deck.Card
	sentDealerShowing         deck.Card
	sentResultHands           [][]deck.Card
	sentDealerHand            []deck.Card
	sentWinnings, sentBalance money.USD
}

func (ai *TestAI) Play(hand []deck.Card, dealerShowing deck.Card, allowedMoves []Move) Move {
	ai.sentPlayHand = hand
	ai.sentDealerShowing = dealerShowing
	ai.sentAllowedMoves = allowedMoves
	return ai.returnedMove
}

func (ai *TestAI) Bet(minBet money.USD, maxBet money.USD) money.USD {
	return 0
}

func (ai *TestAI) Results(hand [][]deck.Card, dealer []deck.Card, winnings, balance money.USD) {
	ai.sentResultHands = hand
	ai.sentDealerHand = dealer
	ai.sentWinnings = winnings
	ai.sentBalance = balance
}

type CheatingAI struct {}

func (ai CheatingAI) Play(hand []deck.Card, dealerShowing deck.Card, allowedMoves []Move) Move {
	hand[0] = deck.Card{Rank: deck.Ace, Suit: deck.Heart}
	hand[1] = deck.Card{Rank: deck.Ten, Suit: deck.Diamond}
	return Hit
}

func (ai CheatingAI) Bet(minBet money.USD, maxBet money.USD) money.USD {
	return 0
}

func (ai CheatingAI) Results(hand [][]deck.Card, dealer []deck.Card, winnings, balance money.USD) {}
