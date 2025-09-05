package strategy

import (
	"testing"
	
	"blackjack/api/game"
)

func TestGetDealerScoreDistribution_Calculator(t *testing.T) {
	config := &game.GameConfig{DealerStandThreshold: 17}
	calc := NewCalculator()
	
	dist := calc.GetDealerScoreDistribution(StrategyHand{Sum: 2, HasAce: false}, config)
	if len(dist) == 0 {
		t.Fatalf("expected non-empty distribution")
	}
	
	dist2 := calc.GetDealerScoreDistribution(StrategyHand{Sum: 1, HasAce: true}, config)
	if len(dist2) == 0 {
		t.Fatalf("expected non-empty distribution for ace")
	}
}

func TestCalculateStandExpectedPayout_Calculator(t *testing.T) {
	config := &game.GameConfig{DealerStandThreshold: 17}
	calc := NewCalculator()
	dealerHand := StrategyHand{Sum: 1, HasAce: true}
	standPayout := calc.CalculateStandExpectedPayout(18, dealerHand, config)
	
	if standPayout < 0 || standPayout > 2 {
		t.Fatalf("expected payout in [0,2], got %f", standPayout)
	}
}

func TestCalculateAllExpectedPayouts_Calculator(t *testing.T) {
	config := &game.GameConfig{DealerStandThreshold: 17}
	calc := NewCalculator()
	
	state := StrategyState{
		Player: StrategyHand{Sum: 20, HasAce: false},
		Dealer: StrategyHand{Sum: 10, HasAce: false},
		HasHit: false,
	}
	
	payouts := calc.CalculateAllExpectedPayouts(state, config)
	if payouts.StandPayout < 0 {
		t.Fatalf("expected non-negative stand payout")
	}
	if payouts.HitPayout < 0 {
		t.Fatalf("expected non-negative hit payout")  
	}
	if payouts.SurrenderPayout != 0.5 {
		t.Fatalf("expected surrender payout 0.5, got %f", payouts.SurrenderPayout)
	}
}