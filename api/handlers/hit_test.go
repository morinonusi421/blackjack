package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blackjack/api/game"
)

// mockHitService は Hit の挙動をテストするためのモックサービスです。
type mockHitService struct{}

func (m mockHitService) Hit(g *game.Game, config *game.GameConfig) error {
	// テスト用に単純にプレイヤーバーストさせる
	g.State = game.Finished
	g.Result = game.DealerWin
	g.Payout = 0
	return nil
}

func TestHitHandler_ReturnsUpdatedGameJSON(t *testing.T) {
	bet := 100

	// 任意のゲーム状態
	playerCards := []game.Card{{Suit: game.Spade, Rank: "10"}, {Suit: game.Heart, Rank: "8"}}
	dealerCards := []game.Card{{Suit: game.Club, Rank: "9"}}

	g := game.Game{
		PlayerHand: game.Hand{Cards: playerCards, Score: game.CalculateScore(playerCards)},
		DealerHand: game.Hand{Cards: dealerCards, Score: game.CalculateScore(dealerCards)},
		Bet:        bet,
		State:      game.PlayerTurn,
		Result:     game.Pending,
	}

	svc := mockHitService{}

	handler := HitHandler(svc)

	req_body := HitRequest{
		Game:   g,
		Config: game.GameConfig{DealerStandThreshold: 17},
	}
	body, _ := json.Marshal(req_body)
	req := httptest.NewRequest(http.MethodPost, "/api/game/hit", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", got)
	}

	var resp game.Game
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Result != game.DealerWin {
		t.Fatalf("expected result DealerWin, got %s", resp.Result)
	}

	if resp.State != game.Finished {
		t.Fatalf("expected state Finished, got %s", resp.State)
	}
}
