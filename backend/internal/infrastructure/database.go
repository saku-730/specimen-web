// internal/infrastructure/database/database.go
package database

import (
	"fmt" 
	"log" 
	"github.com/saku-730/specimen-web/backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"          
)

//connect Database
func NewDatabaseConnection(cfg *configs.Config) (*gorm.DB, error) {
	dsn := cfg.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Printf("データベース接続に失敗しました: %v\n", err)
		return nil, fmt.Errorf("データベース接続エラー: %w", err)
	}

	log.Println("データベースへの接続に成功しました")

	return db, nil
}
