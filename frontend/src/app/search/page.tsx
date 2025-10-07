// src/app/search/page.tsx
'use client';

import { useState, useEffect, FormEvent, ChangeEvent } from 'react';

// --- 型定義 ---
interface SelectOption { id: number; name: string; }
interface UserOption { user_id: number; user_name: string; }

// 検索フォームの状態を管理する型
interface SearchParams {
  identifierId: string;
  collectorId: string;
  dataEntryUserId: string;
  specimenMakerId: string;
  species: string;
  genus: string;
  family: string;
  order: string;
  class: string;
  phylum: string;
  kingdom: string;
  projectId: string;
  institutionId: string;
  collectionId: string;
  latitude: string;
  longitude: string;
  placeName: string;
  occurrenceDateStart: string;
  occurrenceDateEnd: string;
  specimenDateStart: string;
  specimenDateEnd: string;
  observationMethodId: string;
  specimenMethodId: string;
  lifestage: string;
  sex: string;
  note: string;
  behavior: string;
  sourceInfo: string;
}

// 検索結果の型（バックエンドのAPIレスポンスに合わせる）
interface SearchResult {
  occurrence_id: number;
  note: string;
  user_name: string;
  project_name: string;
  species: string;
}

const initialSearchParams: SearchParams = {
  identifierId: '', collectorId: '', dataEntryUserId: '', specimenMakerId: '',
  species: '', genus: '', family: '', order: '', class: '', phylum: '', kingdom: '',
  projectId: '', institutionId: '', collectionId: '',
  latitude: '', longitude: '', placeName: '',
  occurrenceDateStart: '', occurrenceDateEnd: '',
  specimenDateStart: '', specimenDateEnd: '',
  observationMethodId: '', specimenMethodId: '',
  lifestage: '', sex: '', note: '', behavior: '', sourceInfo: '',
};

// --- コンポーネント本体 ---
export default function SearchPage() {
  const [searchParams, setSearchParams] = useState<SearchParams>(initialSearchParams);
  const [results, setResults] = useState<SearchResult[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  // ドロップダウン用の選択肢リスト
  const [users, setUsers] = useState<UserOption[]>([]);
  const [projects, setProjects] = useState<SelectOption[]>([]);
  const [institutionCodes, setInstitutionCodes] = useState<SelectOption[]>([]);
  const [collectionCodes, setCollectionCodes] = useState<SelectOption[]>([]);
  const [observationMethods, setObservationMethods] = useState<SelectOption[]>([]);
  const [specimenMethods, setSpecimenMethods] = useState<SelectOption[]>([]);

  // ページ読み込み時に選択肢をAPIから取得
  useEffect(() => {
    const fetchOptions = async () => {
      try {
        const [usersRes, projectsRes, instRes, collRes, obsRes, specRes] = await Promise.all([
          fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/users`),
          fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/projects`),
          fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/institution-codes`),
          fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/collection-codes`),
          fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/observation-methods`),
          fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/specimen-methods`),
        ]);
        setUsers(await usersRes.json());
        const projectsData = await projectsRes.json();
        setProjects(projectsData.map((p: any) => ({ id: p.project_id, name: p.project_name })));
        // ... 他の選択肢も同様に .map で整形 ...
      } catch (error) {
        console.error("選択肢の取得に失敗しました:", error);
      }
    };
    fetchOptions();
  }, []);

  // フォーム入力ハンドラ
  const handleChange = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    setSearchParams(prev => ({ ...prev, [e.target.name]: e.target.value }));
  };
  
  // フォームクリアハンドラ
  const handleClear = () => {
    setSearchParams(initialSearchParams);
    setResults([]);
  };

  // 検索実行ハンドラ
  const handleSearch = async (e: FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setResults([]);

    // 空でない条件だけをURLクエリパラメータとして組み立てる
    const query = new URLSearchParams();
    Object.entries(searchParams).forEach(([key, value]) => {
      if (value) { // valueが空文字やnullでなければ追加
        query.append(key, value);
      }
    });

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/search?${query.toString()}`);
      if (!response.ok) throw new Error('検索に失敗しました');
      const data = await response.json();
      setResults(data);
    } catch (error) {
      console.error(error);
      alert('検索中にエラーが発生しました');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container mx-auto p-8">
      <h1 className="text-3xl font-bold mb-6">詳細検索</h1>

      <form onSubmit={handleSearch} className="bg-white p-6 rounded-lg shadow-md space-y-6">
        {/* --- 人物関連 --- */}
        <fieldset className="border p-4 rounded">
          <legend className="text-lg font-semibold px-2">人物</legend>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-2">
            <div><label>同定者</label><select name="identifierId" value={searchParams.identifierId} onChange={handleChange} className="w-full mt-1"><option value="">すべて</option>{users.map(u => <option key={u.user_id} value={u.user_id}>{u.user_name}</option>)}</select></div>
            <div><label>採集者</label><select name="collectorId" value={searchParams.collectorId} onChange={handleChange} className="w-full mt-1"><option value="">すべて</option>{users.map(u => <option key={u.user_id} value={u.user_id}>{u.user_name}</option>)}</select></div>
            {/* ... データ入力担当、標本作成者も同様 ... */}
          </div>
        </fieldset>
        
        {/* --- 分類 --- */}
        <fieldset className="border p-4 rounded">
          <legend className="text-lg font-semibold px-2">分類</legend>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-2">
            <div><label>種</label><input name="species" value={searchParams.species} onChange={handleChange} className="w-full mt-1" /></div>
            <div><label>属</label><input name="genus" value={searchParams.genus} onChange={handleChange} className="w-full mt-1" /></div>
            {/* ... 他の分類レベルも同様 ... */}
          </div>
        </fieldset>

        {/* --- プロジェクト・標本の場所 --- */}
        <fieldset className="border p-4 rounded">
            {/* ... プロジェクト、機関コード、コレクションコードのドロップダウン ... */}
        </fieldset>

        {/* --- 採集場所・日時 --- */}
        <fieldset className="border p-4 rounded">
          <legend className="text-lg font-semibold px-2">採集</legend>
           <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-2">
            <div><label>緯度</label><input type="number" step="any" name="latitude" value={searchParams.latitude} onChange={handleChange} className="w-full mt-1" /></div>
            <div><label>経度</label><input type="number" step="any" name="longitude" value={searchParams.longitude} onChange={handleChange} className="w-full mt-1" /></div>
            <div className="col-span-2"><label>地名</label><input name="placeName" value={searchParams.placeName} onChange={handleChange} className="w-full mt-1" /></div>
            <div><label>採集/観察日 (開始)</label><input type="date" name="occurrenceDateStart" value={searchParams.occurrenceDateStart} onChange={handleChange} className="w-full mt-1" /></div>
            <div><label>採集/観察日 (終了)</label><input type="date" name="occurrenceDateEnd" value={searchParams.occurrenceDateEnd} onChange={handleChange} className="w-full mt-1" /></div>
            {/* ... 採集方法など ... */}
          </div>
        </fieldset>
        
        {/* ... 他のfieldsetも同様に作成 ... */}

        <div className="flex justify-end space-x-4">
          <button type="button" onClick={handleClear} className="bg-gray-500 hover:bg-gray-600 text-white font-bold py-2 px-6 rounded-lg">クリア</button>
          <button type="submit" className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-lg">検索</button>
        </div>
      </form>
      
      {/* --- 検索結果表示 --- */}
      <div className="mt-8">
        <h2 className="text-2xl font-bold mb-4">検索結果</h2>
        {isLoading ? (
          <p>検索中...</p>
        ) : (
          <table className="min-w-full bg-white shadow-md rounded-lg">
            {/* ... テーブルのヘッダー ... */}
            <tbody>
              {results.length > 0 ? (
                results.map(res => (
                  <tr key={res.occurrence_id}>
                    <td className="border p-2">{res.occurrence_id}</td>
                    <td className="border p-2">{res.user_name}</td>
                    <td className="border p-2">{res.project_name}</td>
                    <td className="border p-2">{res.species}</td>
                  </tr>
                ))
              ) : (
                <tr><td colSpan={4} className="text-center p-4">検索結果がありません</td></tr>
              )}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
