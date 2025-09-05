'use client';

import { useState, useEffect } from 'react';

interface StartGameFormProps {
  bet: number;
  onBetChange: (val: number) => void;
  loading: boolean;
  onStart: () => void;
  disabled?: boolean;
}

export default function StartGameForm({ bet, onBetChange, loading, onStart, disabled = false }: StartGameFormProps) {
  const [inputValue, setInputValue] = useState(bet.toString());

  // 外部からbetが変更された時は入力値を同期
  useEffect(() => {
    setInputValue(bet.toString());
  }, [bet]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    
    // 数値のみを許可（空文字列も許可）
    if (value === '' || /^\d+$/.test(value)) {
      setInputValue(value);
      
      // 空文字列でない場合のみ数値として親に渡す
      if (value !== '') {
        const numValue = Number(value);
        if (numValue >= 1) {
          onBetChange(numValue);
        }
      }
    }
  };
  return (
    <div
      className="start-form"
      style={{
        border: '1px solid #ccc',
        borderRadius: '8px',
        padding: '12px 16px',
      }}
    >
      <label htmlFor="bet" style={{ fontWeight: 'bold' }}>
        掛け金:
      </label>
      <input
        id="bet"
        type="text"
        value={inputValue}
        onChange={handleInputChange}
        disabled={disabled || loading}
        placeholder="掛け金を入力"
        style={{
          padding: '6px 10px',
          width: '120px',
          border: '1px solid #ccc',
          borderRadius: '4px',
        }}
      />
      <button
        onClick={onStart}
        disabled={loading || disabled || inputValue === '' || bet < 1}
        style={{
          padding: '8px 20px',
          backgroundColor: (loading || disabled || inputValue === '' || bet < 1) ? '#999' : '#0070f3',
          color: '#fff',
          border: 'none',
          borderRadius: '4px',
          cursor: (loading || disabled || inputValue === '' || bet < 1) ? 'not-allowed' : 'pointer',
          fontWeight: 'bold',
          boxShadow: '0 1px 3px rgba(0, 0, 0, 0.3)',
          opacity: (loading || disabled || inputValue === '' || bet < 1) ? 0.6 : 1,
          transition: 'opacity 0.2s',
        }}
      >
        {disabled ? '進行中' : loading ? '開始中...' : 'ゲーム開始'}
      </button>
    </div>
  );
} 
