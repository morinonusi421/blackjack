export interface Card {
  suit: string;
  rank: string;
}

export interface Hand {
  cards: Card[];
  score: number;
}

export type GameState = 'PlayerTurn' | 'Finished';
export type Result = 'Pending' | 'PlayerWin' | 'DealerWin' | 'Push';

export interface Game {
  player_hand: Hand;
  dealer_hand: Hand;
  state: GameState;
  result: Result;
  result_message: string;
  bet: number;
  payout: number;
} 
