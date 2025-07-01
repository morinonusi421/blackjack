import BlackjackClient from './BlackjackClient';
import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'ブラックジャック'
};

export default function Home() {
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
      <BlackjackClient />
    </main>
  );
}
