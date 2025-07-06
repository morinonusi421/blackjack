'use client';

interface StartGameFormProps {
  bet: number;
  onBetChange: (val: number) => void;
  loading: boolean;
  onStart: () => void;
  disabled?: boolean;
}

export default function StartGameForm({ bet, onBetChange, loading, onStart, disabled = false }: StartGameFormProps) {
  return (
    <div
      style={{
        display: 'flex',
        gap: '12px',
        alignItems: 'center',
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
        type="number"
        min={1}
        value={bet}
        onChange={(e) => onBetChange(Number(e.target.value))}
        disabled={disabled || loading}
        style={{
          padding: '6px 10px',
          width: '120px',
          border: '1px solid #ccc',
          borderRadius: '4px',
        }}
      />
      <button
        onClick={onStart}
        disabled={loading || disabled}
        style={{
          padding: '8px 20px',
          backgroundColor: loading || disabled ? '#999' : '#0070f3',
          color: '#fff',
          border: 'none',
          borderRadius: '4px',
          cursor: loading || disabled ? 'not-allowed' : 'pointer',
          fontWeight: 'bold',
          boxShadow: '0 1px 3px rgba(0, 0, 0, 0.3)',
          opacity: loading || disabled ? 0.6 : 1,
          transition: 'opacity 0.2s',
        }}
      >
        {disabled ? '進行中' : loading ? '開始中...' : 'ゲーム開始'}
      </button>
    </div>
  );
} 
