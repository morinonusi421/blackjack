package services

import (
	"testing"

	"blackjack/api/game"
)

// mockDeck はテスト用に決められたカードを順番に配るデッキです
type mockDeck struct {
	cards []game.Card
}

func (m *mockDeck) Deal() game.Card {
	c := m.cards[0]
	m.cards = m.cards[1:]
	return c
}

func TestGameService_NewGame_InvalidBet(t *testing.T) {
	svc := NewGameService(&mockDeck{})
	_, err := svc.NewGame(0)
	if err == nil {
		t.Fatalf("expected error for non-positive bet, got nil")
	}
}

func TestGameService_NewGame_ValidBet(t *testing.T) {
	bet := 100

	cases := []struct {
		name              string
		deckCards         []game.Card
		expectState       game.GameState
		expectResult      game.Result
		expectBalanceDiff int
	}{
		{
			name: "blackjack",
			deckCards: []game.Card{
				{Suit: game.Spade, Rank: "A"},
				{Suit: game.Heart, Rank: "10"}, // プレイヤー 21
				{Suit: game.Club, Rank: "9"},   // ディーラー
			},
			expectState:       game.Finished,
			expectResult:      game.PlayerWin,
			expectBalanceDiff: bet * 3 / 2,
		},
		{
			name: "non-blackjack",
			deckCards: []game.Card{
				{Suit: game.Spade, Rank: "9"},
				{Suit: game.Heart, Rank: "7"}, // プレイヤー 16
				{Suit: game.Club, Rank: "5"},  // ディーラー
			},
			expectState:       game.PlayerTurn,
			expectResult:      game.Pending,
			expectBalanceDiff: 0,
		},
	}

	for _, tc := range cases {
		deck := &mockDeck{cards: tc.deckCards}
		svc := NewGameService(deck)

		g, err := svc.NewGame(bet)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.name, err)
		}

		// 基本プロパティ
		if g.Bet != bet {
			t.Fatalf("%s: expected bet %d, got %d", tc.name, bet, g.Bet)
		}
		if len(g.PlayerHand.Cards) != 2 || len(g.DealerHand.Cards) != 1 {
			t.Fatalf("%s: unexpected card counts p=%d d=%d", tc.name, len(g.PlayerHand.Cards), len(g.DealerHand.Cards))
		}

		// 状態/結果
		if g.State != tc.expectState || g.Result != tc.expectResult {
			t.Fatalf("%s: expected state=%s result=%s, got state=%s result=%s", tc.name, tc.expectState, tc.expectResult, g.State, g.Result)
		}

		if g.BalanceChange != tc.expectBalanceDiff {
			t.Fatalf("%s: expected balance change %d, got %d", tc.name, tc.expectBalanceDiff, g.BalanceChange)
		}
	}
}
