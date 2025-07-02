package game

import "math/rand"

// Deck はカードの供給源を表すインターフェースです。
// Deal() で 1 枚カードを返します。
type Deck interface {
	Deal() Card
}

// RandomDeck はランダムにカードを配る本番用のデッキです。
type RandomDeck struct{}

// Deal はランダムに 1 枚のカードを返します。
func (d *RandomDeck) Deal() Card {
	return Card{
		Suit: suits[rand.Intn(len(suits))],
		Rank: ranks[rand.Intn(len(ranks))],
	}
}
