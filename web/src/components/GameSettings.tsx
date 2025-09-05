'use client';

import { DEFAULT_DEALER_THRESHOLD } from '../constants/config';

interface GameSettingsProps {
  threshold: number;
  onThresholdChange: (threshold: number) => void;
}

export default function GameSettings({ threshold, onThresholdChange }: GameSettingsProps) {
  const handleThresholdChange = (value: number) => {
    onThresholdChange(value);
  };

  return (
    <div
      style={{
        background: '#f8f9fa',
        border: '1px solid #e9ecef',
        borderRadius: '8px',
        padding: '20px',
        maxWidth: '400px',
        margin: '0 auto',
      }}
    >
      <h3
        style={{
          margin: '0 0 16px 0',
          fontSize: '18px',
          fontWeight: '600',
          color: '#333',
        }}
      >
        ゲーム設定
      </h3>

      <div>
        <label
          htmlFor="dealer-threshold"
          style={{
            display: 'block',
            marginBottom: '8px',
            fontSize: '14px',
            fontWeight: '500',
            color: '#555',
          }}
        >
          ディーラーのスタンド閾値: {threshold}
        </label>
        <input
          id="dealer-threshold"
          type="range"
          min="1"
          max="21"
          value={threshold}
          onChange={(e) => handleThresholdChange(parseInt(e.target.value))}
          style={{
            width: '100%',
            height: '6px',
            background: '#ddd',
            borderRadius: '3px',
            outline: 'none',
            cursor: 'pointer',
          }}
        />
        <div
          style={{
            display: 'flex',
            fontSize: '12px',
            color: '#777',
            marginTop: '4px',
            position: 'relative',
            height: '14px',
          }}
        >
          <span style={{ position: 'absolute', left: '0%', transform: 'translateX(-50%)', textAlign: 'center' }}>1</span>
          <span style={{ position: 'absolute', left: '45%', transform: 'translateX(-50%)', textAlign: 'center' }}>10</span>
          <span style={{ position: 'absolute', left: '80%', transform: 'translateX(-50%)', textAlign: 'center' }}>{DEFAULT_DEALER_THRESHOLD}</span>
          <span style={{ position: 'absolute', left: '100%', transform: 'translateX(-50%)', textAlign: 'center' }}>21</span>
        </div>
        <p
          style={{
            fontSize: '13px',
            color: '#666',
            margin: '8px 0 16px 0',
            lineHeight: '1.4',
          }}
        >
          ディーラーがこの値以上になるまでヒットし続けます。
          <br />
          通常のルールでは{DEFAULT_DEALER_THRESHOLD}です。
        </p>
        
        <button
          type="button"
          onClick={() => handleThresholdChange(DEFAULT_DEALER_THRESHOLD)}
          style={{
            padding: '6px 12px',
            fontSize: '12px',
            fontWeight: '500',
            color: '#007bff',
            background: 'transparent',
            border: '1px solid #007bff',
            borderRadius: '4px',
            cursor: 'pointer',
            transition: 'all 0.2s',
          }}
          onMouseEnter={(e) => {
            e.currentTarget.style.background = '#007bff';
            e.currentTarget.style.color = 'white';
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.background = 'transparent';
            e.currentTarget.style.color = '#007bff';
          }}
        >
          デフォルトに戻す
        </button>
      </div>
    </div>
  );
}
