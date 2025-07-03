'use client';

import type { Game } from '../types/game';

interface GameInfoProps {
  game: Game | null;
}

/**
 * ゲームの状態を表示するコンポーネント。
 */
export default function GameInfo({ game }: GameInfoProps) {
  if (!game) return null;

  const formatCard = (card: { suit: string; rank: string }) => `${card.suit} ${card.rank}`;

  return (
    <section
      style={{
        display: 'flex',
        flexDirection: 'column',
        gap: '16px',
        width: '100%',
        maxWidth: '480px',
        border: '1px solid #ccc',
        borderRadius: '8px',
        padding: '16px',
      }}
    >
      <h2 style={{ margin: 0, textAlign: 'center' }}>ゲーム状況</h2>

      <div>
        <h3 style={{ marginBottom: '4px' }}>プレイヤー</h3>
        <p style={{ margin: 0 }}>手札: {game.player_hand.cards.map(formatCard).join(', ')}</p>
        <p style={{ margin: 0 }}>スコア: {game.player_hand.score}</p>
      </div>

      <div>
        <h3 style={{ marginBottom: '4px' }}>ディーラー</h3>
        <p style={{ margin: 0 }}>手札: {game.dealer_hand.cards.map(formatCard).join(', ')}</p>
        <p style={{ margin: 0 }}>スコア: {game.dealer_hand.score}</p>
      </div>

      <div>
        <p style={{ margin: 0 }}>ゲームの進行状況: {game.state}</p>
        <p style={{ margin: 0 }}>結果: {game.result}</p>
        {game.result_message && <p style={{ margin: 0 }}>メッセージ: {game.result_message}</p>}
        <p style={{ margin: 0 }}>掛け金: {game.bet}</p>
        <p style={{ margin: 0 }}>このゲームでのチップ増減: {game.balance_change}</p>
      </div>
    </section>
  );
} 
