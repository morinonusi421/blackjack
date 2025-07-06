package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blackjack/api/game"
)

type mockSurrenderService struct{}

func (m mockSurrenderService) NewGame(bet int) (game.Game, error) { // 未使用
	return game.Game{}, nil
}

func (m mockSurrenderService) Stand(g *game.Game) error { // 未使用
	return nil
}

func (m mockSurrenderService) Hit(g *game.Game) error { // 未使用
	return nil
}

func (m mockSurrenderService) Surrender(g *game.Game) error {
	// ダミーでサレンダー結果を返す
	g.State = game.Finished
	g.Result = game.Surrender
	g.ResultMessage = "Player surrendered."
	g.Payout = g.Bet / 2
	return nil
}

func TestSurrenderHandler_ReturnsUpdatedGameJSON(t *testing.T) {
	bet := 100

	// サレンダーシナリオ: プレイヤー16, ディーラー9
	playerCards := []game.Card{{Suit: game.Spade, Rank: "10"}, {Suit: game.Heart, Rank: "6"}}
	dealerCards := []game.Card{{Suit: game.Club, Rank: "9"}}

	g := game.Game{
		PlayerHand: game.Hand{Cards: playerCards, Score: game.CalculateScore(playerCards)},
		DealerHand: game.Hand{Cards: dealerCards, Score: game.CalculateScore(dealerCards)},
		Bet:        bet,
		State:      game.PlayerTurn,
		Result:     game.Pending,
	}

	svc := mockSurrenderService{}

	handler := SurrenderHandler(svc)

	body, _ := json.Marshal(g)
	req := httptest.NewRequest(http.MethodPost, "/api/game/surrender", bytes.NewReader(body))
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

	if resp.Result != game.Surrender {
		t.Fatalf("expected result Surrender, got %s", resp.Result)
	}

	if resp.Payout != bet/2 {
		t.Fatalf("expected payout %d, got %d", bet/2, resp.Payout)
	}

	if resp.ResultMessage != "Player surrendered." {
		t.Fatalf("expected result message 'Player surrendered.', got '%s'", resp.ResultMessage)
	}
}
