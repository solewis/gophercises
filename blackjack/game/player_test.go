package blackjack

import (
	"gophercises/blackjack/money"
	deck "gophercises/carddeck"
	"reflect"
	"testing"
)

func TestPlayer_Hit(t *testing.T) {
	dealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand := []deck.Card{{Rank: deck.Two, Suit: deck.Club}, {Rank: deck.Two, Suit: deck.Diamond}}
	ai := testAI{returnedMove: Hit}
	testState := buildSinglePlayerGame(playerHand, dealerHand, &ai)
	nextCard := testState.deck.cards[0]

	finished := runPlayerHand(&testState.players[0], &testState.players[0].hands[0], dealerHand[0], &testState.deck, true)
	expectedHand := append(playerHand, nextCard)

	expectHand(t, testState.players[0].hands[0].cards, expectedHand)
	expectBet(t, testState.players[0].hands[0].bet, 200)
	expectFinished(t, finished, false)
	expectAIPlayHand(t, ai.sentPlayHand, playerHand)
	expectAIDealerCard(t, ai.sentDealerShowing, dealerHand[0])
}

func TestPlayer_HitAndBust(t *testing.T) {
	dealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := testAI{returnedMove: Hit}
	testState := buildSinglePlayerGame(playerHand, dealerHand, &ai)
	nextCard := testState.deck.cards[0]

	finished := runPlayerHand(&testState.players[0], &testState.players[0].hands[0], dealerHand[0], &testState.deck, true)
	expectedHand := append(playerHand, nextCard)

	expectHand(t, testState.players[0].hands[0].cards, expectedHand)
	expectBet(t, testState.players[0].hands[0].bet, 200)
	expectFinished(t, finished, true)
	expectAIPlayHand(t, ai.sentPlayHand, playerHand)
	expectAIDealerCard(t, ai.sentDealerShowing, dealerHand[0])
}

func TestPlayer_Split(t *testing.T) {
	dealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand := []deck.Card{{Rank: deck.Eight, Suit: deck.Club}, {Rank: deck.Eight, Suit: deck.Diamond}}
	ai := testAI{returnedMove: Split}
	testState := buildSinglePlayerGame(playerHand, dealerHand, &ai)
	nextCards := testState.deck.cards[0:2]

	finished := runPlayerHand(&testState.players[0], &testState.players[0].hands[0], dealerHand[0], &testState.deck, true)
	expectedHand1 := []deck.Card{playerHand[0], nextCards[0]}
	expectedHand2 := []deck.Card{playerHand[1], nextCards[1]}

	expectHand(t, testState.players[0].hands[0].cards, expectedHand1)
	expectHand(t, testState.players[0].hands[1].cards, expectedHand2)
	expectBet(t, testState.players[0].hands[0].bet, 200)
	expectBet(t, testState.players[0].hands[1].bet, 200)
	expectFinished(t, finished, false)
	if testState.players[0].hands[1].split != true {
		t.Error("Expected second hand to indicate it was split")
	}
}

func TestPlayer_CannotSplitIfHandNotTwoCards(t *testing.T) {
	dealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand := []deck.Card{{Rank: deck.Five, Suit: deck.Club}, {Rank: deck.Five, Suit: deck.Heart}, {Rank: deck.Six, Suit: deck.Diamond}}
	ai := testAI{returnedMove: Split}
	testState := buildSinglePlayerGame(playerHand, dealerHand, &ai)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected split to fail")
		}
	}()

	runPlayerHand(&testState.players[0], &testState.players[0].hands[0], dealerHand[0], &testState.deck, true)
	if containsMove(ai.sentAllowedMoves, Split) {
		t.Error("Split was passed to AI as an allowed move when it should not have been")
	}
}

func TestPlayer_CannotSplitIfCardsNotEqualRank(t *testing.T) {
	dealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand := []deck.Card{{Rank: deck.Five, Suit: deck.Club}, {Rank: deck.Six, Suit: deck.Diamond}}
	ai := testAI{returnedMove: Split}
	testState := buildSinglePlayerGame(playerHand, dealerHand, &ai)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected split to fail")
		}
	}()

	runPlayerHand(&testState.players[0], &testState.players[0].hands[0], dealerHand[0], &testState.deck, true)
	if containsMove(ai.sentAllowedMoves, Split) {
		t.Error("Split was passed to AI as an allowed move when it should not have been")
	}
}

func TestPlayer_Double(t *testing.T) {
	dealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Club}, {Rank: deck.Two, Suit: deck.Diamond}}
	ai := testAI{returnedMove: Double}
	testState := buildSinglePlayerGame(playerHand, dealerHand, &ai)
	nextCard := testState.deck.cards[0]

	finished := runPlayerHand(&testState.players[0], &testState.players[0].hands[0], dealerHand[0], &testState.deck, true)
	expectedHand := append(playerHand, nextCard)
	expectedBet := testState.players[0].initialBet.Multiply(2)

	expectHand(t, testState.players[0].hands[0].cards, expectedHand)
	expectBet(t, testState.players[0].hands[0].bet, expectedBet)
	expectFinished(t, finished, true)
	expectAIPlayHand(t, ai.sentPlayHand, playerHand)
	expectAIDealerCard(t, ai.sentDealerShowing, dealerHand[0])
}

func TestPlayer_CannotDoubleIfHandNotTwoCards(t *testing.T) {
	dealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand := []deck.Card{{Rank: deck.Five, Suit: deck.Club}, {Rank: deck.Five, Suit: deck.Heart}, {Rank: deck.Six, Suit: deck.Diamond}}
	ai := testAI{returnedMove: Double}
	testState := buildSinglePlayerGame(playerHand, dealerHand, &ai)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected double to fail")
		}
	}()

	runPlayerHand(&testState.players[0], &testState.players[0].hands[0], dealerHand[0], &testState.deck, true)
	if containsMove(ai.sentAllowedMoves, Double) {
		t.Error("Double was passed to AI as an allowed move when it should not have been")
	}
}

func TestPlayer_Stand(t *testing.T) {
	dealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
	ai := testAI{returnedMove: Stand}
	testState := buildSinglePlayerGame(playerHand, dealerHand, &ai)

	finished := runPlayerHand(&testState.players[0], &testState.players[0].hands[0], dealerHand[0], &testState.deck, true)

	expectHand(t, testState.players[0].hands[0].cards, playerHand)
	expectBet(t, testState.players[0].hands[0].bet, 200)
	expectFinished(t, finished, true)
	expectAIPlayHand(t, ai.sentPlayHand, playerHand)
	expectAIDealerCard(t, ai.sentDealerShowing, dealerHand[0])
}

func TestPlayer_Surrender(t *testing.T) {
	dealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Club}, {Rank: deck.Six, Suit: deck.Diamond}}
	ai := testAI{returnedMove: Surrender}
	testState := buildSinglePlayerGame(playerHand, dealerHand, &ai)

	finished := runPlayerHand(&testState.players[0], &testState.players[0].hands[0], dealerHand[0], &testState.deck, true)

	expectHand(t, testState.players[0].hands[0].cards, playerHand)
	expectBet(t, testState.players[0].hands[0].bet, 200)
	expectFinished(t, finished, true)
	expectAIPlayHand(t, ai.sentPlayHand, playerHand)
	expectAIDealerCard(t, ai.sentDealerShowing, dealerHand[0])
	if !testState.players[0].hands[0].surrendered {
		t.Error("Expected hand to be surrendered")
	}
}

func TestPlayer_CannotSurrenderIfNotTwoCards(t *testing.T) {
	dealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand := []deck.Card{{Rank: deck.Five, Suit: deck.Club}, {Rank: deck.Five, Suit: deck.Heart}, {Rank: deck.Six, Suit: deck.Diamond}}
	ai := testAI{returnedMove: Surrender}
	testState := buildSinglePlayerGame(playerHand, dealerHand, &ai)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected surrender to fail")
		}
	}()

	runPlayerHand(&testState.players[0], &testState.players[0].hands[0], dealerHand[0], &testState.deck, true)
	if containsMove(ai.sentAllowedMoves, Surrender) {
		t.Error("Surrender was passed to AI as an allowed move when it should not have been")
	}
}

func TestPlayer_CannotSurrenderAfterSplit(t *testing.T) {
	dealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand := []deck.Card{{Rank: deck.Five, Suit: deck.Club}, {Rank: deck.Five, Suit: deck.Heart}}
	ai := testAI{returnedMove: Surrender}
	testState := buildSinglePlayerGame(playerHand, dealerHand, &ai)
	testState.players[0].hands[0].split = true

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected surrender to fail")
		}
	}()

	runPlayerHand(&testState.players[0], &testState.players[0].hands[0], dealerHand[0], &testState.deck, true)
	if containsMove(ai.sentAllowedMoves, Surrender) {
		t.Error("Surrender was passed to AI as an allowed move when it should not have been")
	}
}

func TestPlayer_AICannotModifyHands(t *testing.T) {
	dealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
	playerHand := []deck.Card{{Rank: deck.Two, Suit: deck.Club}, {Rank: deck.Two, Suit: deck.Heart}}
	ai := cheatingAI{}
	testState := buildSinglePlayerGame(playerHand, dealerHand, &ai)
	nextCard := testState.deck.cards[0]

	runPlayerHand(&testState.players[0], &testState.players[0].hands[0], dealerHand[0], &testState.deck, true)
	expectedHand := append([]deck.Card{{Rank: deck.Two, Suit: deck.Club}, {Rank: deck.Two, Suit: deck.Heart}}, nextCard)
	currentPlayerHand := testState.players[0].hands[0].cards
	if !reflect.DeepEqual(currentPlayerHand, expectedHand) {
		t.Errorf("AI was able to change player cards. Expected hand %v, was %v", expectedHand, currentPlayerHand)
	}
}

func buildSinglePlayerGame(playerHand, dealerHand []deck.Card, ai AI) state {
	testPlayer := player{
		hands: []hand{
			{
				cards: playerHand,
				bet:   200,
			},
		},
		balance:    1000,
		ai:         ai,
		initialBet: 200,
	}
	testState := state{
		players:    []player{testPlayer},
		dealerHand: dealerHand,
	}
	shuffle(&testState, Options{NumDecks: 2})
	return testState
}

func expectHand(t *testing.T, hand, expectedHand []deck.Card) {
	if !reflect.DeepEqual(hand, expectedHand) {
		t.Errorf("Expected hand to be %v but was %v", expectedHand, hand)
	}
}

func expectBet(t *testing.T, bet, exectedBet money.USD) {
	if bet != exectedBet {
		t.Errorf("Expected bet to be %s, but was %s", exectedBet, bet)
	}
}

func expectFinished(t *testing.T, finished, expectedFinished bool) {
	if finished != expectedFinished {
		t.Errorf("Expected finished to be %t but was %t", expectedFinished, finished)
	}
}

func expectAIPlayHand(t *testing.T, handSent, expectedHandSent []deck.Card) {
	if !reflect.DeepEqual(expectedHandSent, handSent) {
		t.Errorf("Expected to call AI with hand %v, but called with %v", expectedHandSent, handSent)
	}
}

func expectAIDealerCard(t *testing.T, cardSent, expectedCardSent deck.Card) {
	if !reflect.DeepEqual(expectedCardSent, cardSent) {
		t.Errorf("Expected to call AI with dealer card %s, but called with %s", expectedCardSent.String(), cardSent.String())
	}
}
