package ai

import blackjack "gophercises/blackjack/game"

type playerHand struct {
	score int
	soft  bool
	pair  bool
}
type chooseF func(allowedMoves []blackjack.Move) blackjack.Move

var (
	hitF   = func(allowedMoves []blackjack.Move) blackjack.Move { return blackjack.Hit }
	standF = func(allowedMoves []blackjack.Move) blackjack.Move { return blackjack.Stand }
	splitF = func(allowedMoves []blackjack.Move) blackjack.Move { return blackjack.Split }
	dhF    = func(allowedMoves []blackjack.Move) blackjack.Move {
		if containsMove(allowedMoves, blackjack.Double) {
			return blackjack.Double
		}
		return blackjack.Hit
	}
	dsF = func(allowedMoves []blackjack.Move) blackjack.Move {
		if containsMove(allowedMoves, blackjack.Double) {
			return blackjack.Double
		}
		return blackjack.Stand
	}
	suhF = func(allowedMoves []blackjack.Move) blackjack.Move {
		if containsMove(allowedMoves, blackjack.Surrender) {
			return blackjack.Surrender
		}
		return blackjack.Hit
	}
	susF = func(allowedMoves []blackjack.Move) blackjack.Move {
		if containsMove(allowedMoves, blackjack.Surrender) {
			return blackjack.Surrender
		}
		return blackjack.Stand
	}
	supF = func(allowedMoves []blackjack.Move) blackjack.Move {
		if containsMove(allowedMoves, blackjack.Surrender) {
			return blackjack.Surrender
		}
		return blackjack.Split
	}
	allHit   = map[int]chooseF{2: hitF, 3: hitF, 4: hitF, 5: hitF, 6: hitF, 7: hitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF}
	allStand = map[int]chooseF{2: standF, 3: standF, 4: standF, 5: standF, 6: standF, 7: standF, 8: standF, 9: standF, 10: standF, 11: standF}
)
var lookup = map[playerHand]map[int]chooseF{
	//hard hands
	playerHand{4, false, false}:  allHit,
	playerHand{5, false, false}:  allHit,
	playerHand{6, false, false}:  allHit,
	playerHand{7, false, false}:  allHit,
	playerHand{8, false, false}:  allHit,
	playerHand{9, false, false}:  {2: hitF, 3: dhF, 4: dhF, 5: dhF, 6: dhF, 7: hitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{10, false, false}: {2: dhF, 3: dhF, 4: dhF, 5: dhF, 6: dhF, 7: dhF, 8: dhF, 9: dhF, 10: hitF, 11: hitF,},
	playerHand{11, false, false}: {2: dhF, 3: dhF, 4: dhF, 5: dhF, 6: dhF, 7: dhF, 8: dhF, 9: dhF, 10: dhF, 11: dhF,},
	playerHand{12, false, false}: {2: hitF, 3: hitF, 4: standF, 5: standF, 6: standF, 7: hitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{13, false, false}: {2: standF, 3: standF, 4: standF, 5: standF, 6: standF, 7: hitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{14, false, false}: {2: standF, 3: standF, 4: standF, 5: standF, 6: standF, 7: hitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{15, false, false}: {2: standF, 3: standF, 4: standF, 5: standF, 6: standF, 7: hitF, 8: hitF, 9: hitF, 10: suhF, 11: suhF,},
	playerHand{16, false, false}: {2: standF, 3: standF, 4: standF, 5: standF, 6: standF, 7: hitF, 8: hitF, 9: suhF, 10: suhF, 11: suhF,},
	playerHand{17, false, false}: {2: standF, 3: standF, 4: standF, 5: standF, 6: standF, 7: standF, 8: standF, 9: standF, 10: standF, 11: susF,},
	playerHand{18, false, false}: allStand,
	playerHand{19, false, false}: allStand,
	playerHand{20, false, false}: allStand,
	playerHand{21, false, false}: allStand,
	//soft hands
	playerHand{13, true, false}: {2: hitF, 3: hitF, 4: hitF, 5: dhF, 6: dhF, 7: hitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{14, true, false}: {2: hitF, 3: hitF, 4: hitF, 5: dhF, 6: dhF, 7: hitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{15, true, false}: {2: hitF, 3: hitF, 4: dhF, 5: dhF, 6: dhF, 7: hitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{16, true, false}: {2: hitF, 3: hitF, 4: dhF, 5: dhF, 6: dhF, 7: hitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{17, true, false}: {2: hitF, 3: dhF, 4: dhF, 5: dhF, 6: dhF, 7: hitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{18, true, false}: {2: dsF, 3: dsF, 4: dsF, 5: dsF, 6: dsF, 7: standF, 8: standF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{19, true, false}: {2: standF, 3: standF, 4: standF, 5: standF, 6: dsF, 7: standF, 8: standF, 9: standF, 10: standF, 11: standF,},
	playerHand{20, true, false}: allStand,
	playerHand{21, true, false}: allStand,
	//pairs
	playerHand{4, false, true}:  {2: splitF, 3: splitF, 4: splitF, 5: splitF, 6: splitF, 7: splitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{6, false, true}:  {2: splitF, 3: splitF, 4: splitF, 5: splitF, 6: splitF, 7: splitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{8, false, true}:  {2: hitF, 3: hitF, 4: hitF, 5: splitF, 6: splitF, 7: hitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{10, false, true}: {2: dhF, 3: dhF, 4: dhF, 5: dhF, 6: dhF, 7: dhF, 8: dhF, 9: dhF, 10: hitF, 11: hitF,},
	playerHand{12, false, true}: {2: splitF, 3: splitF, 4: splitF, 5: splitF, 6: splitF, 7: hitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{14, false, true}: {2: splitF, 3: splitF, 4: splitF, 5: splitF, 6: splitF, 7: splitF, 8: hitF, 9: hitF, 10: hitF, 11: hitF,},
	playerHand{16, false, true}: {2: splitF, 3: splitF, 4: splitF, 5: splitF, 6: splitF, 7: splitF, 8: splitF, 9: splitF, 10: splitF, 11: supF,},
	playerHand{18, false, true}: {2: splitF, 3: splitF, 4: splitF, 5: splitF, 6: splitF, 7: standF, 8: splitF, 9: splitF, 10: standF, 11: standF,},
	playerHand{20, false, true}: allStand,
	//ace pair
	playerHand{12, true, true}: {2: splitF, 3: splitF, 4: splitF, 5: splitF, 6: splitF, 7: splitF, 8: splitF, 9: splitF, 10: splitF, 11: splitF,},
}
