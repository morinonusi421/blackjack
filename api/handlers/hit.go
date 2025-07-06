package handlers

import (
	"encoding/json"
	"net/http"

	"blackjack/api/game"
	"blackjack/api/services"
)

// HitRequest はヒット時にクライアントから送られてくる現在のゲーム状態を表します。
type HitRequest game.Game

// HitHandler は Hitter の Hit を呼び出すハンドラを返します。
func HitHandler(gameSvc services.Hitter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req HitRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		g := game.Game(req)

		if err := gameSvc.Hit(&g); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(g)
	}
}
