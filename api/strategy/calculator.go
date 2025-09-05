package strategy

import (
	"sync"

	"blackjack/api/game"
)

// Calculator はブラックジャックの最適戦略を計算する構造体
// メモ化テーブルを内部に保持し、スレッドセーフな計算を提供する
// 完全にステートレスで設定は各メソッドのパラメータとして受け取る
type Calculator struct {
	dealerMemo             map[dealerMemoKey]map[int]float64
	standPayoutMemo        map[standPayoutKey]float64
	allExpectedPayoutsMemo map[strategyStateKey]StrategyExpectedPayouts
	mu                     sync.RWMutex
}

type dealerMemoKey struct {
	hand   StrategyHand
	config game.GameConfig
}

type standPayoutKey struct {
	playerScore int
	dealerHand  StrategyHand
	config      game.GameConfig
}

type strategyStateKey struct {
	state  StrategyState
	config game.GameConfig
}

// NewCalculator は新しいCalculatorのインスタンスを生成
func NewCalculator() *Calculator {
	return &Calculator{
		dealerMemo:             make(map[dealerMemoKey]map[int]float64),
		standPayoutMemo:        make(map[standPayoutKey]float64),
		allExpectedPayoutsMemo: make(map[strategyStateKey]StrategyExpectedPayouts),
	}
}

// ディーラーのアップカードから、ディラーのスコア分布を計算する
func (c *Calculator) GetDealerScoreDistribution(dealerHand StrategyHand, config *game.GameConfig) map[int]float64 {
	key := dealerMemoKey{hand: dealerHand, config: *config}

	c.mu.RLock()
	if dist, found := c.dealerMemo[key]; found {
		c.mu.RUnlock()
		return dist
	}
	c.mu.RUnlock()

	currentScore := calculateScore(dealerHand)

	if currentScore >= config.DealerStandThreshold {
		result := map[int]float64{currentScore: 1.0}
		c.mu.Lock()
		c.dealerMemo[key] = result
		c.mu.Unlock()
		return result
	}
	if currentScore == 0 {
		result := map[int]float64{0: 1.0}
		c.mu.Lock()
		c.dealerMemo[key] = result
		c.mu.Unlock()
		return result
	}

	result := make(map[int]float64)
	for card, prob := range cardProbabilities {
		nextSum := dealerHand.Sum + card
		nextHasAce := dealerHand.HasAce || (card == 1)
		nextHand := StrategyHand{Sum: nextSum, HasAce: nextHasAce}
		subDist := c.GetDealerScoreDistribution(nextHand, config)
		for score, subProb := range subDist {
			result[score] += prob * subProb
		}
	}
	c.mu.Lock()
	c.dealerMemo[key] = result
	c.mu.Unlock()
	return result
}

// プレイヤーのスコアとディーラーの手札から、スタンド時の期待払い戻し（payout）を計算する
func (c *Calculator) CalculateStandExpectedPayout(playerScore int, dealerHand StrategyHand, config *game.GameConfig) float64 {

	// キャッシュがあれば、計算せずにそれを再利用
	key := standPayoutKey{playerScore: playerScore, dealerHand: dealerHand, config: *config}
	c.mu.RLock()
	if v, ok := c.standPayoutMemo[key]; ok {
		c.mu.RUnlock()
		return v
	}
	c.mu.RUnlock()

	// バーストしている場合は0.0
	if playerScore == 0 {
		return 0.0
	}

	// ディーラーのスコア分布を計算
	dealerDist := c.GetDealerScoreDistribution(dealerHand, config)

	// ディーラーのスコア分布とプレイヤーのスコアを比較して、期待払い戻しを計算
	expectedPayout := 0.0
	for dealerScore, prob := range dealerDist {
		if dealerScore < playerScore {
			expectedPayout += 2.0 * prob // 勝ち
		} else if dealerScore == playerScore {
			expectedPayout += 1.0 * prob // 引き分け
		} // 負けは0.0
	}

	// キャッシュに結果を保存
	c.mu.Lock()
	c.standPayoutMemo[key] = expectedPayout
	c.mu.Unlock()
	return expectedPayout
}

// ゲームの状態から、各アクションの期待払い戻しを計算する
func (c *Calculator) CalculateAllExpectedPayouts(state StrategyState, config *game.GameConfig) StrategyExpectedPayouts {

	// キャッシュがあれば、計算せずにそれを再利用
	key := strategyStateKey{state: state, config: *config}
	c.mu.RLock()
	if v, ok := c.allExpectedPayoutsMemo[key]; ok {
		c.mu.RUnlock()
		return v
	}
	c.mu.RUnlock()

	var expectedPayouts StrategyExpectedPayouts

	// プレイヤーのスコアを計算
	playerScore := calculateScore(state.Player)

	// プレイヤーがバーストしている場合
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
		c.mu.Lock()
		c.allExpectedPayoutsMemo[key] = expectedPayouts
		c.mu.Unlock()
		return expectedPayouts
	}

	// スタンドの期待値を計算
	standPayout := c.CalculateStandExpectedPayout(playerScore, state.Dealer, config)
	expectedPayouts.StandPayout = standPayout

	// サレンダーの期待値を計算 (ヒットしている場合はサレンダーができないので、期待利得を0としておく)
	if !state.HasHit {
		expectedPayouts.SurrenderPayout = 0.5
	} else {
		expectedPayouts.SurrenderPayout = 0.0
	}

	// ヒットの期待値を計算
	hit := 0.0
	for card, prob := range cardProbabilities {
		nextSum := state.Player.Sum + card
		nextHasAce := state.Player.HasAce || (card == 1)
		nextHand := StrategyHand{Sum: nextSum, HasAce: nextHasAce}
		nextState := StrategyState{Player: nextHand, Dealer: state.Dealer, HasHit: true}
		hit += c.CalculateAllExpectedPayouts(nextState, config).BestPayout * prob
	}
	expectedPayouts.HitPayout = hit

	// 最適行動を求める
	best := expectedPayouts.StandPayout
	if expectedPayouts.HitPayout > best {
		best = expectedPayouts.HitPayout
	}
	if expectedPayouts.SurrenderPayout > best {
		best = expectedPayouts.SurrenderPayout
	}
	expectedPayouts.BestPayout = best

	// キャッシュに結果を保存
	c.mu.Lock()
	c.allExpectedPayoutsMemo[key] = expectedPayouts
	c.mu.Unlock()
	return expectedPayouts
}
