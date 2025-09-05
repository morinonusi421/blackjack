package handlers

import (
	"encoding/json"
	"net/http"

	"blackjack/api/game"
	"blackjack/api/services"
)

// StandRequest はスタンド時にクライアントから送られてくる現在のゲーム状態と設定を表します。
type StandRequest struct {
	Game   game.Game       `json:"game"`
	Config game.GameConfig `json:"config"`
}

// StandHandler は Stander の Stand を呼び出すハンドラを返します。
func StandHandler(gameSvc services.Stander) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req StandRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		g := req.Game

		if err := gameSvc.Stand(&g, &req.Config); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(g)
	}
}
