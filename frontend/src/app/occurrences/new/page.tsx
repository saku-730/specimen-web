// src/app/occurrences/new/page.tsx
'use client';

import { useState, useEffect, FormEvent, ChangeEvent } from 'react';
import { useRouter } from 'next/navigation';

// --- 型定義: APIから取得する選択肢のデータの形 ---
interface SelectOption {
  id: number;
  name: string;
}

interface UserOption {
  user_id: number;
  user_name: string;
}

// --- フォーム全体のデータ状態を管理する型 ---
interface FormData {
  // === Occurrence ===
  project_id: string;
  user_id: string;
  individual_id: string;
  lifestage: string;
  sex: string;
  body_length: string;
  created_at: string;
  timezone: number;
  language_id: string;
  note: string;
  // === Classification ===
  species: string;
  genus: string;
  family: string;
  order: string;
  class: string;
  phylum: string;
  kingdom: string;
  // === Place ===
  latitude: string;
  longitude: string;
  place_name: string;
  // === Observation ===
  observation_user_id: string;
  observation_method_id: string;
  behavior: string;
  observed_at: string;
  observation_timezone: number;
  // === Specimen & MakeSpecimen ===
  make_specimen_user_id: string;
  specimen_method_id: string;
  make_specimen_created_at: string;
  make_specimen_timezone: number;
  institution_id: string;
  collection_id: string;
  // === Identification ===
  identification_user_id: string;
  source_info: string;
  identificated_at: string;
  identification_timezone: number;
}

// --- タイムスタンプを "YYYY-MM-DDTHH:MM" 形式で取得するヘルパー関数 ---
const getCurrentTimestamp = () => {
  const now = new Date();
  now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
  return now.toISOString().slice(0, 16);
};

// --- ここからがコンポーネント本体 ---
export default function NewOccurrencePage() {
  const router = useRouter();

  // フォーム全体の入力データを一つのstateで管理
  const [formData, setFormData] = useState<FormData>({
    project_id: '',
    user_id: '1', // 仮: ログインユーザーID
    individual_id: '',
    lifestage: '',
    sex: '',
    body_length: '',
    created_at: getCurrentTimestamp(),
    timezone: 9, // 初期値: JST
    language_id: '',
    note: '',
    species: '',
    genus: '',
    family: '',
    order: '',
    class: '',
    phylum: '',
    kingdom: '',
    latitude: '',
    longitude: '',
    place_name: '',
    observation_user_id: '1', // 仮: ログインユーザーID
    observation_method_id: '',
    behavior: '',
    observed_at: getCurrentTimestamp(),
    observation_timezone: 9,
    make_specimen_user_id: '1', // 仮: ログインユーザーID
    specimen_method_id: '',
    make_specimen_created_at: getCurrentTimestamp(),
    make_specimen_timezone: 9,
    institution_id: '',
    collection_id: '',
    identification_user_id: '1', // 仮: ログインユーザーID
    source_info: '',
    identificated_at: getCurrentTimestamp(),
    identification_timezone: 9,
  });

  // ドロップダウン用の選択肢リスト
  const [projects, setProjects] = useState<SelectOption[]>([]);
  const [users, setUsers] = useState<UserOption[]>([]);
  const [languages, setLanguages] = useState<SelectOption[]>([]);
  const [observationMethods, setObservationMethods] = useState<SelectOption[]>([]);
  const [specimenMethods, setSpecimenMethods] = useState<SelectOption[]>([]);
  const [institutionCodes, setInstitutionCodes] = useState<SelectOption[]>([]);
  const [collectionCodes, setCollectionCodes] = useState<SelectOption[]>([]);
  const [isLoadingOptions, setIsLoadingOptions] = useState(true);

  // フォームの入力値をまとめて更新するハンドラ
  const handleChange = (e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  // ページ読み込み時に、ドロップダウンの選択肢をAPIからまとめて取得
  useEffect(() => {
    const fetchOptions = async () => {
      try {
        const [
          projectsRes, usersRes, languagesRes, obsMethodsRes, 
          specMethodsRes, instCodesRes, collCodesRes
        ] = await Promise.all([
          fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/projects`),
          fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/users`),
          fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/languages`),
          fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/observation-methods`),
          fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/specimen-methods`),
          fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/institution-codes`),
          fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/collection-codes`),
        ]);

        // APIから取得したデータを整形してstateに保存
        const projectsData = await projectsRes.json();
        setProjects(projectsData.map((p: any) => ({ id: p.project_id, name: p.project_name })));

        setUsers(await usersRes.json());
        
        const languagesData = await languagesRes.json();
        setLanguages(languagesData.map((l: any) => ({ id: l.language_id, name: l.language_common })));
        
        const obsMethodsData = await obsMethodsRes.json();
        setObservationMethods(obsMethodsData.map((m: any) => ({ id: m.observation_method_id, name: m.method_common_name })));

        const specMethodsData = await specMethodsRes.json();
        setSpecimenMethods(specMethodsData.map((m: any) => ({ id: m.specimen_methods_id, name: m.method_common_name })));

        const instCodesData = await instCodesRes.json();
        setInstitutionCodes(instCodesData.map((i: any) => ({ id: i.institution_id, name: i.institution_code })));

        const collCodesData = await collCodesRes.json();
        setCollectionCodes(collCodesData.map((c: any) => ({ id: c.collection_id, name: c.collection_code })));

      } catch (error) {
        console.error("選択肢の取得に失敗しました:", error);
        alert("選択肢の取得に失敗しました。");
      } finally {
        setIsLoadingOptions(false);
      }
    };
    fetchOptions();
  }, []);

  // --- フォーム送信処理 ---
  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    
    // バックエンドに送信するデータ形式に整形
    const payload = {
      occurrence: {
        project_id: Number(formData.project_id),
        user_id: Number(formData.user_id),
        individual_id: formData.individual_id ? Number(formData.individual_id) : null,
        lifestage: formData.lifestage,
        sex: formData.sex,
        body_length: formData.body_length ? Number(formData.body_length) : null,
        created_at: new Date(formData.created_at).toISOString(),
        timezone: formData.timezone,
        language_id: Number(formData.language_id),
        note: formData.note,
      },
      classification: {
        class_classification: { // jsonb形式に整形
          species: formData.species, genus: formData.genus, family: formData.family,
          order: formData.order, class: formData.class, phylum: formData.phylum,
          kingdom: formData.kingdom,
        },
      },
      place: {
        coordinates: `POINT(${formData.longitude} ${formData.latitude})`, // PostGIS形式に変換
        place_name_json: { class_place_name: { name: formData.place_name } }, // jsonb形式に整形
      },
      observation: {
        user_id: Number(formData.observation_user_id),
        observation_method_id: Number(formData.observation_method_id),
        behavior: formData.behavior,
        observed_at: new Date(formData.observed_at).toISOString(),
        timezone: formData.observation_timezone,
      },
      specimen: {
        specimen_method_id: Number(formData.specimen_method_id),
        institution_id: Number(formData.institution_id),
        collection_id: Number(formData.collection_id),
      },
      make_specimen: {
        user_id: Number(formData.make_specimen_user_id),
        date: new Date(formData.make_specimen_created_at).toISOString().split('T')[0], // YYYY-MM-DD
        created_at: new Date(formData.make_specimen_created_at).toISOString(),
        timezone: formData.make_specimen_timezone,
      },
      identification: {
        user_id: Number(formData.identification_user_id),
        source_info: formData.source_info,
        identificated_at: new Date(formData.identificated_at).toISOString(),
        timezone: formData.identification_timezone,
      },
    };

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/full-occurrence`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (!response.ok) throw new Error('登録に失敗しました');
      
      alert('登録に成功しました！');
      router.push('/occurrences');
    } catch (error) {
      console.error(error);
      alert(`エラーが発生しました: ${error}`);
    }
  };

  if (isLoadingOptions) {
    return <div>フォームを準備中なのだ...</div>;
  }
  
  return (
    <form onSubmit={handleSubmit} className="space-y-6 bg-white p-8 rounded-lg shadow-md">
      <h1 className="text-2xl font-bold mb-6">新規データ 入力フォーム</h1>

      {/* --- 基本情報 --- */}
      <fieldset className="border p-4 rounded">
        <legend className="text-lg font-semibold px-2">基本情報</legend>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-2">
          <div>
            <label className="block text-sm font-medium text-gray-700">プロジェクト名*</label>
            <select name="project_id" value={formData.project_id} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm">
              <option value="">選択してください</option>
              {projects.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">データ入力ユーザ名*</label>
            <select name="user_id" value={formData.user_id} onChange={handleChange} required className="mt-1 block w-full rounded-md border-gray-300 shadow-sm">
              {users.map(u => <option key={u.user_id} value={u.user_id}>{u.user_name}</option>)}
            </select>
          </div>
          <div><label>個体ID</label><input type="number" name="individual_id" value={formData.individual_id} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
          <div><label>生育段階</label><input name="lifestage" value={formData.lifestage} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
          <div><label>性</label><input name="sex" value={formData.sex} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
          <div><label>体長</label><input type="number" step="0.01" name="body_length" value={formData.body_length} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
          <div><label>入力日時*</label><input type="datetime-local" name="created_at" value={formData.created_at} onChange={handleChange} required className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
          <div><label>タイムゾーン*</label><input type="number" name="timezone" value={formData.timezone} onChange={handleChange} required className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
          <div>
            <label>言語名</label>
            <select name="language_id" value={formData.language_id} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm">
                <option value="">選択してください</option>
                {languages.map(l => <option key={l.id} value={l.id}>{l.name}</option>)}
            </select>
          </div>
          <div className="md:col-span-3"><label>備考</label><textarea name="note" value={formData.note} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
        </div>
      </fieldset>

      {/* --- 分類 --- */}
      <fieldset className="border p-4 rounded">
        <legend className="text-lg font-semibold px-2">分類</legend>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-2">
            <div><label>界</label><input name="kingdom" value={formData.kingdom} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div><label>門</label><input name="phylum" value={formData.phylum} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div><label>鋼</label><input name="class" value={formData.class} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div><label>目</label><input name="order" value={formData.order} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div><label>科</label><input name="family" value={formData.family} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div><label>属</label><input name="genus" value={formData.genus} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div><label>種</label><input name="species" value={formData.species} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
        </div>
      </fieldset>

      {/* --- 場所 --- */}
      <fieldset className="border p-4 rounded">
        <legend className="text-lg font-semibold px-2">場所</legend>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-2">
            <div><label>緯度</label><input type="number" step="any" name="latitude" value={formData.latitude} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div><label>経度</label><input type="number" step="any" name="longitude" value={formData.longitude} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div className="md:col-span-2"><label>地名</label><textarea name="place_name" value={formData.place_name} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
        </div>
      </fieldset>
      
      {/* --- 観察 --- */}
      <fieldset className="border p-4 rounded">
        <legend className="text-lg font-semibold px-2">観察</legend>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-2">
            <div>
                <label>観察者*</label>
                <select name="observation_user_id" value={formData.observation_user_id} onChange={handleChange} required className="mt-1 block w-full rounded-md border-gray-300 shadow-sm">
                    {users.map(u => <option key={u.user_id} value={u.user_id}>{u.user_name}</option>)}
                </select>
            </div>
            <div>
                <label>観察・採集方法</label>
                <select name="observation_method_id" value={formData.observation_method_id} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm">
                    <option value="">選択してください</option>
                    {observationMethods.map(m => <option key={m.id} value={m.id}>{m.name}</option>)}
                </select>
            </div>
            <div className="md:col-span-2"><label>行動</label><input name="behavior" value={formData.behavior} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div><label>観察日時</label><input type="datetime-local" name="observed_at" value={formData.observed_at} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div><label>観察日時タイムゾーン</label><input type="number" name="observation_timezone" value={formData.observation_timezone} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
        </div>
      </fieldset>

      {/* --- 標本 --- */}
      <fieldset className="border p-4 rounded">
        <legend className="text-lg font-semibold px-2">標本</legend>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-2">
            <div>
                <label>標本作成者</label>
                <select name="make_specimen_user_id" value={formData.make_specimen_user_id} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm">
                    {users.map(u => <option key={u.user_id} value={u.user_id}>{u.user_name}</option>)}
                </select>
            </div>
            <div>
                <label>標本作成方法</label>
                <select name="specimen_method_id" value={formData.specimen_method_id} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm">
                    <option value="">選択してください</option>
                    {specimenMethods.map(m => <option key={m.id} value={m.id}>{m.name}</option>)}
                </select>
            </div>
            <div><label>標本作成日時</label><input type="datetime-local" name="make_specimen_created_at" value={formData.make_specimen_created_at} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div><label>標本作成日時タイムゾーン*</label><input type="number" name="make_specimen_timezone" value={formData.make_specimen_timezone} onChange={handleChange} required className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div>
                <label>機関コード</label>
                <select name="institution_id" value={formData.institution_id} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm">
                    <option value="">選択してください</option>
                    {institutionCodes.map(i => <option key={i.id} value={i.id}>{i.name}</option>)}
                </select>
            </div>
            <div>
                <label>コレクションコード</label>
                <select name="collection_id" value={formData.collection_id} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm">
                    <option value="">選択してください</option>
                    {collectionCodes.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
                </select>
            </div>
        </div>
      </fieldset>

      {/* --- 同定 --- */}
      <fieldset className="border p-4 rounded">
        <legend className="text-lg font-semibold px-2">同定</legend>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-2">
            <div>
                <label>同定者</label>
                <select name="identification_user_id" value={formData.identification_user_id} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm">
                    {users.map(u => <option key={u.user_id} value={u.user_id}>{u.user_name}</option>)}
                </select>
            </div>
            <div className="md:col-span-2"><label>参考情報</label><textarea name="source_info" value={formData.source_info} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div><label>同定日時</label><input type="datetime-local" name="identificated_at" value={formData.identificated_at} onChange={handleChange} className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
            <div><label>同定日時タイムゾーン*</label><input type="number" name="identification_timezone" value={formData.identification_timezone} onChange={handleChange} required className="mt-1 block w-full rounded-md border-gray-300 shadow-sm" /></div>
        </div>
      </fieldset>
      
      <div className="flex justify-end">
        <button type="submit" className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-lg">
          登録する
        </button>
      </div>
    </form>
  );
}
