package main

//TODO
//allow user to pick rules
//calculate house edge based on rules
//rebuild deck after 75% through
import (
	"fmt"
	"gophercises/carddeck"
	"math"
	"strings"
)

const (
	minBet                     USD     = 200
	maxBet                     USD     = 50000
	naturalBlackjackMultiplier float64 = 1.5
	numDecks                   int     = 4
	playerStartingUSD          USD     = 10000
)

type BlackjackDeck struct {
	Cards []deck.Card
}

func (d *BlackjackDeck) Deal() deck.Card {
	card, cards := d.Cards[0], d.Cards[1:]
	d.Cards = cards
	return card
}

type USD int64

func (m USD) String() string {
	x := float64(m)
	x /= 100
	return fmt.Sprintf("$%.2f", x)
}

func (m USD) Float64() float64 {
	x := float64(m)
	x = x / 100
	return x
}

func (m USD) Multiply(f float64) USD {
	x := (float64(m) * f) + 0.5
	return USD(x)
}

func ToUSD(f float64) USD {
	return USD((f * 100) + 0.5)
}

type Player struct {
	USD        USD
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

	player := Player{USD: playerStartingUSD}
	blackjackDeck := BlackjackDeck{deck.New(deck.DeckCount(numDecks), deck.Shuffle)}

	var playAgain string
	for playAgain != "n" {
		playRound(&blackjackDeck, &player)

		if player.USD < minBet {
			fmt.Printf("You have played %d games. You have %s. You do not have enough money to continue\n", player.TotalGames, player.USD.String())
			playAgain = "n"
		} else {
			fmt.Printf("You have played %d game(s). You have %s. Play again? (yn)\n", player.TotalGames, player.USD.String())
			fmt.Scanln(&playAgain)
		}
	}
}

func playRound(blackjackDeck *BlackjackDeck, player *Player) {
	bet := retrieveBet(player)

	var dealerHand, playerHand Hand
	dealerHand = []deck.Card{blackjackDeck.Deal(), blackjackDeck.Deal()}
	playerHand = []deck.Card{blackjackDeck.Deal(), blackjackDeck.Deal()}

	if dealerHand.IsNaturalBlackjack() || playerHand.IsNaturalBlackjack() {
		handleResults(playerHand, dealerHand, player, bet)
		return
	}

	fmt.Printf("Dealer showing: %s\n", dealerHand.FirstCardString())
	fmt.Printf("Your hand: %s\n", playerHand.String())

	playerHand, bet = runPlayer(playerHand, blackjackDeck, bet, player)
	dealerHand = runDealer(dealerHand, blackjackDeck)

	handleResults(playerHand, dealerHand, player, bet)
}

func handleResults(playerHand Hand, dealerHand Hand, player *Player, bet USD) {
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

func retrieveBet(player *Player) USD {
	fmt.Printf("You have %s.\n", player.USD.String())
	var betInput float64
	maxBet := ToUSD(math.Min(player.USD.Float64(), maxBet.Float64()))
	fmt.Printf("How much would you like to bet? %s (min) %s (max)\n", minBet.String(), maxBet.String())
	for betInput < minBet.Float64() || betInput > maxBet.Float64() {
		_, e := fmt.Scanln(&betInput)
		if e != nil {
			fmt.Println("Invalid bet... Try again")
		}
	}
	bet := ToUSD(betInput)
	fmt.Printf("Your bet is %s\n", bet.String())
	return bet
}

func runPlayer(playerHand Hand, blackjackDeck *BlackjackDeck, bet USD, player *Player) (Hand, USD) {
	var playerChoice int
	for playerChoice != 2 {
		fmt.Println("Would you like to:")
		fmt.Println("1. Hit")
		fmt.Println("2. Stand")
		if len(playerHand) == 2 && player.USD >= bet.Multiply(2) {
			fmt.Println("3. Double down")
		}
		fmt.Scanln(&playerChoice)
		switch playerChoice {
		case 1:
			playerHand = append(playerHand, blackjackDeck.Deal())
			fmt.Printf("Your hand: %s\n", playerHand.String())
		case 3:
			//TODO can only do if has enough money, and only on first run
			playerHand = append(playerHand, blackjackDeck.Deal())
			bet = bet.Multiply(2)
			fmt.Printf("Bet is now %s\n", bet.String())
			fmt.Printf("Your hand: %s\n", playerHand.String())
			return playerHand, bet
		}
		if playerHand.Score().Value > 21 {
			return playerHand, bet
		}
	}
	return playerHand, bet
}

func runDealer(dealerHand Hand, blackjackDeck *BlackjackDeck) Hand {
	for dealerHand.Score().Value <= 16 || dealerHand.Score().IsSoft17() {
		dealerHand = append(dealerHand, blackjackDeck.Deal())
	}
	return dealerHand
}
