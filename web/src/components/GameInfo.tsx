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

  const getCardImagePath = (card: { suit: string; rank: string }) => {
    const suit = card.suit.toLowerCase();
    const rank = card.rank; // "A", "2".."10", "J", "Q", "K"
    return `/cards/${suit}_${rank}.png`;
  };

  const HandView = ({
    title,
    cards,
    score,
  }: {
    title: string;
    cards: { suit: string; rank: string }[];
    score: number;
  }) => (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '8px' }}>
        <h3 style={{ margin: 0 }}>{title}</h3>
        <span
          aria-label={`${title} のスコア`}
          style={{
            display: 'inline-block',
            padding: '2px 8px',
            fontSize: '12px',
            fontWeight: 700,
            color: '#1e3a8a',
            background: '#dbeafe',
            border: '1px solid #bfdbfe',
            borderRadius: '999px',
          }}
        >
          スコア {score}
        </span>
      </div>
      <div
        style={{
          display: 'flex',
          gap: '8px',
          flexWrap: 'wrap',
          alignItems: 'center',
        }}
      >
        {cards.map((c, idx) => (
          <img
            key={`${c.suit}-${c.rank}-${idx}`}
            src={getCardImagePath(c)}
            alt={`${c.suit} ${c.rank}`}
            style={{
              width: '72px',
              height: 'auto',
              borderRadius: '6px',
              boxShadow: '0 1px 3px rgba(0,0,0,0.15)',
              background: '#fff',
            }}
          />
        ))}
      </div>
    </div>
  );

  const StatCard = ({ label, value, emphasis = false }: { label: string; value: number | string; emphasis?: boolean }) => (
    <div
      style={{
        flex: 1,
        minWidth: '140px',
        border: '1px solid #e5e7eb',
        borderRadius: '8px',
        padding: '12px',
        background: emphasis ? '#ecfdf5' : '#ffffff',
        boxShadow: '0 1px 2px rgba(0,0,0,0.04)',
      }}
    >
      <div style={{ fontSize: '12px', color: '#6b7280', marginBottom: '6px' }}>{label}</div>
      <div style={{ fontSize: '20px', fontWeight: 700 }}>{value}</div>
    </div>
  );

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

      <HandView title="プレイヤー" cards={game.player_hand.cards} score={game.player_hand.score} />

      <HandView title="ディーラー" cards={game.dealer_hand.cards} score={game.dealer_hand.score} />

      {game.result_message && (
        <div
          style={{
            padding: '10px 12px',
            border: '1px solid #fecaca',
            background: '#fee2e2',
            color: '#7f1d1d',
            borderRadius: '8px',
          }}
        >
          {game.result_message}
        </div>
      )}

      <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap' }}>
        <StatCard label="掛け金" value={game.bet} />
        <StatCard 
          label="払戻金" 
          value={game.state === 'Finished' ? game.payout : '?'} 
          emphasis={game.state === 'Finished' && game.payout > 0} 
        />
      </div>
    </section>
  );
} 
