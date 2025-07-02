package game

import (
	"testing"
)

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name     string
		cards    []Card
		expected int
	}{
		{
			name:     "simple numbers",
			cards:    []Card{{Suit: Spade, Rank: "2"}, {Suit: Heart, Rank: "3"}},
			expected: 5,
		},
		{
			name:     "face cards",
			cards:    []Card{{Suit: Spade, Rank: "J"}, {Suit: Diamond, Rank: "Q"}},
			expected: 20,
		},
		{
			name:     "ace as 11",
			cards:    []Card{{Suit: Club, Rank: "A"}, {Suit: Heart, Rank: "9"}},
			expected: 20,
		},
		{
			name:     "ace adjustment",
			cards:    []Card{{Suit: Club, Rank: "A"}, {Suit: Heart, Rank: "9"}, {Suit: Spade, Rank: "3"}},
			expected: 13,
		},
	}

	for _, tc := range tests {
		got := CalculateScore(tc.cards)
		if got != tc.expected {
			t.Fatalf("%s: expected %d, got %d", tc.name, tc.expected, got)
		}
	}
}
