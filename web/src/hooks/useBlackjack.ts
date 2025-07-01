'use client';

import { useState, useCallback } from 'react';

interface RandomNumberResponse {
  number: number;
}

/**
 * ブラックジャックのカード取得と状態管理を行うカスタムフック。
 *
 * @param apiUrl エンドポイントのベース URL（例: https://example.com）
 */
export default function useBlackjack(apiUrl?: string | null) {
  const [total, setTotal] = useState(0);
  const [lastCard, setLastCard] = useState<number | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const drawCard = useCallback(async () => {
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
  }, [apiUrl]);

  return {
    total,
    lastCard,
    loading,
    error,
    drawCard,
  } as const;
} 
