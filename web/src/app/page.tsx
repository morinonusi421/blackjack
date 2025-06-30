'use client';

import { useEffect, useState } from 'react';

interface ApiResponse {
  message: string;
}

export default function Home() {
  const [data, setData] = useState<ApiResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const apiUrl = process.env.NEXT_PUBLIC_API_URL;

  useEffect(() => {
    if (!apiUrl) {
      setError('APIのURLが設定されていません。');
      setLoading(false);
      return;
    }
    const fetchData = async () => {
      try {
        const res = await fetch(`${apiUrl}/api/hello`);
        if (!res.ok) throw new Error('APIからの応答がありませんでした');
        const result: ApiResponse = await res.json();
        setData(result);
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
    fetchData();
  }, [apiUrl]);

  return (
    <main style={{ fontFamily: 'sans-serif', textAlign: 'center', marginTop: '50px' }}>
      <h1>Blackjack Game</h1>
      <p>Go APIからのメッセージ:</p>
      <div style={{ marginTop: '20px', padding: '20px', border: '1px solid #ccc', borderRadius: '8px', display: 'inline-block' }}>
        {loading && <p>読み込み中...</p>}
        {error && <p style={{ color: 'red' }}>エラー: {error}</p>}
        {data && <p style={{ fontSize: '2em', fontWeight: 'bold' }}>{data.message}</p>}
      </div>
    </main>
  );
}
