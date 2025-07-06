package handlers

import (
	"encoding/json"
	"net/http"

	"blackjack/api/game"
	"blackjack/api/services"
)

// SurrenderRequest はサレンダー時にクライアントから送られてくる現在のゲーム状態を表します。
type SurrenderRequest game.Game

// SurrenderHandler は GameService の Surrender を呼び出すハンドラを返します。
func SurrenderHandler(gameSvc services.GameService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req SurrenderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		g := game.Game(req)

		if err := gameSvc.Surrender(&g); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(g)
	}
}
