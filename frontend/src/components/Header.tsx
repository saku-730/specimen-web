// src/components/Header.tsx
import Link from 'next/link';

// Headerコンポーネントの定義なのだ
export default function Header() {
  return (
    // headerタグで、ページのヘッダーであることを示すのだ
    <header className="bg-white shadow-md">
      <div className="container mx-auto px-4 py-4 flex gap-x-10 items-center">
        {/* 左側のロゴやタイトル */}
        <div className="text-xl text-gray-800 font-bold">
          <Link href="/" className="hover:text-gray-700">
            標本管理アプリ
          </Link>
        </div>

        {/* 右側のナビゲーションメニュー */}
        <nav>
          <ul className="flex space-x-6">
            <li>
              <Link href="/occurrences" className="text-gray-800 hover:text-blue-500">
                発生情報
              </Link>
            </li>
            <li>
              <Link href="/projects" className="text-gray-800 hover:text-blue-500">
                プロジェクト
              </Link>
            </li>
            <li>
              <Link href="/users" className="text-gray-800 hover:text-blue-500">
                ユーザー
              </Link>
            </li>
	    <li>
              <Link href="/login" className="text-gray-800 hover:text-blue-500">
                再ログイン
              </Link>
            </li>

          </ul>
        </nav>
      </div>
    </header>
  );
}
