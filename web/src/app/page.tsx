'use client';

import { useState } from 'react';

// 乱数 API のレスポンス型
interface RandomNumberResponse {
  number: number;
}

export default function Home() {
  const [total, setTotal] = useState(0); // 表示する合計値
  const [lastCard, setLastCard] = useState<number | null>(null); // 直近で引いたカード
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const apiUrl = process.env.NEXT_PUBLIC_API_URL;

  const handleDraw = async () => {
    if (!apiUrl) {
      setError('APIのURLが設定されていません。');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const res = await fetch(`${apiUrl}/api/random_number`);
      if (!res.ok) throw new Error('APIからの応答がありませんでした');
      const result: RandomNumberResponse = await res.json();
      setLastCard(result.number);
      setTotal((prev) => prev + result.number);
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('不明なエラーが発生しました');
      }
    } finally {
      setLoading(false);
    }
  };

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
        {/* 合計 */}
        <div
          style={{
            padding: '20px',
            border: '1px solid #ccc',
            borderRadius: '8px',
            minWidth: '120px',
          }}
        >
          <p style={{ margin: 0, fontSize: '0.9em', color: '#555' }}>合計</p>
          <span style={{ fontSize: '2.4em', fontWeight: 'bold' }}>{total}</span>
        </div>

        {/* 今回のカード */}
        <div
          style={{
            padding: '20px',
            border: '1px solid #ccc',
            borderRadius: '8px',
            minWidth: '120px',
          }}
        >
          <p style={{ margin: 0, fontSize: '0.9em', color: '#555' }}>今回のカード</p>
          <span style={{ fontSize: '2.4em', fontWeight: 'bold' }}>{lastCard ?? 0}</span>
        </div>
      </div>

      {error && <p style={{ color: 'red' }}>エラー: {error}</p>}

      <button
        onClick={handleDraw}
        disabled={loading}
        style={{ padding: '12px 32px', fontSize: '1.1em', borderRadius: '4px', cursor: 'pointer' }}
      >
        {loading ? '取得中...' : 'カードを引く'}
      </button>
    </main>
  );
}
