package strategy

var cardProbabilities = map[int]float64{
	1: 1.0 / 13.0, 2: 1.0 / 13.0, 3: 1.0 / 13.0, 4: 1.0 / 13.0,
	5: 1.0 / 13.0, 6: 1.0 / 13.0, 7: 1.0 / 13.0, 8: 1.0 / 13.0,
	9: 1.0 / 13.0, 10: 4.0 / 13.0,
}

// 戦略計算用の状態 game.goのGameStateより簡素
type StrategyState struct {
	Player StrategyHand
	Dealer StrategyHand
	HasHit bool
}

// 戦略計算用の手札 game.goのHandより簡素
type StrategyHand struct {
	Sum    int
	HasAce bool
}

// 各アクションの期待払い戻しをまとめて返す型
type StrategyExpectedPayouts struct {
	HitPayout       float64
	StandPayout     float64
	SurrenderPayout float64
	BestPayout      float64 //上記3つの中で最も高い期待値
}

// 新たな表現用のスコア計算関数
func calculateScore(sHand StrategyHand) int {
	// バースト
	if sHand.Sum > 21 {
		return 0
	}
	if sHand.HasAce && sHand.Sum+10 <= 21 {
		return sHand.Sum + 10
	}
	return sHand.Sum
}
