package handlers

import (
	"encoding/json"
	"net/http"

	"blackjack/api/game"
	"blackjack/api/services"
)

// StrategyRequest は現在のゲーム状態と設定を入力として受け取る
type StrategyRequest struct {
	Game   game.Game       `json:"game"`
	Config game.GameConfig `json:"config"`
}

// StrategyResponse は各アクションの期待払い戻しを返す
type StrategyResponse struct {
	HitPayout       float64 `json:"hit_payout"`
	StandPayout     float64 `json:"stand_payout"`
	SurrenderPayout float64 `json:"surrender_payout"`
}

// StrategyHandler は最適戦略の期待払い戻しを返すハンドラ
func StrategyHandler(strategyAdvisor services.StrategyAdvisor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req StrategyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		// 不正なconfigじゃないかバリデーション
		if req.Config.DealerStandThreshold < 1 || req.Config.DealerStandThreshold > 21 {
			http.Error(w, "dealer stand threshold must be between 1 and 21", http.StatusBadRequest)
			return
		}

		payouts, err := strategyAdvisor.Advise(req.Game, &req.Config)
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
