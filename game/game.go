package game

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
}

func (g GameState) GetID() string {
	return g.ID
}

func (g GameState) HasPlayer(player string) bool {
	return hasPlayer(g.Players, player)
}

// state of the game visible to a player during the game
type VisibleGameState struct {
	Hand          []Card `json:"hand"`
	Table         []Card `json:"table"`
	CurrentPlayer int    `json:"currentPlayer"`
	Trump         Color  `json:"trump"`
	Bid           int    `json:"bid"`
	BidWinner     int    `json:"bidWinner"`
}

func (g GameState) Visible(player int) VisibleGameState {
	return VisibleGameState{
		Hand:          g.Hands[player],
		Table:         g.Table,
		CurrentPlayer: g.CurrentPlayer,
		Trump:         g.Trump,
		Bid:           g.Bid,
		BidWinner:     g.BidWinner,
	}
}

func hasPlayer(players [4]string, player string) bool {
	for _, p := range players {
		if p == player {
			return true
		}
	}
	return false
}

type HasID interface {
	GetID() string
}
