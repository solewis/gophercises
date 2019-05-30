package ai

import (
	blackjack "gophercises/blackjack/game"
	"gophercises/blackjack/money"
	deck "gophercises/carddeck"
	"math"
)

func BasicHiLo(opts blackjack.Options) blackjack.AI {
	return &basicHiLo{totalCards: opts.NumDecks * 52}
}

type basicHiLo struct {
	count, cardsPlayed, totalCards int
}

func (_ basicHiLo) Name() string {
	return "Basic Hi-Lo Card Counter"
}

func (_ basicHiLo) Play(hand []deck.Card, dealerShowing deck.Card, allowedMoves []blackjack.Move) blackjack.Move {
	pScore, soft := blackjack.Score(hand)
	pPair := len(hand) == 2 && hand[0].Rank == hand[1].Rank
	dScore, _ := blackjack.Score([]deck.Card{dealerShowing})
	return lookup[playerHand{pScore, soft, pPair}][dScore](allowedMoves)
}

func (ai *basicHiLo) Bet(minBet money.USD, maxBet money.USD, shuffled bool) money.USD {
	if shuffled {
		ai.count = 0
		ai.cardsPlayed = 0
	}
	//TODO integer division is floor math, can it be modified to round (i.e. considered 4 decks left until halfway through the 4th
	remainingDecks := (ai.totalCards - ai.cardsPlayed) / 52
	trueCount := ai.count / int(math.Max(float64(remainingDecks), 1))
	var bet float64 = 5
	//TODO represent bet spread in a data structure
	switch {
	case trueCount >= 5:
		bet *= 12
	case trueCount == 4:
		bet *= 8
	case trueCount == 3:
		bet *= 4
	case trueCount == 2:
		bet *= 2
	}
	usd := money.ToUSD(bet)
	return usd
}

func (ai basicHiLo) HandResults(hand, dealer []deck.Card, winnings, balance money.USD) {
	//fmt.Printf("Hand: %s, dealer: %s, winnings: %s, balance: %s. count: %d. Cards played: %d\n", blackjack.HandString(hand), blackjack.HandString(dealer), winnings.String(), balance.String(), ai.count, ai.cardsPlayed)
}

func (ai *basicHiLo) RoundRecap(allHands [][]deck.Card) {
	for _, h := range allHands {
		for _, c := range h {
			switch {
			case c.Rank <= 6:
				ai.count++
			case c.Rank >= 10:
				ai.count--
			}
		}
		ai.cardsPlayed += len(h)
	}
}
