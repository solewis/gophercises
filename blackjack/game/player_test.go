package blackjack

//func TestPlayer_Hit(t *testing.T) {
//	dealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
//	playerHand := []deck.Card{{Rank: deck.Two, Suit: deck.Club}, {Rank: deck.Two, Suit: deck.Diamond}}
//	ai := testAI{returnedMove: Hit}
//	testState := setupGame(playerHand, dealerHand, &ai)
//	nextCard := testState.deck.cards[0]
//
//	runPlayerTurn(&testState)
//	expectedHand := append(playerHand, nextCard)
//	verifySingleHand(t, testState, 0, testState.player.initialBet, ai, playerHand, dealerHand, expectedHand)
//}
//
//func TestPlayer_HitAndBust(t *testing.T) {
//	dealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
//	playerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
//	ai := testAI{returnedMove: Hit}
//	testState := setupGame(playerHand, dealerHand, &ai)
//	nextCard := testState.deck.cards[0]
//
//	runPlayerTurn(&testState)
//	expectedHand := append(playerHand, nextCard)
//	verifySingleHand(t, testState, 1, testState.player.initialBet, ai, playerHand, dealerHand, expectedHand)
//}
//
//func TestPlayer_Split(t *testing.T) {
//	dealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
//	playerHand := []deck.Card{{Rank: deck.Eight, Suit: deck.Club}, {Rank: deck.Eight, Suit: deck.Diamond}}
//	ai := testAI{returnedMove: Split}
//	testState := setupGame(playerHand, dealerHand, &ai)
//	nextCards := testState.deck.cards[0:2]
//
//	runPlayerTurn(&testState)
//	expectedHand1 := []deck.Card{playerHand[0], nextCards[0]}
//	expectedHand2 := []deck.Card{playerHand[1], nextCards[1]}
//	if !reflect.DeepEqual(testState.player.hands[0].cards, expectedHand1) {
//		t.Errorf("Expected first hand to be %v but was %v", expectedHand1, testState.player.hands[0].cards)
//	}
//	if !reflect.DeepEqual(testState.player.hands[1].cards, expectedHand2) {
//		t.Errorf("Expected second hand to be %v but was %v", expectedHand2, testState.player.hands[1].cards)
//	}
//	if testState.player.hands[1].bet != testState.player.initialBet {
//		t.Errorf("Expected bet for second hand to be %s, but was %s", testState.player.initialBet, testState.player.hands[1].bet)
//	}
//	if testState.handIdx != 0 {
//		t.Error("Expected hand index to stay at 0")
//	}
//}
//
//func TestPlayer_CannotSplitIfHandNotTwoCards(t *testing.T) {
//	dealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
//	playerHand := []deck.Card{{Rank: deck.Five, Suit: deck.Club}, {Rank:deck.Five, Suit: deck.Heart}, {Rank: deck.Six, Suit: deck.Diamond}}
//	ai := testAI{returnedMove: Split}
//	testState := setupGame(playerHand, dealerHand, &ai)
//
//	defer func() {
//		if r := recover(); r == nil {
//			t.Error("Expected split to fail")
//		}
//	}()
//
//	runPlayerTurn(&testState)
//	if containsMove(ai.sentAllowedMoves, Split) {
//		t.Error("Split was passed to AI as an allowed move when it should not have been")
//	}
//}
//
//func TestPlayer_CannotSplitIfCardsNotEqualRank(t *testing.T) {
//	dealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
//	playerHand := []deck.Card{{Rank: deck.Five, Suit: deck.Club}, {Rank: deck.Six, Suit: deck.Diamond}}
//	ai := testAI{returnedMove: Split}
//	testState := setupGame(playerHand, dealerHand, &ai)
//
//	defer func() {
//		if r := recover(); r == nil {
//			t.Error("Expected split to fail")
//		}
//	}()
//
//	runPlayerTurn(&testState)
//	if containsMove(ai.sentAllowedMoves, Split) {
//		t.Error("Split was passed to AI as an allowed move when it should not have been")
//	}
//}
//
//func TestPlayer_Double(t *testing.T) {
//	dealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
//	playerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Club}, {Rank: deck.Two, Suit: deck.Diamond}}
//	ai := testAI{returnedMove: Double}
//	testState := setupGame(playerHand, dealerHand, &ai)
//	nextCard := testState.deck.cards[0]
//
//	runPlayerTurn(&testState)
//	expectedHand := append(playerHand, nextCard)
//	expectedBet := testState.player.initialBet.Multiply(2)
//	verifySingleHand(t, testState, 1, expectedBet, ai, playerHand, dealerHand, expectedHand)
//}
//
//func TestPlayer_CannotDoubleIfHandNotTwoCards(t *testing.T) {
//	dealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
//	playerHand := []deck.Card{{Rank: deck.Five, Suit: deck.Club}, {Rank:deck.Five, Suit: deck.Heart}, {Rank: deck.Six, Suit: deck.Diamond}}
//	ai := testAI{returnedMove: Double}
//	testState := setupGame(playerHand, dealerHand, &ai)
//
//	defer func() {
//		if r := recover(); r == nil {
//			t.Error("Expected double to fail")
//		}
//	}()
//
//	runPlayerTurn(&testState)
//	if containsMove(ai.sentAllowedMoves, Double) {
//		t.Error("Double was passed to AI as an allowed move when it should not have been")
//	}
//}
//
//func TestPlayer_Stand(t *testing.T) {
//	dealerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
//	playerHand := []deck.Card{{Rank: deck.Nine, Suit: deck.Club}, {Rank: deck.Ten, Suit: deck.Diamond}}
//	ai := testAI{returnedMove: Stand}
//	testState := setupGame(playerHand, dealerHand, &ai)
//
//	runPlayerTurn(&testState)
//	verifySingleHand(t, testState, 1, testState.player.initialBet, ai, playerHand, dealerHand, playerHand)
//}
//
//func TestPlayer_Surrender(t *testing.T) {
//	dealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
//	playerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Club}, {Rank: deck.Six, Suit: deck.Diamond}}
//	ai := testAI{returnedMove: Surrender}
//	testState := setupGame(playerHand, dealerHand, &ai)
//
//	runPlayerTurn(&testState)
//	verifySingleHand(t, testState, 1, testState.player.initialBet, ai, playerHand, dealerHand, playerHand)
//	if !testState.player.hands[0].surrendered {
//		t.Error("Expected hand to be surrendered")
//	}
//}
//
//func TestPlayer_CannotSurrenderIfNotTwoCards(t *testing.T) {
//	dealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
//	playerHand := []deck.Card{{Rank: deck.Five, Suit: deck.Club}, {Rank:deck.Five, Suit: deck.Heart}, {Rank: deck.Six, Suit: deck.Diamond}}
//	ai := testAI{returnedMove: Surrender}
//	testState := setupGame(playerHand, dealerHand, &ai)
//
//	defer func() {
//		if r := recover(); r == nil {
//			t.Error("Expected surrender to fail")
//		}
//	}()
//
//	runPlayerTurn(&testState)
//	if containsMove(ai.sentAllowedMoves, Surrender) {
//		t.Error("Surrender was passed to AI as an allowed move when it should not have been")
//	}
//}
//
//func TestPlayer_AICannotModifyHands(t *testing.T) {
//	dealerHand := []deck.Card{{Rank: deck.Ten, Suit: deck.Spade}, {Rank: deck.Ten, Suit: deck.Club}}
//	playerHand := []deck.Card{{Rank: deck.Two, Suit: deck.Club}, {Rank:deck.Two, Suit: deck.Heart}}
//	ai := cheatingAI{}
//	testState := setupGame(playerHand, dealerHand, &ai)
//	nextCard := testState.deck.cards[0]
//
//	runPlayerTurn(&testState)
//	expectedHand := append([]deck.Card{{Rank: deck.Two, Suit: deck.Club}, {Rank: deck.Two, Suit: deck.Heart}}, nextCard)
//	currentPlayerHand := testState.player.hands[0].cards
//	if !reflect.DeepEqual(currentPlayerHand, expectedHand) {
//		t.Errorf("AI was able to change player cards. Expected hand %v, was %v", expectedHand, currentPlayerHand)
//	}
//}
//
//func setupGame(playerHand, dealerHand []deck.Card, ai AI) state {
//	testPlayer := player{
//		hands: []hand{
//			{
//				cards: playerHand,
//				bet:   200,
//			},
//		},
//		balance:    1000,
//		ai:         ai,
//		initialBet: 200,
//	}
//	testState := state{
//		handIdx:    0,
//		player:     testPlayer,
//		dealerHand: dealerHand,
//	}
//	shuffle(&testState, Options{NumDecks: 2})
//	return testState
//}
//
//func verifySingleHand(
//	t *testing.T, testState state, expectedEndingHandIdx int,
//	expectedEndingBet money.USD, ai testAI,
//	startingPlayerHand, startingDealerHand, expectedEndingPlayerHand []deck.Card) {
//
//	if !reflect.DeepEqual(testState.player.hands[0].cards, expectedEndingPlayerHand) {
//		t.Errorf("Expected hand to be %v but was %v", expectedEndingPlayerHand, testState.player.hands[0].cards)
//	}
//	if testState.player.hands[0].bet != expectedEndingBet {
//		t.Errorf("Expected bet to be %s, but was %s", expectedEndingBet, testState.player.hands[0].bet)
//	}
//	if testState.handIdx != expectedEndingHandIdx {
//		t.Error("Expected hand index to be", expectedEndingHandIdx)
//	}
//	if !reflect.DeepEqual(startingPlayerHand, ai.sentPlayHand) {
//		t.Errorf("Expected to call AI with hand %v, but called with %v", startingPlayerHand, ai.sentPlayHand)
//	}
//	if !reflect.DeepEqual(startingDealerHand[0], ai.sentDealerShowing) {
//		t.Errorf("Expected to call AI with dealer card %s, but called with %s", startingDealerHand[0].String(), ai.sentDealerShowing.String())
//	}
//}
