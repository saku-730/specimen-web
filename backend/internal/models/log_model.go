// internal/model/log_models.go
package model

import "time"

// ChangeLog は "change_logs" テーブルに対応するのだ
type ChangeLog struct {
	LogID       uint      `gorm:"primaryKey" json:"log_id"`
	Type        string    `json:"type"`
	ChangedID   uint      `json:"changed_id"`
	BeforeValue string    `json:"before_value"`
	AfterValue  string    `json:"after_value"`
	UserID      uint      `json:"user_id"`
	Date        time.Time `gorm:"default:now()" json:"date"`
	Row         string    `gorm:"column:row" json:"row"` // SQLのRowは予約語の可能性があるのでカラム名を指定するのだ

	// 関連
	User User `gorm:"foreignKey:UserID" json:"user"`
}
