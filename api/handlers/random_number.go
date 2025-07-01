package handlers

import (
	"encoding/json"
	"net/http"
)

// RandomNumberResponse は乱数を JSON で返すレスポンス構造体
type RandomNumberResponse struct {
	Number int `json:"number"`
}

// RandomNumberGenerator は乱数を返す機能を定義するインターフェースです。
type RandomNumberGenerator interface {
	Generate() int
}

// NewRandomNumberHandler は依存として受け取った RandomNumberGenerator を使い乱数を返すハンドラを生成します。
func NewRandomNumberHandler(generator RandomNumberGenerator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resp := RandomNumberResponse{Number: generator.Generate()}
		json.NewEncoder(w).Encode(resp)
	}
}
