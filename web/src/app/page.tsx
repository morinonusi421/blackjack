'use client';

import { useState } from 'react';

// 乱数 API のレスポンス型
interface RandomNumberResponse {
  number: number;
}

export default function Home() {
  const [total, setTotal] = useState(0); // 表示する合計値
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
    <main style={{ fontFamily: 'sans-serif', textAlign: 'center', marginTop: '50px' }}>
      <h1>Blackjack Game</h1>

      <p>現在の合計:</p>
      <div
        style={{
          marginTop: '20px',
          padding: '20px',
          border: '1px solid #ccc',
          borderRadius: '8px',
          display: 'inline-block',
          minWidth: '120px',
        }}
      >
        <span style={{ fontSize: '2em', fontWeight: 'bold' }}>{total}</span>
      </div>

      {error && <p style={{ color: 'red' }}>エラー: {error}</p>}

      <button
        onClick={handleDraw}
        disabled={loading}
        style={{ marginTop: '20px', padding: '10px 20px', fontSize: '1em' }}
      >
        {loading ? '取得中...' : 'カードを引く'}
      </button>
    </main>
  );
}
