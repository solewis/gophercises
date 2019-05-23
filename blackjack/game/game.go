package blackjack

import (
	"errors"
	"fmt"
	"gophercises/blackjack/money"
	"gophercises/carddeck"
	"math"
	"strings"
)

//TODO
// surrender only allowed on first returnedMove, not after split

type AI interface {
	Name() string
	Bet(minBet money.USD, maxBet money.USD) money.USD
	Play(hand []deck.Card, dealerShowing deck.Card, allowedMoves []Move) Move
	Results(hands [][]deck.Card, dealer []deck.Card, winnings, balance money.USD)
}

type Move int

const (
	Hit Move = iota
	Stand
	Double
	Split
	Surrender
)

func (m Move) String() string {
	return []string{"Hit", "Stand", "Double", "Split", "Surrender"}[m]
}

func Play(opts Options) map[string]money.USD {
	var s = state{}
	shuffle(&s, opts)
	for name, ai := range opts.AIs {
		s.players = append(s.players, player{name: name, ai: ai, balance: 0})
	}

	for i := 0; i < opts.NumRounds; i++ {
		if s.deck.needsReshuffle(opts.PercentDeckUsage) {
			shuffle(&s, opts)
		}
		bet(&s, opts)
		deal(&s)

		if NaturalBlackjack(s.dealerHand) {
			handleResults(&s, opts)
			continue
		}

		for playerIdx := range s.players {
			for handIdx := 0; handIdx < len(s.players[playerIdx].hands); handIdx++ {
				for {
					//TODO can this be recursive and testable?
					finished := runHand(&s.players[playerIdx], &s.players[playerIdx].hands[handIdx], s.dealerHand[0], s.deck)
					if finished {
						break
					}
				}
			}
		}

		for score, soft := Score(s.dealerHand); score <= 16 || (score == 17 && soft); score, soft = Score(s.dealerHand) {
			s.dealerHand = append(s.dealerHand, s.deck.draw())
		}

		handleResults(&s, opts)
	}
	fmt.Printf("%v", s.players)
	var balances = make(map[string]money.USD)
	for _, p := range s.players {
		balances[p.name] = p.balance
	}
	return balances
}

func Score(cards []deck.Card) (value int, soft bool) {
	score := 0
	lowAces := 0
	for _, card := range cards {
		score += int(math.Min(float64(card.Rank), 10))
		if card.Rank == deck.Ace {
			lowAces += 1
		}
	}

	if score <= 11 && lowAces > 0 {
		return score + 10, true
	} else {
		return score, false
	}
}

func HandString(hand []deck.Card) string {
	strList := make([]string, len(hand))
	for i, card := range hand {
		strList[i] = card.String()
	}
	return strings.Join(strList, ", ")
}

func NaturalBlackjack(cards []deck.Card) bool {
	score, _ := Score(cards)
	return len(cards) == 2 && score == 21
}

type Options struct {
	MinBet                     money.USD
	MaxBet                     money.USD
	NaturalBlackjackMultiplier float64
	NumDecks                   int
	NumRounds                  int
	PercentDeckUsage           float64
	AIs                        map[string]AI
}

type state struct {
	deck       blackjackDeck
	players    []player
	dealerHand []deck.Card
	handIdx    int
}

type blackjackDeck struct {
	cards   []deck.Card
	discard []deck.Card
}

func (d *blackjackDeck) draw() deck.Card {
	card, cards := d.cards[0], d.cards[1:]
	d.cards = cards
	d.discard = append(d.discard, card)
	return card
}

func (d *blackjackDeck) needsReshuffle(percentToUse float64) bool {
	size := len(d.cards) + len(d.discard)
	return (float64(len(d.discard)) / float64(size)) > percentToUse
}

type player struct {
	name       string
	ai         AI
	balance    money.USD
	hands      []hand
	initialBet money.USD
}

type hand struct {
	cards       []deck.Card
	bet         money.USD
	surrendered bool
}

func (h hand) string() string {
	strList := make([]string, len(h.cards))
	for i, card := range h.cards {
		strList[i] = card.String()
	}
	return strings.Join(strList, ", ")
}

func shuffle(s *state, opts Options) {
	s.deck = blackjackDeck{
		deck.New(deck.Count(opts.NumDecks), deck.Shuffle),
		make([]deck.Card, 0, 10),
	}
}

func bet(s *state, opts Options) {
	for pIdx := range s.players {
		p := &s.players[pIdx]
		bet := p.ai.Bet(opts.MinBet, opts.MaxBet)
		if bet.Float64() < opts.MinBet.Float64() {
			panic(errors.New("bet must be greater than or equal to " + opts.MinBet.String()))
		}
		if bet.Float64() > opts.MaxBet.Float64() {
			panic(errors.New("bet must be less than or equal to " + opts.MaxBet.String()))
		}
		p.initialBet = bet
	}
}

func deal(s *state) {
	playerHands := make([][]deck.Card, len(s.players))
	s.dealerHand = make([]deck.Card, 0, 21)

	for pIdx := range playerHands {
		playerHands[pIdx] = make([]deck.Card, 0, 21)
	}

	for card := 0; card < 2; card++ {
		for pIdx:=range s.players {
			playerHands[pIdx] = append(playerHands[pIdx], s.deck.draw())
		}
		s.dealerHand = append(s.dealerHand, s.deck.draw())
	}

	for pIdx := range s.players {
		p:= &s.players[pIdx]
		p.hands = []hand{
			{
				cards: playerHands[pIdx],
				bet:   p.initialBet,
			},
		}
	}

	s.handIdx = 0
}

func runHand(p *player, h *hand, ds deck.Card, d blackjackDeck) (finished bool) {
	//currentPlayer := s.currentPlayer()
	//playerHand := &s.player.hands[s.handIdx].cards
	cards := &(*h).cards
	playerHandCopy := make([]deck.Card, len(*cards))
	copy(playerHandCopy, *cards)

	allowedMoves := allowedMoves(*cards)
	move := p.ai.Play(playerHandCopy, ds, allowedMoves)

	switch move {
	case Hit:
		*cards = append(append(*cards, d.draw()))
		score, _ := Score(*cards)
		return score > 21
	case Stand:
		return true
	case Double:
		if !containsMove(allowedMoves, Double) {
			panic(errors.New("can only double on a hand with two cards"))
		}
		*cards = append(append(*cards, d.draw()))
		h.bet *= 2
		return true
	case Split:
		if !containsMove(allowedMoves, Split) {
			panic(errors.New("can only split on a hand with two identical cards"))
		}
		splitCard := (*cards)[1]
		*cards = []deck.Card{(*cards)[0], d.draw()}
		p.hands = append(p.hands, hand{
			cards: []deck.Card{splitCard, d.draw()},
			bet:   p.initialBet,
		})
		return false
	case Surrender:
		if !containsMove(allowedMoves, Surrender) {
			panic(errors.New("can only surrender when you have two cards"))
		}
		h.surrendered = true
		return true
	default:
		panic(errors.New("unhandled move"))
	}
}

//func runPlayerTurn(s *state) {
//	currentPlayer := s.currentPlayer()
//	playerHand := &s.player.hands[s.handIdx].cards
//	playerHandCopy := make([]deck.Card, len(*playerHand))
//	copy(playerHandCopy, *playerHand)
//
//	allowedMoves := allowedMoves(*playerHand)
//	move := s.player.ai.Play(playerHandCopy, s.dealerHand[0], allowedMoves)
//
//	switch move {
//	case Hit:
//		*playerHand = append(append(*playerHand, s.deck.draw()))
//		if score, _ := Score(*playerHand); score > 21 {
//			s.handIdx++
//		}
//	case Stand:
//		s.handIdx++
//	case Double:
//		if !containsMove(allowedMoves, Double) {
//			panic(errors.New("can only double on a hand with two cards"))
//		}
//		*playerHand = append(append(*playerHand, s.deck.draw()))
//		s.player.hands[s.handIdx].bet *= 2
//		s.handIdx++
//	case Split:
//		if !containsMove(allowedMoves, Split) {
//			panic(errors.New("can only split on a hand with two identical cards"))
//		}
//		splitCard := (*playerHand)[1]
//		*playerHand = []deck.Card{(*playerHand)[0], s.deck.draw()}
//		s.player.hands = append(s.player.hands, hand{
//			cards: []deck.Card{splitCard, s.deck.draw()},
//			bet:   s.player.initialBet,
//		})
//	case Surrender:
//		if !containsMove(allowedMoves, Surrender) {
//			panic(errors.New("can only surrender when you have two cards"))
//		}
//		s.player.hands[s.handIdx].surrendered = true
//		s.handIdx++
//	}
//}

func containsMove(moves []Move, move Move) bool {
	for _, m := range moves {
		if m == move {
			return true
		}
	}
	return false
}

func allowedMoves(playerHand []deck.Card) []Move {
	allowedMoves := []Move{Hit, Stand}
	if len(playerHand) == 2 {
		allowedMoves = append(allowedMoves, []Move{Double, Surrender}...)
		if (playerHand)[0].Rank == (playerHand)[1].Rank {
			allowedMoves = append(allowedMoves, Split)
		}
	}
	return allowedMoves
}

func handleResults(s *state, opts Options) {
	dealerScore, _ := Score(s.dealerHand)

	for pIdx := range s.players {
		p := &s.players[pIdx]
		allHands := make([][]deck.Card, len(p.hands))
		var winnings money.USD = 0
		for i, playerHand := range p.hands {
			allHands[i] = playerHand.cards
			handScore, _ := Score(playerHand.cards)
			switch {
			case NaturalBlackjack(s.dealerHand) && NaturalBlackjack(playerHand.cards):
				//draw
			case NaturalBlackjack(s.dealerHand):
				winnings -= playerHand.bet
			case playerHand.surrendered == true:
				winnings -= playerHand.bet.Multiply(.5)
			case NaturalBlackjack(playerHand.cards):
				winnings += playerHand.bet.Multiply(opts.NaturalBlackjackMultiplier)
			case handScore > 21:
				winnings -= playerHand.bet
			case dealerScore > 21:
				winnings += playerHand.bet
			case handScore > dealerScore:
				winnings += playerHand.bet
			case dealerScore > handScore:
				winnings -= playerHand.bet
			case dealerScore == handScore:
				//draw
			}
		}
		p.balance += winnings
		p.ai.Results(allHands, s.dealerHand, winnings, p.balance)
	}
}
