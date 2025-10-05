// src/app/page.tsx

import Link from 'next/link';
import { FaChevronRight } from 'react-icons/fa'; 

export default function HomePage() {
  return (
    <main className="p-8">
      <h1 className="text-3xl font-bold text-gray-800 mb-2">
        標本管理ダッシュボード
      </h1>
      <p className="text-gray-600 mb-8">
      </p>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        
        <Link href="/occurrences/new" className="block p-6 bg-white rounded-lg border border-gray-200 shadow-md hover:shadow-lg hover:bg-gray-50 transition">
          <h2 className="text-xl font-semibold mb-2 text-gray-900">データ入力</h2>
          <p className="text-gray-700 mb-4">データを新規入力します。</p>
          <div className="flex items-center text-blue-600 font-medium">
            入力画面へ <FaChevronRight className="ml-2" />
          </div>
        </Link>

        <Link href="/projects" className="block p-6 bg-white rounded-lg border border-gray-200 shadow-md hover:shadow-lg hover:bg-gray-50 transition">
          <h2 className="text-xl font-semibold mb-2 text-gray-900">プロジェクト</h2>
          <p className="text-gray-700 mb-4">プロジェクトの管理や、メンバーの確認を行います。</p>
          <div className="flex items-center text-blue-600 font-medium">
            一覧へ <FaChevronRight className="ml-2" />
          </div>
        </Link>

        <Link href="/users" className="block p-6 bg-white rounded-lg border border-gray-200 shadow-md hover:shadow-lg hover:bg-gray-50 transition">
          <h2 className="text-xl font-semibold mb-2 text-gray-900">標本・採集手法</h2>
          <p className="text-gray-700 mb-4">標本作成方法、採集方法を確認、設定します</p>
          <div className="flex items-center text-blue-600 font-medium">
            一覧へ <FaChevronRight className="ml-2" />
          </div>
        </Link>
        
      </div>
    </main>
  );
}
