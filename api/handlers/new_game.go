package handlers

import (
	"encoding/json"
	"net/http"

	"blackjack/api/services"
)

// NewGameRequest は新規ゲーム開始時に受け取るリクエストボディ
// 例: {"bet": 100}
// Bet は必須で 1 以上の整数であることを想定します。
type NewGameRequest struct {
	Bet int `json:"bet"`
}

// NewGameHandler は GameStarter を用いて新規ゲームを開始するハンドラを生成します。
func NewGameHandler(gameSvc services.GameStarter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// リクエストパース
		var req NewGameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		g, err := gameSvc.NewGame(req.Bet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(g)
	}
}
