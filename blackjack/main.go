package main

//TODO
//allow user to pick rules
//calculate house edge based on rules
//rebuild deck after 75% through
import (
	"fmt"
	"gophercises/blackjack/util"
	"gophercises/carddeck"
	"math"
	"strings"
)

type AI interface {
	Bet() util.USD
	Play(hand []deck.Card, dealerShowing deck.Card) Move
}

type Move func(*GameState)

const (
	minBet                     util.USD = 200
	maxBet                     util.USD = 50000
	naturalBlackjackMultiplier float64  = 1.5
	numDecks                   int      = 4
	playerStartingUSD          util.USD = 10000
)

type GameState struct {
	deck BlackjackDeck
	player Player
	playerHand Hand
	dealerHand Hand
}

type BlackjackDeck struct {
	Cards []deck.Card
}

func (d *BlackjackDeck) Deal() deck.Card {
	card, cards := d.Cards[0], d.Cards[1:]
	d.Cards = cards
	return card
}

type Player struct {
	USD        util.USD
	TotalGames int
}

type Hand []deck.Card

func (h Hand) String() string {
	strList := make([]string, len(h))
	for i, card := range h {
		strList[i] = card.String()
	}
	return strings.Join(strList, ", ")
}

func (h Hand) FirstCardString() string {
	return h[0].String()
}

func (h Hand) IsNaturalBlackjack() bool {
	return h.Score().Value == 21 && len(h) == 2
}

type Score struct {
	Value int
	Soft  bool
}

func (s Score) IsSoft17() bool {
	return s.Value == 17 && s.Soft
}

func (h Hand) Score() Score {
	score := 0
	lowAces := 0
	for _, card := range h {
		score += int(math.Min(float64(card.Rank), 10))
		if card.Rank == deck.Ace {
			lowAces += 1
		}
	}

	if score <= 11 && lowAces > 0 {
		return Score{score + 10, true}
	} else {
		return Score{score, false}
	}
}

func main() {
	fmt.Println("Welcome to blackjack")
	fmt.Println("Table rules:")
	fmt.Println("1. Dealer hits on soft 17")
	fmt.Printf("2. Minimum bet %s\n", minBet.String())
	fmt.Printf("3. Maximum bet %s\n", maxBet.String())
	fmt.Println("4. Dealer does not automatically win on blackjack")
	fmt.Printf("5. Payout on natural blackjack %.2f\n", naturalBlackjackMultiplier)
	fmt.Printf("6. %d decks\n", numDecks)
	//fmt.Println("7. No double down restrictions (can double down on any hand")
	//fmt.Println("8. No split restrictions")
	//fmt.Println("9. Early and late surrender allowed")
	//fmt.Println("10. Continuous shuffle?")

	var gs = GameState{}
	shuffle(&gs)
	addPlayer(&gs)

	var playAgain string
	for playAgain != "n" {
		playRound(&gs)

		if gs.player.USD < minBet {
			fmt.Printf("You have played %d games. You have %s. You do not have enough money to continue\n", gs.player.TotalGames, gs.player.USD.String())
			playAgain = "n"
		} else {
			fmt.Printf("You have played %d game(s). You have %s. Play again? (yn)\n", gs.player.TotalGames, gs.player.USD.String())
			fmt.Scanln(&playAgain)
		}
	}
}

func shuffle(gs *GameState) {
	gs.deck = BlackjackDeck{deck.New(deck.DeckCount(numDecks), deck.Shuffle)}
}

func addPlayer(gs *GameState) {
	gs.player = Player{USD: playerStartingUSD}
}

func playRound(gs *GameState) {
	bet := retrieveBet(gs.player)
	deal(gs)

	if gs.dealerHand.IsNaturalBlackjack() || gs.playerHand.IsNaturalBlackjack() {
		handleResults(gs.playerHand, gs.dealerHand, &gs.player, bet)
		return
	}

	fmt.Printf("Dealer showing: %s\n", gs.dealerHand.FirstCardString())
	fmt.Printf("Your hand: %s\n", gs.playerHand.String())

	bet = runPlayer(gs, bet)
	gs.dealerHand = runDealer(gs.dealerHand, &gs.deck)

	handleResults(gs.playerHand, gs.dealerHand, &gs.player, bet)
}

func deal(gs *GameState) {
	gs.playerHand = make(Hand, 0, 21)
	gs.dealerHand = make(Hand, 0, 21)
	gs.playerHand = append(gs.playerHand, gs.deck.Deal())
	gs.dealerHand = append(gs.dealerHand, gs.deck.Deal())
	gs.playerHand = append(gs.playerHand, gs.deck.Deal())
	gs.dealerHand = append(gs.dealerHand, gs.deck.Deal())
}

func handleResults(playerHand Hand, dealerHand Hand, player *Player, bet util.USD) {
	dealerScore, playerScore := dealerHand.Score().Value, playerHand.Score().Value
	fmt.Println("Final hands:")
	fmt.Printf("Dealer: %s. Score: %d\n", dealerHand.String(), dealerScore)
	fmt.Printf("You: %s. Score: %d\n", playerHand.String(), playerScore)

	switch {
	case dealerHand.IsNaturalBlackjack() && playerHand.IsNaturalBlackjack():
		fmt.Println("Both have natural blackjack. Draw")
	case dealerHand.IsNaturalBlackjack():
		fmt.Printf("Dealer has natural blackjack, you lose %s\n", bet.String())
		player.USD -= bet
	case playerHand.IsNaturalBlackjack():
		payout := bet.Multiply(naturalBlackjackMultiplier)
		fmt.Printf("You have natural blackjack, you win %s\n", payout.String())
		player.USD += payout
	case playerScore > 21:
		fmt.Printf("You busted, you lose %s\n", bet.String())
		player.USD -= bet
	case dealerScore > 21:
		fmt.Printf("Dealer busted, you win %s\n", bet.String())
		player.USD += bet
	case playerScore > dealerScore:
		fmt.Printf("You win %s\n", bet.String())
		player.USD += bet
	case dealerScore > playerScore:
		fmt.Printf("You lose %s\n", bet.String())
		player.USD -= bet
	default:
		fmt.Println("Draw")
	}
	player.TotalGames++
}

func retrieveBet(player Player) util.USD {
	fmt.Printf("You have %s.\n", player.USD.String())
	var betInput float64
	maxBet := util.ToUSD(math.Min(player.USD.Float64(), maxBet.Float64()))
	fmt.Printf("How much would you like to bet? %s (min) %s (max)\n", minBet.String(), maxBet.String())
	for {
		_, e := fmt.Scanf("%f\n", &betInput)
		switch {
		case e != nil:
			fmt.Println("Invalid bet... Must be a float with a max of 2 decimal places")
		case betInput < minBet.Float64():
			fmt.Println("Bet must be greater than", minBet.String())
		case betInput > maxBet.Float64():
			fmt.Println("Bet must be less than", maxBet.String())
		default:
			bet := util.ToUSD(betInput)
			fmt.Printf("Your bet is %s\n", bet.String())
			return bet
		}
	}
}

func runPlayer(gs *GameState, bet util.USD) (updatedBet util.USD) {
	var playerChoice int
	for playerChoice != 2 {
		fmt.Println("Would you like to:")
		fmt.Println("1. Hit")
		fmt.Println("2. Stand")
		if len(gs.playerHand) == 2 && gs.player.USD >= bet.Multiply(2) {
			fmt.Println("3. Double down")
		}
		fmt.Scanln(&playerChoice)
		switch playerChoice {
		case 1:
			gs.playerHand = append(gs.playerHand, gs.deck.Deal())
			fmt.Printf("Your hand: %s\n", gs.playerHand.String())
		case 3:
			gs.playerHand = append(gs.playerHand, gs.deck.Deal())
			bet = bet.Multiply(2)
			fmt.Printf("Bet is now %s\n", bet.String())
			fmt.Printf("Your hand: %s\n", gs.playerHand.String())
			return bet
		}
		if gs.playerHand.Score().Value > 21 {
			return bet
		}
	}
	return bet
}

func runDealer(dealerHand Hand, blackjackDeck *BlackjackDeck) Hand {
	for dealerHand.Score().Value <= 16 || dealerHand.Score().IsSoft17() {
		dealerHand = append(dealerHand, blackjackDeck.Deal())
	}
	return dealerHand
}
