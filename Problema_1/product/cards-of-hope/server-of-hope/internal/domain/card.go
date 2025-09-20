package domain

// Card representa uma carta do jogo.
//
// Campos:
//   - Type: tipo da carta (ex: pedra, papel, tesoura).
//   - Stars: quantidade de estrelas da carta.
type Card struct {
	Type  string `json:"type"`
	Stars int    `json:"stars"`
}

// CardPackage representa um pacote de três cartas.
type CardPackage [3]Card
