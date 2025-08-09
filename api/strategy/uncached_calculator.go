package strategy

// メモ化なしのブラックジャック戦略計算用構造体
// あらゆる可能性を全探索し、期待値を足し合わせることで、最適解を求める。
type UncachedCalculator struct{}

func NewUncachedCalculator() *UncachedCalculator {
	return &UncachedCalculator{}
}

// ディーラーのアップカードから、ディーラーのスコア分布を計算する（メモ化なし）
func (c *UncachedCalculator) GetDealerScoreDistribution(dealerHand StrategyHand) map[int]float64 {
	currentScore := calculateScore(dealerHand)

	if currentScore >= 17 {
		return map[int]float64{currentScore: 1.0}
	}
	if currentScore == 0 {
		return map[int]float64{0: 1.0}
	}

	result := make(map[int]float64)
	for card, prob := range cardProbabilities {
		nextSum := dealerHand.Sum + card
		nextHasAce := dealerHand.HasAce || (card == 1)
		nextHand := StrategyHand{Sum: nextSum, HasAce: nextHasAce}
		subDist := c.GetDealerScoreDistribution(nextHand)
		for score, subProb := range subDist {
			result[score] += prob * subProb
		}
	}
	return result
}

// プレイヤーのスコアとディーラーの手札から、スタンド時の期待払い戻し（payout）を計算する（メモ化なし）
func (c *UncachedCalculator) CalculateStandExpectedPayout(playerScore int, dealerHand StrategyHand) float64 {
	if playerScore == 0 {
		return 0.0
	}
	dealerDist := c.GetDealerScoreDistribution(dealerHand)
	expectedPayout := 0.0
	for dealerScore, prob := range dealerDist {
		if dealerScore < playerScore {
			expectedPayout += 2.0 * prob
		} else if dealerScore == playerScore {
			expectedPayout += 1.0 * prob
		}
	}
	return expectedPayout
}

// ゲームの状態から、各アクションの期待払い戻しを計算する（メモ化なし）
func (c *UncachedCalculator) CalculateAllExpectedPayouts(state StrategyState) StrategyExpectedPayouts {
	var expectedPayouts StrategyExpectedPayouts
	playerScore := calculateScore(state.Player)

	if playerScore == 0 {
		expectedPayouts.StandPayout = 0
		expectedPayouts.HitPayout = 0
		if !state.HasHit {
			expectedPayouts.SurrenderPayout = 0.5
			expectedPayouts.BestPayout = 0.5
		} else {
			expectedPayouts.SurrenderPayout = 0.0
			expectedPayouts.BestPayout = 0.0
		}
		return expectedPayouts
	}

	standPayout := c.CalculateStandExpectedPayout(playerScore, state.Dealer)
	expectedPayouts.StandPayout = standPayout

	if !state.HasHit {
		expectedPayouts.SurrenderPayout = 0.5
	} else {
		expectedPayouts.SurrenderPayout = 0.0
	}

	hit := 0.0
	for card, prob := range cardProbabilities {
		nextSum := state.Player.Sum + card
		nextHasAce := state.Player.HasAce || (card == 1)
		nextHand := StrategyHand{Sum: nextSum, HasAce: nextHasAce}
		nextState := StrategyState{Player: nextHand, Dealer: state.Dealer, HasHit: true}
		hit += c.CalculateAllExpectedPayouts(nextState).BestPayout * prob
	}
	expectedPayouts.HitPayout = hit

	best := expectedPayouts.StandPayout
	if expectedPayouts.HitPayout > best {
		best = expectedPayouts.HitPayout
	}
	if expectedPayouts.SurrenderPayout > best {
		best = expectedPayouts.SurrenderPayout
	}
	expectedPayouts.BestPayout = best

	return expectedPayouts
}
