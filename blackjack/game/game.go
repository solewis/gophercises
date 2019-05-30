package blackjack

import (
	"errors"
	"fmt"
	"gophercises/blackjack/money"
	"gophercises/carddeck"
	"math"
	"sort"
	"strings"
	"time"
)

var (
	shuffles, bets, deals, playerHands, resultHandlings []int64
)

//TODO
// if all players bust, does dealer still draw? Probably not

type AI interface {
	Name() string
	Bet(minBet money.USD, maxBet money.USD, shuffled bool) money.USD
	Play(hand []deck.Card, dealerShowing deck.Card, allowedMoves []Move) Move
	HandResults(hand, dealer []deck.Card, winnings, balance money.USD)
	RoundRecap(allHands [][]deck.Card)
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

func Play(ais map[string]AI, opts Options) map[string]money.USD {
	var s = state{}
	shuffle(&s, opts)
	for name, ai := range ais {
		s.players = append(s.players, player{name: name, ai: ai, balance: 0})
	}

	for i := 0; i < opts.NumRounds; i++ {
		shuffled := false
		if s.deck.needsReshuffle(opts.PercentDeckUsage) {
			shuffle(&s, opts)
			shuffled = true
		}
		bet(&s, opts, shuffled)
		deal(&s)

		if NaturalBlackjack(s.dealerHand) {
			handleResults(&s, opts)
			continue
		}

		for playerIdx := range s.players {
			for handIdx := 0; handIdx < len(s.players[playerIdx].hands); handIdx++ {
				firstPlay := true
				for {
					finished := runPlayerHand(&s.players[playerIdx], &s.players[playerIdx].hands[handIdx], s.dealerHand[0], &s.deck, firstPlay)
					firstPlay = false
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

	var balances = make(map[string]money.USD)
	for _, p := range s.players {
		balances[p.name] = p.balance
	}
	//shuffles, bets, deals, playerHands, resultHandlings
	shuffles = removeOutliers(shuffles)
	bets = removeOutliers(bets)
	deals = removeOutliers(deals)
	playerHands = removeOutliers(playerHands)
	resultHandlings = removeOutliers(resultHandlings)
	fmt.Printf("Shuffles: %.2f. Count: %d. Total: %d. Stddev: %.2f\n", avg(shuffles), len(shuffles), sum(shuffles), stddev(shuffles))
	fmt.Printf("Bets: %.2f. Count: %d. Total: %d. Stddev: %.2f\n", avg(bets), len(bets), sum(bets), stddev(bets))
	fmt.Printf("Deals: %.2f. Count: %d. Total: %d. Stddev: %.2f\n", avg(deals), len(deals), sum(deals), stddev(deals))
	fmt.Printf("Player hands: %.2f. Count: %d. Total: %d. Stddev: %.2f\n", avg(playerHands), len(playerHands), sum(playerHands), stddev(playerHands))
	fmt.Printf("Results handlings: %.2f. Count: %d. Total: %d. Stddev: %.2f\n", avg(resultHandlings), len(resultHandlings), sum(resultHandlings), stddev(resultHandlings))
	return balances
}

func removeOutliers(store []int64) []int64 {
	sort.Slice(store, func(i, j int) bool { return store[i] < store[j] })
	numToRemove := int( .1 * float64(len(store)))
	return store[0+numToRemove: len(store) - numToRemove]
}

func stddev(store []int64) float64 {
	mean := avg(store)
	var intermediates = make([]float64, len(store))
	for i, s := range store {
		intermediates[i] = math.Pow(float64(s) - mean, 2)
	}
	return math.Sqrt(avgF(intermediates))
}
func avgF(store []float64) float64 {
	var total float64
	for _, s := range store {
		total+=s
	}
	return total / float64(len(store))
}

func avg(store []int64) float64 {
	var total int64
	for _, s := range store {
		total+=s
	}
	return float64(total) / float64(len(store))
}

func sum(store []int64) int64 {
	var total int64
	for _, s := range store {
		total+=s
	}
	return total
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
	Double                     bool
	Surrender                  bool
	DoubleAfterSplit           bool
}

type state struct {
	deck       blackjackDeck
	players    []player
	dealerHand []deck.Card
}

type blackjackDeck struct {
	cards      []deck.Card
	discard    []deck.Card
	discardIdx int
}

func (d *blackjackDeck) draw() deck.Card {
	card, cards := d.cards[0], d.cards[1:]
	d.cards = cards
	d.discard[d.discardIdx] = card
	d.discardIdx++
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
	split       bool
}

func (h hand) string() string {
	strList := make([]string, len(h.cards))
	for i, card := range h.cards {
		strList[i] = card.String()
	}
	return strings.Join(strList, ", ")
}

func shuffle(s *state, opts Options) {
	defer track(time.Now(), &shuffles)
	cards := deck.New(deck.Count(opts.NumDecks), deck.Shuffle)
	s.deck = blackjackDeck{
		cards:   cards,
		discard: make([]deck.Card, len(cards)),
	}
}

func bet(s *state, opts Options, shuffled bool) {
	defer track(time.Now(), &bets)
	for pIdx := range s.players {
		p := &s.players[pIdx]
		bet := p.ai.Bet(opts.MinBet, opts.MaxBet, shuffled)
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
	defer track(time.Now(), &deals)

	cardsToDeal := 2
	playerHands := make([][]deck.Card, len(s.players))
	s.dealerHand = make([]deck.Card, cardsToDeal, 21)

	for pIdx := range playerHands {
		playerHands[pIdx] = make([]deck.Card, cardsToDeal, 21)
	}

	for card := 0; card < cardsToDeal; card++ {
		for pIdx := range s.players {
			playerHands[pIdx][card] = s.deck.draw()
		}
		s.dealerHand[card] = s.deck.draw()
	}

	for pIdx := range s.players {
		p := &s.players[pIdx]
		p.hands = []hand{
			{
				cards: playerHands[pIdx],
				bet:   p.initialBet,
			},
		}
	}
}

func track(start time.Time, store *[]int64) {
	*store = append(*store, time.Since(start).Nanoseconds())
}

func runPlayerHand(
	p *player,
	h *hand,
	ds deck.Card,
	d *blackjackDeck,
	firstPlay bool) (finished bool) {

	defer track(time.Now(), &playerHands)

	cards := &(*h).cards
	playerHandCopy := make([]deck.Card, len(*cards))
	copy(playerHandCopy, *cards)

	allowedMoves := allowedMoves(*h, firstPlay)
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
			split: true,
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

func containsMove(moves []Move, move Move) bool {
	for _, m := range moves {
		if m == move {
			return true
		}
	}
	return false
}

func allowedMoves(playerHand hand, firstPlay bool) []Move {
	allowedMoves := []Move{Hit, Stand}
	if len(playerHand.cards) == 2 {
		allowedMoves = append(allowedMoves, Double)
		if firstPlay && !playerHand.split {
			allowedMoves = append(allowedMoves, Surrender)
		}
		if playerHand.cards[0].Rank == playerHand.cards[1].Rank {
			allowedMoves = append(allowedMoves, Split)
		}
	}
	return allowedMoves
}

func handleResults(s *state, opts Options) {
	defer track(time.Now(), &resultHandlings)
	dealerScore, _ := Score(s.dealerHand)

	allHands := make([][]deck.Card, 0, len(s.players))
	for pIdx := range s.players {
		p := &s.players[pIdx]
		for _, playerHand := range p.hands {
			var winnings money.USD = 0
			allHands = append(allHands, playerHand.cards)
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
			p.balance += winnings
			p.ai.HandResults(playerHand.cards, s.dealerHand, winnings, p.balance)
		}
	}
	allHands = append(allHands, s.dealerHand)
	for pIdx := range s.players {
		s.players[pIdx].ai.RoundRecap(allHands)
	}
}
