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

// mockStrategyAdvisor は固定の期待払い戻し（金額ベース, 事前スケール済み）を返すモック
type mockStrategyAdvisor struct {
	payouts strategy.StrategyExpectedPayouts
	err     error
}

func (m mockStrategyAdvisor) Advise(g game.Game) (strategy.StrategyExpectedPayouts, error) {
	return m.payouts, m.err
}

func TestStrategyHandler_ReturnsExpectedPayouts(t *testing.T) {
	advisor := mockStrategyAdvisor{
		payouts: strategy.StrategyExpectedPayouts{
			HitPayout:       90.0,
			StandPayout:     80.0,
			SurrenderPayout: 50.0,
			BestPayout:      90.0,
		},
		err: nil,
	}

	handler := StrategyHandler(advisor)

	// ゲームのダミー状態
	g := game.Game{
		PlayerHand: game.Hand{Cards: []game.Card{{Suit: game.Spade, Rank: "A"}, {Suit: game.Spade, Rank: "9"}}, Score: 20},
		DealerHand: game.Hand{Cards: []game.Card{{Suit: game.Heart, Rank: "7"}}, Score: 7},
		State:      game.PlayerTurn,
		Result:     game.Pending,
		Bet:        1,
	}

	body, _ := json.Marshal(StrategyRequest(g))
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

	// 金額（bet倍）で返る（このテストではbet=1なので固定値のまま）
	if resp.HitPayout != 90.0 || resp.StandPayout != 80.0 || resp.SurrenderPayout != 50.0 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}
