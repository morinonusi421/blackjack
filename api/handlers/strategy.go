package handlers

import (
	"encoding/json"
	"net/http"

	"blackjack/api/game"
	"blackjack/api/services"
)

// StrategyRequest は現在のゲーム状態を入力として受け取る
type StrategyRequest game.Game

// StrategyResponse は各アクションの期待払い戻しを返す
type StrategyResponse struct {
	HitPayout       float64 `json:"hit_payout"`
	StandPayout     float64 `json:"stand_payout"`
	SurrenderPayout float64 `json:"surrender_payout"`
}

// StrategyHandler は最適戦略の期待払い戻しを返すハンドラ
func StrategyHandler(advisor services.StrategyAdvisor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req StrategyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		g := game.Game(req)

		payouts, err := advisor.Advise(g)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp := StrategyResponse{
			HitPayout:       payouts.HitPayout,
			StandPayout:     payouts.StandPayout,
			SurrenderPayout: payouts.SurrenderPayout,
		}

		json.NewEncoder(w).Encode(resp)
	}
}
