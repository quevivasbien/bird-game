package game

import (
	"fmt"

	"github.com/quevivasbien/bird-game/utils"
)

const GAME_ID_LENGTH = 8

type Color int

const (
	Red Color = iota + 1
	Yellow
	Green
	Black
)

type Card struct {
	Color Color `json:"color"`
	Value int   `json:"value"`
}

func (c Card) Beats(other Card, trump Color) bool {
	if c.Color != other.Color {
		if c.Color == trump || c == Bird {
			return true
		} else {
			return false
		}
	}
	return c.Value == 1 || c.Value > other.Value
}

func (c Card) Score() int {
	switch c.Value {
	case 1:
		return 15
	case 13:
		return 10
	case 10:
		return 10
	case 5:
		return 5
	case Bird.Value:
		return 20
	default:
		return 0
	}
}

var Bird Card = Card{0, 0}

type GameState struct {
	ID            string    `json:"id"`
	Players       [4]string `json:"players"`
	Hands         [4][]Card `json:"hands"`
	Discarded     [2][]Card `json:"discarded"`
	Widow         [5]Card   `json:"widow"`
	Table         []Card    `json:"table"`
	CurrentPlayer int       `json:"currentPlayer"`
	Trump         Color     `json:"trump"`
	Bid           int       `json:"bid"`
	BidWinner     int       `json:"bidWinner"`
	Done          bool      `json:"done"`
}

func (g GameState) GetID() string {
	return g.ID
}

func (g GameState) HasPlayer(player string) bool {
	return utils.Contains(g.Players[:], player)
}

func (g *GameState) ExchangeWithWidow(toWidow []Card, fromWidow []Card) error {
	if len(toWidow) != len(fromWidow) {
		return fmt.Errorf("Tried to take and give different amounts of cards from the widow")
	}
	newWidow := g.Widow // copy so we don't make changes if something is wrong
	newHand := g.Hands[g.BidWinner]
	for i := range toWidow {
		handIndex := utils.IndexOf(newHand, toWidow[i])
		if handIndex == -1 {
			return fmt.Errorf("Tried to put a card in the widow that was not in the bin winner's hand")
		}
		widowIndex := utils.IndexOf(newWidow[:], fromWidow[i])
		if widowIndex == -1 {
			return fmt.Errorf("Tried to take a card from the widow that was not in the widow")
		}
		newWidow[widowIndex], newHand[handIndex] = newHand[handIndex], newWidow[widowIndex]
	}
	g.Widow = newWidow
	g.Hands[g.BidWinner] = newHand
	return nil
}

func (g *GameState) PlayCard(player string, card Card) error {
	playerIndex := utils.IndexOf(g.Players[:], player)
	if playerIndex == -1 {
		return fmt.Errorf("Player is not in this game")
	}
	cards := g.Hands[playerIndex]
	cardIndex := utils.IndexOf(cards, card)
	if cardIndex == -1 {
		return fmt.Errorf("Card is not in player's hand")
	}
	g.Table = append(g.Table, card)
	g.Hands[playerIndex] = utils.Remove(g.Hands[playerIndex], cardIndex)
	if len(g.Table) == 4 {
		return g.finishPlay()
	}
	return nil
}

func (g *GameState) finishPlay() error {
	if len(g.Table) != 4 {
		return fmt.Errorf("Attempted to finish a play before all players have played")
	}
	// figure out winner
	winner := (g.CurrentPlayer + 1) % 4
	bestCard := g.Table[0]
	for i := 1; i <= 3; i++ {
		player := (winner + i) % 4
		card := g.Table[player]
		if card.Beats(bestCard, g.Trump) {
			winner = player
			bestCard = card
		}
	}
	// remove cards from table
	if winner%2 == 0 {
		g.Discarded[0] = append(g.Discarded[0], g.Table...)
	} else {
		g.Discarded[1] = append(g.Discarded[1], g.Table...)
	}
	g.Table = []Card{}
	// check if game is done
	done := true
	for i := range g.Hands {
		if len(g.Hands[i]) != 0 {
			done = false
		}
	}
	if done {
		g.Done = done
		if winner%2 == 0 {
			g.Discarded[0] = append(g.Discarded[0], g.Widow[:]...)
		} else {
			g.Discarded[1] = append(g.Discarded[1], g.Widow[:]...)
		}
	}
	return nil
}

func (g *GameState) Score() (int, int, error) {
	if !g.Done {
		return -1, -1, fmt.Errorf("Cannot calculate score before game is finished")
	}
	score0 := 0
	for _, card := range g.Discarded[0] {
		score0 += card.Score()
	}
	if len(g.Hands[0]) > len(g.Hands[1]) {
		score0 += 20
	}
	score1 := 200 - score0
	return score0, score1, nil
}

// state of the game visible to a player during the game
type VisibleGameState struct {
	Hand          []Card `json:"hand"`
	Table         []Card `json:"table"`
	CurrentPlayer int    `json:"currentPlayer"`
	Trump         Color  `json:"trump"`
	Bid           int    `json:"bid"`
	BidWinner     int    `json:"bidWinner"`
	Done          bool   `json:"done"`
}

func (g GameState) Visible(player int) VisibleGameState {
	return VisibleGameState{
		Hand:          g.Hands[player],
		Table:         g.Table,
		CurrentPlayer: g.CurrentPlayer,
		Trump:         g.Trump,
		Bid:           g.Bid,
		BidWinner:     g.BidWinner,
		Done:          g.Done,
	}
}
