'use client';

import type { CSSProperties } from 'react';

interface ActionButtonsProps {
  onStand: () => void;
  onHit: () => void;
  onSurrender: () => void;
  canSurrender: boolean;
  disabled?: boolean;
}

export default function ActionButtons({ onStand, onHit, onSurrender, canSurrender, disabled = false }: ActionButtonsProps) {
  // サレンダーは初回アクションでのみ可能
  const surrenderDisabled = disabled || !canSurrender;

  const baseStyle: CSSProperties = {
    padding: '10px 20px',
    backgroundColor: disabled ? '#999' : '#0070f3',
    color: '#fff',
    border: 'none',
    borderRadius: '4px',
    fontSize: '16px',
    cursor: disabled ? 'not-allowed' : 'pointer',
    opacity: disabled ? 0.6 : 1,
    transition: 'filter 0.2s',
  };

  const surrenderStyle: CSSProperties = {
    ...baseStyle,
    backgroundColor: surrenderDisabled ? '#999' : '#ff4444',
    cursor: surrenderDisabled ? 'not-allowed' : 'pointer',
    opacity: surrenderDisabled ? 0.6 : 1,
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
        Hit
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
      <button
        onClick={onSurrender}
        disabled={surrenderDisabled}
        style={surrenderStyle}
        onMouseOver={(e) => {
          if (!surrenderDisabled) (e.currentTarget.style.filter = 'brightness(1.1)');
        }}
        onMouseOut={(e) => {
          if (!surrenderDisabled) (e.currentTarget.style.filter = 'brightness(1)');
        }}
      >
        Surrender
      </button>
    </div>
  );
}
