// src/app/page.tsx

import Link from 'next/link';
import { FaChevronRight } from 'react-icons/fa'; // アイコン用のライブラリなのだ

/* * アイコンを使うには、ターミナルで下のコマンドを実行してね
 * npm install react-icons
 */

export default function HomePage() {
  return (
    <main className="p-8">
      {/* ページのメインタイトル */}
      <h1 className="text-3xl font-bold text-gray-800 mb-2">
        標本管理ダッシュボード
      </h1>
      <p className="text-gray-600 mb-8">
        ようこそ！ここから各機能にアクセスできます。
      </p>

      {/* 機能カードのグリッド表示 */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        
        {/* カード1: 発生情報 */}
        <Link href="/occurrences/new" className="block p-6 bg-white rounded-lg border border-gray-200 shadow-md hover:shadow-lg hover:bg-gray-50 transition">
          <h2 className="text-xl font-semibold mb-2 text-gray-900">データ登録</h2>
          <p className="text-gray-700 mb-4">データを新規登録します。</p>
          <div className="flex items-center text-blue-600 font-medium">
            入力画面へ <FaChevronRight className="ml-2" />
          </div>
        </Link>

        {/* カード2: プロジェクト */}
        <Link href="/projects" className="block p-6 bg-white rounded-lg border border-gray-200 shadow-md hover:shadow-lg hover:bg-gray-50 transition">
          <h2 className="text-xl font-semibold mb-2 text-gray-900">プロジェクト</h2>
          <p className="text-gray-700 mb-4">プロジェクトの管理や、メンバーの確認を行います。</p>
          <div className="flex items-center text-blue-600 font-medium">
            一覧へ <FaChevronRight className="ml-2" />
          </div>
        </Link>

        {/* カード3: ユーザー管理 */}
        <Link href="/users" className="block p-6 bg-white rounded-lg border border-gray-200 shadow-md hover:shadow-lg hover:bg-gray-50 transition">
          <h2 className="text-xl font-semibold mb-2 text-gray-900">ユーザー管理</h2>
          <p className="text-gray-700 mb-4">ユーザーの登録や、権限の変更を行います。</p>
          <div className="flex items-center text-blue-600 font-medium">
            一覧へ <FaChevronRight className="ml-2" />
          </div>
        </Link>
        
      </div>
    </main>
  );
}
