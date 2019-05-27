package blackjack

import (
	"gophercises/blackjack/money"
	deck "gophercises/carddeck"
	"reflect"
	"testing"
)

var testOptions = Options{
	NaturalBlackjackMultiplier: 1.5,
}
//TODO
// test multiple players
func TestHandleResults_PlayerHandSurrendered(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Two, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Ace, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := testAI{}
	testState := setupSingleDeckGameSurrendered(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)

	expectBalance(t, testState.players[0].balance, 900)
	expectAIResultHands(t, ai.sentResultHands, [][]deck.Card{playerHand1})
	expectDealerHand(t, ai.sentDealerHand, testDealerHand)
	expectWinnings(t, ai.sentWinnings, -100)
}

func TestHandleResults_BothNaturalBlackjack(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Ace, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Ace, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := testAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)

	expectBalance(t, testState.players[0].balance, 1000)
	expectAIResultHands(t, ai.sentResultHands, [][]deck.Card{playerHand1})
	expectDealerHand(t, ai.sentDealerHand, testDealerHand)
	expectWinnings(t, ai.sentWinnings, 0)
}

func TestHandleResults_PlayerNaturalBlackjackOnly(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Eight, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Ace, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := testAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)

	expectBalance(t, testState.players[0].balance, 1300)
	expectAIResultHands(t, ai.sentResultHands, [][]deck.Card{playerHand1})
	expectDealerHand(t, ai.sentDealerHand, testDealerHand)
	expectWinnings(t, ai.sentWinnings, 300)
}

func TestHandleResults_DealerNaturalBlackjackOnly(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Ace, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Eight, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := testAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)

	expectBalance(t, testState.players[0].balance, 800)
	expectAIResultHands(t, ai.sentResultHands, [][]deck.Card{playerHand1})
	expectDealerHand(t, ai.sentDealerHand, testDealerHand)
	expectWinnings(t, ai.sentWinnings, -200)
}

func TestHandleResults_PlayerBust(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Heart}}
	playerHand1 := []deck.Card{{Rank: deck.Eight, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}, {Rank: deck.Five, Suit: deck.Club}}
	ai := testAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)

	expectBalance(t, testState.players[0].balance, 800)
	expectAIResultHands(t, ai.sentResultHands, [][]deck.Card{playerHand1})
	expectDealerHand(t, ai.sentDealerHand, testDealerHand)
	expectWinnings(t, ai.sentWinnings, -200)
}

func TestHandleResults_DealerBust(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Heart}}
	playerHand1 := []deck.Card{{Rank: deck.Three, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := testAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)

	expectBalance(t, testState.players[0].balance, 1200)
	expectAIResultHands(t, ai.sentResultHands, [][]deck.Card{playerHand1})
	expectDealerHand(t, ai.sentDealerHand, testDealerHand)
	expectWinnings(t, ai.sentWinnings, 200)
}

func TestHandleResults_PlayerLessThanDealer(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Nine, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := testAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)

	expectBalance(t, testState.players[0].balance, 800)
	expectAIResultHands(t, ai.sentResultHands, [][]deck.Card{playerHand1})
	expectDealerHand(t, ai.sentDealerHand, testDealerHand)
	expectWinnings(t, ai.sentWinnings, -200)
}

func TestHandleResults_PlayerMoreThanDealer(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Ten, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := testAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)

	expectBalance(t, testState.players[0].balance, 1200)
	expectAIResultHands(t, ai.sentResultHands, [][]deck.Card{playerHand1})
	expectDealerHand(t, ai.sentDealerHand, testDealerHand)
	expectWinnings(t, ai.sentWinnings, 200)
}

func TestHandleResults_PlayerEqualToDealer(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Nine, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := testAI{}
	testState := setupSingleDeckGame(playerHand1, testDealerHand, &ai)
	handleResults(&testState, testOptions)

	expectBalance(t, testState.players[0].balance, 1000)
	expectAIResultHands(t, ai.sentResultHands, [][]deck.Card{playerHand1})
	expectDealerHand(t, ai.sentDealerHand, testDealerHand)
	expectWinnings(t, ai.sentWinnings, 0)
}

func TestHandleResults_MultiplePlayerHands(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Seven, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Eight, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	playerHand2 := []deck.Card{{Rank: deck.Six, Suit: deck.Heart}, {Rank: deck.Ten, Suit: deck.Heart}}
	ai := &testAI{}
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
		players:    []player{testPlayer},
		dealerHand: testDealerHand,
	}
	handleResults(&testState, testOptions)

	expectBalance(t, testState.players[0].balance, 1200)
	expectAIResultHands(t, ai.sentResultHands, [][]deck.Card{playerHand1, playerHand2})
	expectDealerHand(t, ai.sentDealerHand, testDealerHand)
	expectWinnings(t, ai.sentWinnings, 200)
}

func TestHandleResults_MultiplePlayers(t *testing.T) {
	testDealerHand := []deck.Card{{Rank: deck.Seven, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand1 := []deck.Card{{Rank: deck.Eight, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	playerHand2 := []deck.Card{{Rank: deck.Six, Suit: deck.Heart}, {Rank: deck.Ten, Suit: deck.Heart}}
	ai1 := &testAI{}
	ai2 := &testAI{}
	testPlayer1 := player{
		hands: []hand{
			{
				cards: playerHand1,
				bet:   400,
			},
		},
		balance: 1000,
		ai:      ai1,
	}
	testPlayer2 := player{
		hands: []hand{
			{
				cards: playerHand2,
				bet:   200,
			},
		},
		balance: 1000,
		ai:      ai2,
	}
	testState := state{
		players:    []player{testPlayer1, testPlayer2},
		dealerHand: testDealerHand,
	}
	handleResults(&testState, testOptions)

	//player 1
	expectBalance(t, testState.players[0].balance, 1400)
	expectWinnings(t, ai1.sentWinnings, 400)
	expectAIResultHands(t, ai1.sentResultHands, [][]deck.Card{playerHand1})
	expectDealerHand(t, ai1.sentDealerHand, testDealerHand)

	//player 2
	expectBalance(t, testState.players[1].balance, 800)
	expectWinnings(t, ai2.sentWinnings, -200)
	expectAIResultHands(t, ai2.sentResultHands, [][]deck.Card{playerHand2})
	expectDealerHand(t, ai2.sentDealerHand, testDealerHand)
}

func setupSingleDeckGameSurrendered(playerHand, dealerHand []deck.Card, ai AI) state {
	testState := setupSingleDeckGame(playerHand, dealerHand, ai)
	testState.players[0].hands[0].surrendered = true
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
		players:    []player{testPlayer},
		dealerHand: dealerHand,
	}
	return testState
}

func expectBalance(t *testing.T, balance, expectedBalance money.USD) {
	if balance != expectedBalance {
		t.Errorf("Expected balance %s but was %s", expectedBalance.String(), balance.String())
	}
}

func expectAIResultHands(t *testing.T, hands, expectedHands [][]deck.Card) {
	if !reflect.DeepEqual(hands, expectedHands) {
		t.Errorf("Expected %v hands to be sent to AI, but was %v", expectedHands, hands)
	}
}

func expectDealerHand(t *testing.T, hand, expectedHand []deck.Card) {
	if !reflect.DeepEqual(hand, expectedHand) {
		t.Errorf("Expected %v dealer hand to be sent to AI, but was %v", expectedHand, hand)
	}
}

func expectWinnings(t *testing.T, winnings, expectedWinnings money.USD) {
	if winnings != expectedWinnings {
		t.Errorf("Expected winnings %s but was %s", expectedWinnings.String(), winnings.String())
	}
}
