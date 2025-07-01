package services

import (
	"math/rand"
)

// defaultRandomNumberGenerator は乱数生成処理のデフォルト実装です。
type defaultRandomNumberGenerator struct{}

// NewRandomNumberGenerator は defaultRandomNumberGenerator のインスタンスを返します。
// 呼び出し側では handlers.RandomNumberGenerator インターフェースとして受け取る想定です。
func NewRandomNumberGenerator() *defaultRandomNumberGenerator {
	return &defaultRandomNumberGenerator{}
}

// Generate は 1 〜 13 の乱数を返します。
func (g *defaultRandomNumberGenerator) Generate() int {
	return rand.Intn(13) + 1
}
