package strategy

import (
	"fmt"
	"testing"
)

func TestGetDealerScoreDistribution_Calculator(t *testing.T) {
	calc := NewCalculator()
	fmt.Println("ディラーが2を持つ場合の、ディラースコア分布:")
	dist := calc.GetDealerScoreDistribution(StrategyHand{Sum: 2, HasAce: false})
	fmt.Println(dist)

	fmt.Println("ディラーがAを持つ場合の、ディラースコア分布:")
	dist2 := calc.GetDealerScoreDistribution(StrategyHand{Sum: 1, HasAce: true})
	fmt.Println(dist2)
}

func TestCalculateStandExpectedPayout_Calculator(t *testing.T) {
	calc := NewCalculator()
	dealerHand := StrategyHand{Sum: 1, HasAce: true}
	standExpectedPayout := calc.CalculateStandExpectedPayout(18, dealerHand)
	fmt.Println("プレイヤーが18、ディーラーがAを持つ場合の、スタンド期待払い戻し:")
	fmt.Printf("StandExpectedPayout: %f\n", standExpectedPayout)
}

func TestCalculateAllExpectedPayouts_Calculator(t *testing.T) {
	calc := NewCalculator()
	// 例: プレイヤーA, ディーラー7, ヒットしていない
	state := StrategyState{
		Player: StrategyHand{Sum: 1, HasAce: true},
		Dealer: StrategyHand{Sum: 7, HasAce: false},
		HasHit: false,
	}
	expectedPayouts := calc.CalculateAllExpectedPayouts(state)
	fmt.Println("プレイヤーがAを持ち、ディーラーが7を持つ場合の、各アクションの期待払い戻し:")
	fmt.Printf("Hit: %f, Stand: %f, Surrender: %f, Best: %f\n", expectedPayouts.HitPayout, expectedPayouts.StandPayout, expectedPayouts.SurrenderPayout, expectedPayouts.BestPayout)

	// プレイヤー16, ディーラー10, ヒットしていない
	state2 := StrategyState{
		Player: StrategyHand{Sum: 16, HasAce: false},
		Dealer: StrategyHand{Sum: 10, HasAce: false},
		HasHit: false,
	}
	expectedPayouts2 := calc.CalculateAllExpectedPayouts(state2)
	fmt.Println("プレイヤーが11と5を持ち、ディーラーが10を持つ場合の、各アクションの期待払い戻し:")
	fmt.Printf("Hit: %f, Stand: %f, Surrender: %f, Best: %f\n", expectedPayouts2.HitPayout, expectedPayouts2.StandPayout, expectedPayouts2.SurrenderPayout, expectedPayouts2.BestPayout)
}

// ベンチマーク本体
func benchmarkCalculateAllExpectedPayouts(b *testing.B, calc interface {
	CalculateAllExpectedPayouts(StrategyState) StrategyExpectedPayouts
}) {
	states := []StrategyState{
		// 適当な状態を3つ用意
		{Player: StrategyHand{Sum: 2, HasAce: true}, Dealer: StrategyHand{Sum: 1, HasAce: true}, HasHit: false},
		{Player: StrategyHand{Sum: 18, HasAce: false}, Dealer: StrategyHand{Sum: 10, HasAce: false}, HasHit: false},
		{Player: StrategyHand{Sum: 15, HasAce: false}, Dealer: StrategyHand{Sum: 9, HasAce: false}, HasHit: true},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, state := range states {
			_ = calc.CalculateAllExpectedPayouts(state)
		}
	}
}

func BenchmarkCalculateAllExpectedPayouts_Calculator(b *testing.B) {
	benchmarkCalculateAllExpectedPayouts(b, NewCalculator())
}

func BenchmarkCalculateAllExpectedPayouts_UncachedCalculator(b *testing.B) {
	benchmarkCalculateAllExpectedPayouts(b, NewUncachedCalculator())
}
