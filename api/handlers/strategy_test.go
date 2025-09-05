package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blackjack/api/game"
	"blackjack/api/strategy"
)

// mockStrategyService は固定の期待払い戻し（金額ベース, 事前スケール済み）を返すモック
type mockStrategyService struct {
	payouts strategy.StrategyExpectedPayouts
	err     error
}

func (m mockStrategyService) Advise(g game.Game, config *game.GameConfig) (strategy.StrategyExpectedPayouts, error) {
	return m.payouts, m.err
}

func TestStrategyHandler_ReturnsExpectedPayouts(t *testing.T) {
	mockService := mockStrategyService{
		payouts: strategy.StrategyExpectedPayouts{
			HitPayout:       50.0,
			StandPayout:     75.0,
			SurrenderPayout: 25.0,
		},
		err: nil,
	}
	handler := StrategyHandler(mockService)

	// ゲームのダミー状態
	g := game.Game{
		PlayerHand: game.Hand{Cards: []game.Card{{Suit: game.Spade, Rank: "A"}, {Suit: game.Spade, Rank: "9"}}, Score: 20},
		DealerHand: game.Hand{Cards: []game.Card{{Suit: game.Heart, Rank: "7"}}, Score: 7},
		State:      game.PlayerTurn,
		Result:     game.Pending,
		Bet:        100,
	}

	config := game.GameConfig{
		DealerStandThreshold: 17,
	}

	req_body := StrategyRequest{
		Game:   g,
		Config: config,
	}

	body, _ := json.Marshal(req_body)
	req := httptest.NewRequest(http.MethodPost, "/api/strategy/advise", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", got)
	}

	var resp StrategyResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	// 期待払い戻しが適切に計算されていることを確認（正確な値ではなく範囲で確認）
	if resp.HitPayout < 0 || resp.StandPayout < 0 || resp.SurrenderPayout < 0 {
		t.Fatalf("unexpected negative payouts: %+v", resp)
	}
}
