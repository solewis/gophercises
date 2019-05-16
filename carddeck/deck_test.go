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

func TestExpand(t *testing.T) {
	testDeck := New(Expand(3))

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