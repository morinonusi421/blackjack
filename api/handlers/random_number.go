package handlers

import (
	"encoding/json"
	"net/http"

	"blackjack/api/service"
)

// RandomNumberResponse は乱数を JSON で返すレスポンス構造体
type RandomNumberResponse struct {
	Number int `json:"number"`
}

// RandomNumberHandler は1〜3の乱数を返します。
func RandomNumberHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := RandomNumberResponse{Number: service.RandomOneToThirteen()}
	json.NewEncoder(w).Encode(resp)
}
