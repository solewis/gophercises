package deck

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	testDeck := New()
	// 13 ranks * 4 suits
	if len(testDeck) != 13*4 {
		t.Error("Wrong number of cards in a new deck.")
	}
}

func TestDeckCount(t *testing.T) {
	testDeck := New(DeckCount(3))

	if len(testDeck) != 13*4*3 {
		t.Errorf("Expected %d cards, received %d cards.", 13*4*3, len(testDeck))
	}
}

func TestShuffle(t *testing.T) {
	testDeckShuffled := New(Shuffle)
	testDeckDefault := New()

	if reflect.DeepEqual(testDeckShuffled, testDeckDefault) {
		t.Error("Expected shuffled deck, but deck was in order")
	}
}

func TestShuffleWithMultipleDecks(t *testing.T) {
	testDeck := New(DeckCount(2), Shuffle)

	if testDeck[0] == testDeck[52] && testDeck[1] == testDeck[53] {
		t.Error("Expected shuffling to happen on expanded deck, instead it happened on original deck then was duplicated to expand")
	}
}