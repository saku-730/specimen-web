// src/app/users/page.tsx

// 1. 'use client' を宣言するのだ
'use client'; 

import { useState, useEffect } from 'react';

// 2. APIから返ってくるUserの「形」を定義する
interface User {
  user_id: number;
  user_name: string;
  display_name: string;
}

// 3. これがページの本体となるコンポーネントなのだ
export default function UsersPage() {
  // 4. stateを定義する
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 5. useEffectを使って、コンポーネントが最初に表示された時にデータを取得する
  useEffect(() => {
    // 6. データを取得するための非同期関数を定義
    const fetchUsers = async () => {
      try {
        // 7. 環境変数からAPIのURLを取得して、fetchでリクエストを送る
        const apiUrl = `${process.env.NEXT_PUBLIC_API_BASE_URL}/users`;
        const response = await fetch(apiUrl);

        // 8. レスポンスが成功でなければエラーを投げる
        if (!response.ok) {
          throw new Error('データの取得に失敗しました');
        }

        // 9. レスポンスのJSONをパースして、stateに保存する
        const data: User[] = await response.json();
        setUsers(data);
      } catch (err: any) {
        setError(err.message);
      } finally {
        // 10. 成功しても失敗しても、ローディング状態を解除する
        setLoading(false);
      }
    };

    fetchUsers(); // 11. 関数を実行！
  }, []); // 12. []が空なので、この処理は最初に1回だけ実行される

  // 13. ローディング中やエラー時の表示
  if (loading) return <div>読み込み中なのだ...</div>;
  if (error) return <div>エラー: {error}</div>;

  // 14. 取得したデータを表示する
  return (
    <div>
      <h1>ユーザー一覧</h1>
      <ul>
        {users.map((user) => (
          <li key={user.user_id}>
            {user.display_name} (@{user.user_name})
          </li>
        ))}
      </ul>
    </div>
  );
}
