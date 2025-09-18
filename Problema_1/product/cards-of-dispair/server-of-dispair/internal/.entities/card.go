// TO CHECK
package entities

var (
	CardTypes = []string{
		"rock",
		"paper",
		"scissors",
		"surrender",
	}
	CardWins = map[string]string{
		"rock":     "scissors",
		"paper":    "rock",
		"scissors": "paper",
	}
)

type Card struct {
	Type  string `json:"type"`
	Stars uint   `json:"stars"`
}

func NewCard(cardType string, stars uint) Card {
	return Card{
		Type:  cardType,
		Stars: stars,
	}
}

type CardInterface interface {
	Against(opponent Card) int
	GetType() string
	GetStars() uint
}

func (c Card) Against(opponent Card) int {
	if c.Type == "surrender" {
		return -1 // Lose
	}
	if opponent.Type == "surrender" {
		return 1 // Win
	}
	if c.Type == opponent.Type {
		if c.Stars > opponent.Stars {
			return 1 // Win
		} else if c.Stars < opponent.Stars {
			return -1 // Lose
		}
		return 0 // Draw
	}
	if CardWins[c.Type] == opponent.Type {
		return 1 // Win
	}
	return -1 // Lose
}

func (c Card) GetType() string {
	return c.Type
}

func (c Card) GetStars() uint {
	return c.Stars
}
