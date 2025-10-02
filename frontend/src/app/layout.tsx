// src/app/layout.tsx
import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import Header from '@/components/Header'; // ① 作ったHeaderコンポーネントをインポート

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: '標本管理アプリ',
  description: '標本を管理するためのWebアプリケーション',
};

// これが全ページの共通レイアウトになるのだ
export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body className={inter.className}>
        {/* ② ここにヘッダーを配置！ */}
        <Header />

        {/* ③ メインコンテンツ部分 */}
        <main className="container mx-auto px-4 py-8">
          {children} {/* ← ここに各ページ(page.tsx)の中身が自動で入る */}
        </main>
      </body>
    </html>
  );
}
