package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blackjack/api/game"
)

// mockGameService はテスト用に固定の Game を返すサービス実装です。
type mockGameService struct {
	expectedBet int
	retGame     game.Game
	retErr      error
}

func (m mockGameService) NewGame(bet int) (game.Game, error) {
	if bet != m.expectedBet {
		return game.Game{}, m.retErr
	}
	return m.retGame, m.retErr
}

// ダミー実装: テスト対象外
func (m mockGameService) Stand(g *game.Game) error {
	return nil
}

// 未使用
func (m mockGameService) Hit(g *game.Game) error {
	return nil
}

func TestNewGameHandler_ReturnsGameJSON(t *testing.T) {
	// 期待する Game オブジェクト
	expectedBet := 100
	g := game.Game{
		Bet:    expectedBet,
		State:  game.PlayerTurn,
		Result: game.Pending,
		Payout: 0,
	}

	// モックサービス
	svc := mockGameService{
		expectedBet: expectedBet,
		retGame:     g,
		retErr:      nil,
	}

	handler := NewGameHandler(svc)

	// リクエストボディ
	body, _ := json.Marshal(NewGameRequest{Bet: expectedBet})
	req := httptest.NewRequest(http.MethodPost, "/api/game/new", bytes.NewReader(body))
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

	if resp.Bet != expectedBet || resp.State != game.PlayerTurn || resp.Result != game.Pending {
		t.Fatalf("unexpected response: %+v", resp)
	}
}
