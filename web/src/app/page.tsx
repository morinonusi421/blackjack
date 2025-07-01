'use client';

import useBlackjack from '../hooks/useBlackjack';

import StatCard from '../components/StatCard';
import ErrorMessage from '../components/ErrorMessage';

export default function Home() {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;

  const { total, lastCard, loading, error, drawCard } = useBlackjack(apiUrl);

  return (
    <main
      style={{
        fontFamily: 'sans-serif',
        textAlign: 'center',
        marginTop: '40px',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        gap: '24px',
      }}
    >
      <h1>Blackjack Game</h1>

      <div
        style={{
          display: 'flex',
          gap: '40px',
          justifyContent: 'center',
          flexWrap: 'wrap',
        }}
      >
        <StatCard label="合計" value={total} />
        <StatCard label="今回のカード" value={lastCard} />
      </div>

      <ErrorMessage message={error ?? ''} />

      <button
        onClick={drawCard}
        disabled={loading}
        style={{ padding: '12px 32px', fontSize: '1.1em', borderRadius: '4px', cursor: 'pointer' }}
      >
        {loading ? '取得中...' : 'カードを引く'}
      </button>
    </main>
  );
}
