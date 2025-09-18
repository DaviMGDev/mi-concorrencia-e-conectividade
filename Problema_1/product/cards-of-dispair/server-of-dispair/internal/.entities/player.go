package entities

import "sync"

type Player struct {
	ID   string  `json:"id"`
	HP   uint    `json:"hp"`
	Hand [3]Card `json:"hand"`

	mutex sync.RWMutex `json:"-"`
}

type PlayerInterface interface {
	TakeDamage(damage uint)
	Reset()
	GetHP() uint
	GetHand() [3]Card
	ChangeHand(hand [3]Card)
}

func NewPlayer(id string) *Player {
	player := &Player{
		ID: id,
	}
	player.Reset()
	return player
}

func (p *Player) TakeDamage(damage uint) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if damage >= p.HP {
		p.HP = 0
	} else {
		p.HP -= damage
	}
}

func (p *Player) Reset() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.HP = 10
	p.Hand = [3]Card{
		NewCard("rock", 1),
		NewCard("paper", 1),
		NewCard("scissors", 1),
	}
}

func (p *Player) ChangeHand(hand [3]Card) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Hand = hand
}

func (p *Player) GetHP() uint {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.HP
}

func (p *Player) GetHand() [3]Card {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.Hand
}