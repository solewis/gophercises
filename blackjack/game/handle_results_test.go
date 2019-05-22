package blackjack

import (
	"gophercises/blackjack/money"
	"gophercises/carddeck"
	"reflect"
	"testing"
)

var testOptions = Options{
	NaturalBlackjackMultiplier: 1.5,
}

func TestHandleResults_PlayerHandSurrendered(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Two, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Ace, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := TestAI{}
	testState := setupSingleDeckGameSurrendered(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)
	verify(t, ai, -100, 900, testState)
}

func TestHandleResults_BothNaturalBlackjack(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Ace, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Ace, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := TestAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)
	verify(t, ai, 0, 1000, testState)
}

func TestHandleResults_PlayerNaturalBlackjackOnly(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Eight, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Ace, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := TestAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)
	verify(t, ai, 300, 1300, testState)
}

func TestHandleResults_DealerNaturalBlackjackOnly(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Ace, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Eight, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := TestAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)
	verify(t, ai, -200, 800, testState)
}

func TestHandleResults_PlayerBust(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Heart}}
	playerHand1 := []deck.Card{{Rank: deck.Eight, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}, {Rank: deck.Five, Suit: deck.Club}}
	ai := TestAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)
	verify(t, ai, -200, 800, testState)
}

func TestHandleResults_DealerBust(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Heart}}
	playerHand1 := []deck.Card{{Rank: deck.Three, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := TestAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)
	verify(t, ai, 200, 1200, testState)
}

func TestHandleResults_PlayerLessThanDealer(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Nine, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := TestAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)
	verify(t, ai, -200, 800, testState)
}

func TestHandleResults_PlayerMoreThanDealer(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Ten, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := TestAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)
	verify(t, ai, 200, 1200, testState)
}

func TestHandleResults_PlayerEqualToDealer(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Nine, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := TestAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)
	verify(t, ai, 0, 1000, testState)
}

func TestHandleResults_MultiplePlayerHands(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Seven, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Eight, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	playerHand2 := []deck.Card{{Rank: deck.Six, Suit: deck.Heart}, {Rank: deck.Ten, Suit: deck.Heart}}
	ai := &TestAI{}
	testPlayer := player{
		hands: []hand{
			{
				cards: playerHand1,
				bet:   400,
			},
			{
				cards: playerHand2,
				bet:   200,
			},
		},
		balance: 1000,
		ai:      ai,
	}
	testState := state{
		player:     testPlayer,
		dealerHand: testDealerHand,
	}
	handleResults(&testState, testOptions)
	verify(t, *ai, 200, 1200, testState)
}

func setupSingleDeckGameSurrendered(playerHand, dealerHand []deck.Card, ai AI) state {
	testState := setupSingleDeckGame(playerHand, dealerHand, ai)
	testState.player.hands[0].surrendered = true
	return testState
}

func setupSingleDeckGame(playerHand, dealerHand []deck.Card, ai AI) state {
	testPlayer := player{
		hands: []hand{
			{
				cards: playerHand,
				bet:   200,
			},
		},
		balance: 1000,
		ai:      ai,
	}
	testState := state{
		player:     testPlayer,
		dealerHand: dealerHand,
	}
	return testState
}

func verify(
	t *testing.T,
	ai TestAI,
	expectedWinnings, expectedBalance money.USD,
	testState state) {

	var playerHands [][]deck.Card
	for _, hand := range testState.player.hands {
		playerHands = append(playerHands, hand.cards)
	}

	if testState.player.balance != expectedBalance {
		t.Errorf("Expected balance to be %s, but was %s", expectedBalance.String(), testState.player.balance.String())
	}
	if !reflect.DeepEqual(ai.sentResultHands, playerHands) {
		t.Errorf("Expected player hands passed to AI to be %v but was %v", playerHands, ai.sentResultHands)
	}
	if !reflect.DeepEqual(ai.sentDealerHand, testState.dealerHand) {
		t.Errorf("Expected dealer hand passed to AI to be %v but was %v", testState.dealerHand, ai.sentDealerHand)
	}
	if ai.sentWinnings != expectedWinnings {
		t.Errorf("Expected winnings to be %s but was %s", expectedWinnings, ai.sentWinnings)
	}
}
