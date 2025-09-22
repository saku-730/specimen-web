// internal/infrastructure/database/database.go
package database

import (
	"fmt" // エラーメッセージを整形するためにインポートするのだ
	"log" 
	"github.com/saku-730/specimen-web/backend/config"
	"gorm.io/driver/postgres" // GORMのPostgreSQL用ドライバなのだ
	"gorm.io/gorm"            // GORM本体なのだ
)

// NewDatabaseConnectionは、渡された設定情報を使ってデータベースに接続し、
// その接続オブジェクト(*gorm.DB)を返す関数なのだ
func NewDatabaseConnection(cfg *configs.Config) (*gorm.DB, error) {
	// 1. configsパッケージに作ったDSN()メソッドを呼び出して、接続文字列を取得するのだ
	dsn := cfg.DSN()

	// 2. GORMを使ってPostgreSQLデータベースへの接続を試みるのだ
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// 3. もし接続でエラーが発生したら、ログにエラーを記録して、エラーを返すのだ
	if err != nil {
		log.Printf("データベース接続に失敗しました: %v\n", err)
		return nil, fmt.Errorf("データベース接続エラー: %w", err)
	}

	// 4. 接続に成功したことをログに表示するのだ
	log.Println("データベースへの接続に成功しました！")

	// 5. 成功したら、データベース操作ができるdbオブジェクトと、nil(エラーなし)を返すのだ
	return db, nil
}
