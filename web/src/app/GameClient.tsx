'use client';

import { useState } from 'react';
import useGame from '../hooks/useGame';
import GameInfo from '../components/GameInfo';
import ErrorMessage from '../components/ErrorMessage';
import StartGameForm from '../components/StartGameForm';
import Balance from '../components/Balance';

// シンプルなボタンコンポーネント
function ActionButtons({
  onStand,
  onHit,
  disabled,
}: {
  onStand: () => void;
  onHit: () => void;
  disabled?: boolean;
}) {
  const baseStyle: React.CSSProperties = {
    padding: '10px 20px',
    backgroundColor: '#0070f3',
    color: '#fff',
    border: 'none',
    borderRadius: '4px',
    fontSize: '16px',
    cursor: disabled ? 'not-allowed' : 'pointer',
    opacity: disabled ? 0.6 : 1,
    transition: 'filter 0.2s',
  };

  return (
    <div style={{ display: 'flex', gap: '12px' }}>
      <button
        onClick={onHit}
        disabled={disabled}
        style={baseStyle}
        onMouseOver={(e) => {
          if (!disabled) (e.currentTarget.style.filter = 'brightness(1.1)');
        }}
        onMouseOut={(e) => {
          if (!disabled) (e.currentTarget.style.filter = 'brightness(1)');
        }}
      >
        Hit (未実装)
      </button>
      <button
        onClick={onStand}
        disabled={disabled}
        style={baseStyle}
        onMouseOver={(e) => {
          if (!disabled) (e.currentTarget.style.filter = 'brightness(1.1)');
        }}
        onMouseOut={(e) => {
          if (!disabled) (e.currentTarget.style.filter = 'brightness(1)');
        }}
      >
        Stand
      </button>
    </div>
  );
}

export default function GameClient() {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;
  const { game, loading, error, startGame, stand, balance } = useGame(apiUrl);
  const [bet, setBet] = useState(100);

  const handleStart = () => {
    startGame(bet);
  };

  const gameInProgress = game?.state === 'PlayerTurn';

  return (
    <section
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        gap: '24px',
      }}
    >
      {/* 所持金表示 */}
      <Balance balance={balance} />

      {/* 掛け金入力 & ゲーム開始ボタン */}
      <StartGameForm
        bet={bet}
        onBetChange={setBet}
        loading={loading}
        onStart={handleStart}
        disabled={gameInProgress}
      />

      {/* エラーメッセージ */}
      <ErrorMessage message={error ?? ''} />

      {/* ゲーム状況表示 */}
      <GameInfo game={game} />

      {/* Hit / Stand アクション */}
      {game && game.state === 'PlayerTurn' && (
        <ActionButtons
          onHit={() => alert('Hit は未実装です')}
          onStand={stand}
          disabled={loading}
        />
      )}
    </section>
  );
} 
