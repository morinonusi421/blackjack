package service

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomOneToThree は1〜3の乱数を返します。
func RandomOneToThree() int {
	return rand.Intn(3) + 1
}
