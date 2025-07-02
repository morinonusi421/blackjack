package services

import (
	"testing"

	"blackjack/api/game"
)

func TestGameService_NewGame_InvalidBet(t *testing.T) {
	svc := NewGameService()
	_, err := svc.NewGame(0)
	if err == nil {
		t.Fatalf("expected error for non-positive bet, got nil")
	}
}

func TestGameService_NewGame_ValidBet(t *testing.T) {
	bet := 100
	svc := NewGameService()

	g, err := svc.NewGame(bet)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Bet が反映されている
	if g.Bet != bet {
		t.Fatalf("expected bet %d, got %d", bet, g.Bet)
	}

	// 手札枚数の検証
	if len(g.PlayerHand.Cards) != 2 {
		t.Fatalf("player should have 2 cards, got %d", len(g.PlayerHand.Cards))
	}
	if len(g.DealerHand.Cards) != 1 {
		t.Fatalf("dealer should have 1 card, got %d", len(g.DealerHand.Cards))
	}

	// スコア計算が一致しているか
	expectedPlayerScore := game.CalculateScore(g.PlayerHand.Cards)
	expectedDealerScore := game.CalculateScore(g.DealerHand.Cards)

	if g.PlayerHand.Score != expectedPlayerScore {
		t.Fatalf("player score mismatch: expected %d, got %d", expectedPlayerScore, g.PlayerHand.Score)
	}
	if g.DealerHand.Score != expectedDealerScore {
		t.Fatalf("dealer score mismatch: expected %d, got %d", expectedDealerScore, g.DealerHand.Score)
	}

	// 状態・結果の検証
	if expectedPlayerScore == 21 {
		if g.State != game.Finished || g.Result != game.PlayerWin {
			t.Fatalf("expected Finished/PlayerWin when player has blackjack, got state=%s result=%s", g.State, g.Result)
		}
		if g.BalanceChange != bet*3/2 {
			t.Fatalf("expected balance change %d, got %d", bet*3/2, g.BalanceChange)
		}
	} else {
		if g.State != game.PlayerTurn || g.Result != game.Pending {
			t.Fatalf("expected PlayerTurn/Pending for non-BJ, got state=%s result=%s", g.State, g.Result)
		}
		if g.BalanceChange != 0 {
			t.Fatalf("expected balance change 0, got %d", g.BalanceChange)
		}
	}
}
