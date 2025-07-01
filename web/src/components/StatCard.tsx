'use client';

import React from 'react';

interface StatCardProps {
  label: string;
  value: number | null;
}

/**
 * ラベルと数値を表示するカードコンポーネント。
 */
export default function StatCard({ label, value }: StatCardProps) {
  return (
    <div
      style={{
        padding: '20px',
        border: '1px solid #ccc',
        borderRadius: '8px',
        minWidth: '120px',
      }}
    >
      <p style={{ margin: 0, fontSize: '0.9em', color: '#555' }}>{label}</p>
      <span style={{ fontSize: '2.4em', fontWeight: 'bold' }}>{value ?? 0}</span>
    </div>
  );
} 
