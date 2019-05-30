package blackjack

import (
	"gophercises/blackjack/money"
	deck "gophercises/carddeck"
)

type testAI struct {
	returnedMove              Move
	sentAllowedMoves          []Move
	sentPlayHand              []deck.Card
	sentDealerShowing         deck.Card
	sentResultHands           [][]deck.Card
	sentDealerHand            []deck.Card
	sentWinnings, sentBalance money.USD
}

func (_ *testAI) Name() string {
	return "Test AI"
}

func (ai *testAI) Play(hand []deck.Card, dealerShowing deck.Card, allowedMoves []Move) Move {
	ai.sentPlayHand = hand
	ai.sentDealerShowing = dealerShowing
	ai.sentAllowedMoves = allowedMoves
	return ai.returnedMove
}

func (_ *testAI) Bet(minBet money.USD, maxBet money.USD) money.USD {
	return 0
}

func (ai *testAI) HandResults(hand, dealer []deck.Card, winnings, balance money.USD) {
	//ai.sentResultHands = hand
	ai.sentDealerHand = dealer
	ai.sentWinnings = winnings
	ai.sentBalance = balance
}

func (_ testAI) RoundRecap(allHands [][]deck.Card) {}

type cheatingAI struct {}

func (h cheatingAI) Name() string {
	return "Cheating AI"
}

func (_ cheatingAI) Play(hand []deck.Card, dealerShowing deck.Card, allowedMoves []Move) Move {
	hand[0] = deck.Card{Rank: deck.Ace, Suit: deck.Heart}
	hand[1] = deck.Card{Rank: deck.Ten, Suit: deck.Diamond}
	return Hit
}

func (_ cheatingAI) Bet(minBet money.USD, maxBet money.USD) money.USD {
	return 0
}

func (_ cheatingAI) HandResults(hand, dealer []deck.Card, winnings, balance money.USD) {}

func (_ cheatingAI) RoundRecap(allHands [][]deck.Card) {}
