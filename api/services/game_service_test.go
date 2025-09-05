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

// MockDeck is an improved mock deck for more complex testing scenarios
type MockDeck struct {
	nextCard game.Card
	calls    int
}

func (m *MockDeck) Deal() game.Card {
	m.calls++
	return m.nextCard
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
		name         string
		deckCards    []game.Card
		expectState  game.GameState
		expectResult game.Result
		expectPayout int
	}{
		{
			name: "blackjack",
			deckCards: []game.Card{
				{Suit: game.Spade, Rank: "A"},
				{Suit: game.Heart, Rank: "10"}, // プレイヤー 21
				{Suit: game.Club, Rank: "9"},   // ディーラー
			},
			expectState:  game.Finished,
			expectResult: game.PlayerWin,
			expectPayout: 250,
		},
		{
			name: "non-blackjack",
			deckCards: []game.Card{
				{Suit: game.Spade, Rank: "9"},
				{Suit: game.Heart, Rank: "7"}, // プレイヤー 16
				{Suit: game.Club, Rank: "5"},  // ディーラー
			},
			expectState:  game.PlayerTurn,
			expectResult: game.Pending,
			expectPayout: 0,
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

		if g.Payout != tc.expectPayout {
			t.Fatalf("%s: expected payout %d, got %d", tc.name, tc.expectPayout, g.Payout)
		}
	}
}

func TestGameService_Stand(t *testing.T) {
	bet := 100

	type scenario struct {
		name         string
		playerCards  []game.Card // プレイヤー初期手札
		dealerCards  []game.Card // ディーラー初期手札 (1 枚を想定)
		deckCards    []game.Card // ディーラーがヒットで引くカード（必要に応じて）
		expectResult game.Result
		expectPayout int
	}

	cases := []scenario{
		{
			name:         "dealer bust -> player win",
			playerCards:  []game.Card{{Suit: game.Spade, Rank: "10"}, {Suit: game.Heart, Rank: "7"}},  // 17
			dealerCards:  []game.Card{{Suit: game.Club, Rank: "10"}},                                  // 10 (<17)
			deckCards:    []game.Card{{Suit: game.Diamond, Rank: "6"}, {Suit: game.Spade, Rank: "8"}}, // 10+6=16, then 24→bust
			expectResult: game.PlayerWin,
			expectPayout: bet * 2,
		},
		{
			name:         "dealer higher score -> dealer win",
			playerCards:  []game.Card{{Suit: game.Spade, Rank: "10"}, {Suit: game.Heart, Rank: "6"}}, // 16
			dealerCards:  []game.Card{{Suit: game.Club, Rank: "10"}},                                 // 10
			deckCards:    []game.Card{{Suit: game.Diamond, Rank: "8"}},                               // 18
			expectResult: game.DealerWin,
			expectPayout: 0,
		},
		{
			name:         "push",
			playerCards:  []game.Card{{Suit: game.Spade, Rank: "10"}, {Suit: game.Heart, Rank: "8"}}, // 18
			dealerCards:  []game.Card{{Suit: game.Club, Rank: "9"}},                                  // 9
			deckCards:    []game.Card{{Suit: game.Diamond, Rank: "9"}},                               // 18 tie
			expectResult: game.Push,
			expectPayout: bet, // 掛け金返却
		},
		{
			name:         "player higher score -> player win",
			playerCards:  []game.Card{{Suit: game.Spade, Rank: "10"}, {Suit: game.Heart, Rank: "9"}}, // 19
			dealerCards:  []game.Card{{Suit: game.Club, Rank: "10"}},                                 // 10
			deckCards:    []game.Card{{Suit: game.Diamond, Rank: "8"}},                               // 18 (<21, >=17)
			expectResult: game.PlayerWin,
			expectPayout: bet * 2,
		},
	}

	for _, tc := range cases {
		deck := &mockDeck{cards: tc.deckCards}
		svc := NewGameService(deck)

		// ゲーム状態を構築
		g := game.Game{
			PlayerHand: game.Hand{Cards: tc.playerCards, Score: game.CalculateScore(tc.playerCards)},
			DealerHand: game.Hand{Cards: tc.dealerCards, Score: game.CalculateScore(tc.dealerCards)},
			Bet:        bet,
			State:      game.PlayerTurn,
			Result:     game.Pending,
		}

		config := &game.GameConfig{DealerStandThreshold: 17}
		if err := svc.Stand(&g, config); err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.name, err)
		}

		if g.Result != tc.expectResult {
			t.Fatalf("%s: expected result %s, got %s", tc.name, tc.expectResult, g.Result)
		}

		if g.Payout != tc.expectPayout {
			t.Fatalf("%s: expected payout %d, got %d", tc.name, tc.expectPayout, g.Payout)
		}
	}
}

func TestGameService_Hit(t *testing.T) {
	bet := 100

	type scenario struct {
		name         string
		playerCards  []game.Card
		dealerCards  []game.Card
		deckCards    []game.Card // Hit で引くカード、その後 Stand で使うカード
		expectState  game.GameState
		expectResult game.Result
		expectPayout int
	}

	cases := []scenario{
		{
			name:         "player busts",
			playerCards:  []game.Card{{Suit: game.Spade, Rank: "10"}, {Suit: game.Heart, Rank: "2"}}, // 12
			dealerCards:  []game.Card{{Suit: game.Club, Rank: "9"}},                                  // 9
			deckCards:    []game.Card{{Suit: game.Diamond, Rank: "K"}},                               // 12+10 = 22 -> bust
			expectState:  game.Finished,
			expectResult: game.DealerWin,
			expectPayout: 0,
		},
		{
			name:        "hit to 21 -> player win after dealer play",
			playerCards: []game.Card{{Suit: game.Spade, Rank: "5"}, {Suit: game.Heart, Rank: "6"}}, // 11
			dealerCards: []game.Card{{Suit: game.Club, Rank: "10"}},                                // 10
			deckCards: []game.Card{
				{Suit: game.Diamond, Rank: "K"}, // player hits to 21
				{Suit: game.Spade, Rank: "7"},   // dealer hits to 17 and stands
			},
			expectState:  game.Finished,
			expectResult: game.PlayerWin,
			expectPayout: bet * 2,
		},
		{
			name:         "hit under 21 continue",
			playerCards:  []game.Card{{Suit: game.Spade, Rank: "7"}, {Suit: game.Heart, Rank: "5"}}, // 12
			dealerCards:  []game.Card{{Suit: game.Club, Rank: "10"}},                                // 10
			deckCards:    []game.Card{{Suit: game.Diamond, Rank: "6"}},                              // player 18 (<21)
			expectState:  game.PlayerTurn,
			expectResult: game.Pending,
			expectPayout: 0,
		},
	}

	for _, tc := range cases {
		deck := &mockDeck{cards: tc.deckCards}
		svc := NewGameService(deck)

		g := game.Game{
			PlayerHand: game.Hand{Cards: tc.playerCards, Score: game.CalculateScore(tc.playerCards)},
			DealerHand: game.Hand{Cards: tc.dealerCards, Score: game.CalculateScore(tc.dealerCards)},
			Bet:        bet,
			State:      game.PlayerTurn,
			Result:     game.Pending,
		}

		config := &game.GameConfig{DealerStandThreshold: 17}
		if err := svc.Hit(&g, config); err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.name, err)
		}

		if g.State != tc.expectState {
			t.Fatalf("%s: expected state %s, got %s", tc.name, tc.expectState, g.State)
		}
		if g.Result != tc.expectResult {
			t.Fatalf("%s: expected result %s, got %s", tc.name, tc.expectResult, g.Result)
		}
		if g.Payout != tc.expectPayout {
			t.Fatalf("%s: expected payout %d, got %d", tc.name, tc.expectPayout, g.Payout)
		}
	}
}

func TestGameService_Surrender(t *testing.T) {
	bet := 100

	deck := &mockDeck{cards: []game.Card{}}
	svc := NewGameService(deck)

	g := game.Game{
		PlayerHand: game.Hand{
			Cards: []game.Card{
				{Suit: game.Spade, Rank: "10"},
				{Suit: game.Heart, Rank: "6"},
			},
			Score: 16,
		},
		DealerHand: game.Hand{
			Cards: []game.Card{{Suit: game.Club, Rank: "9"}},
			Score: 9,
		},
		Bet:    bet,
		State:  game.PlayerTurn,
		Result: game.Pending,
	}

	config := &game.GameConfig{DealerStandThreshold: 17}
	err := svc.Surrender(&g, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if g.State != game.Finished {
		t.Fatalf("expected state %s, got %s", game.Finished, g.State)
	}
	if g.Result != game.Surrender {
		t.Fatalf("expected result %s, got %s", game.Surrender, g.Result)
	}
	if g.Payout != bet/2 {
		t.Fatalf("expected payout %d, got %d", bet/2, g.Payout)
	}
	if g.ResultMessage != game.MessagePlayerSurrendered {
		t.Fatalf("expected result message '%s', got '%s'", game.MessagePlayerSurrendered, g.ResultMessage)
	}
}
