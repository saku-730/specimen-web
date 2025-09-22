// internal/model/wiki_models.go
package model

import "time"

// WikiPage は "wiki_pages" テーブルに対応するのだ
type WikiPage struct {
	PageID      uint      `gorm:"primaryKey" json:"page_id"`
	Title       string    `json:"title"`
	UserID      uint      `json:"user_id"`
	CreatedDate time.Time `gorm:"default:now()" json:"created_date"`
	UpdatedDate *time.Time `json:"updated_date"`
	ContentPath string    `json:"content_path"`
}
