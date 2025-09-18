package entities

import "sync"

type Game struct {
	ID      string          `json:"id"`
	RoomID  string          `json:"room_id"`
	Choice  map[string]Card `json:"choice"` // playerID -> cardType
	started bool            `json:"started"`

	mutex sync.RWMutex `json:"-"`
}

type GameInterface interface {
	AddPlayer(playerID string) error
	Started() bool
	AllReady() bool
	MakeChoice(playerID, card string, stars uint) error
	ResetChoices() error
	Start()
	End()
}

func NewGame(id, roomID string) *Game {
	return &Game{
		ID:      id,
		RoomID:  roomID,
		Choice:  make(map[string]Card),
		started: false,
	}
}

func (g *Game) Started() bool {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.started
}

func (g *Game) AllReady() bool {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return len(g.Choice) >= 2 // Assuming a minimum of 2 players to start
}

func (g *Game) ResetChoices() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.Choice = make(map[string]Card)
	return nil
}

func (g *Game) AddPlayer(playerID string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if _, exists := g.Choice[playerID]; exists {
		return nil // Player already in the game
	}
	g.Choice[playerID] = Card{} // Initialize with an empty card
	return nil
}

func (g *Game) MakeChoice(playerID, cardType string, stars uint) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if _, exists := g.Choice[playerID]; !exists {
		return nil // Player not in the game
	}
	g.Choice[playerID] = NewCard(cardType, stars)
	return nil
}

func (g *Game) Start() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.started = true
}

func (g *Game) End() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.started = false
	g.Choice = make(map[string]Card) // Clear choices
}