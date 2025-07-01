'use client';

interface ErrorMessageProps {
  message: string;
}

export default function ErrorMessage({ message }: ErrorMessageProps) {
  if (!message) return null;
  return <p style={{ color: 'red' }}>エラー: {message}</p>;
} 
