package service

import (
	"math/rand"
)

// RandomOneToThirteen は1〜13の乱数を返します。
func RandomOneToThirteen() int {
	return rand.Intn(13) + 1
}
