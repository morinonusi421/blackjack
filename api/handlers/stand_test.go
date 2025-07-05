package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blackjack/api/game"
)

type mockStandService struct{}

func (m mockStandService) NewGame(bet int) (game.Game, error) { // 未使用
	return game.Game{}, nil
}

func (m mockStandService) Stand(g *game.Game) error {
	// ダミーで引き分けを返す
	g.State = game.Finished
	g.Result = game.Push
	g.Payout = g.Bet
	return nil
}

func TestStandHandler_ReturnsUpdatedGameJSON(t *testing.T) {
	bet := 100

	// Push シナリオ: プレイヤー18, ディーラー18
	playerCards := []game.Card{{Suit: game.Spade, Rank: "9"}, {Suit: game.Heart, Rank: "9"}}
	dealerCards := []game.Card{{Suit: game.Club, Rank: "9"}, {Suit: game.Diamond, Rank: "9"}}

	g := game.Game{
		PlayerHand: game.Hand{Cards: playerCards, Score: game.CalculateScore(playerCards)},
		DealerHand: game.Hand{Cards: dealerCards, Score: game.CalculateScore(dealerCards)},
		Bet:        bet,
		State:      game.PlayerTurn,
		Result:     game.Pending,
	}

	svc := mockStandService{}

	handler := StandHandler(svc)

	body, _ := json.Marshal(g)
	req := httptest.NewRequest(http.MethodPost, "/api/game/stand", bytes.NewReader(body))
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

	if resp.Result != game.Push {
		t.Fatalf("expected result Push, got %s", resp.Result)
	}

	if resp.Payout != bet {
		t.Fatalf("expected payout %d, got %d", bet, resp.Payout)
	}
}
