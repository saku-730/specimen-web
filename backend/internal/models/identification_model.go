// internal/model/identification_models.go
package model

import "time"

// Identification は "identifications" テーブルに対応するのだ
type Identification struct {
	IdentificationID uint      `gorm:"primaryKey" json:"identification_id"`
	UserID           uint      `json:"user_id"`
	OccurrenceID     uint      `json:"occurrence_id"`
	SourceInfo       string    `json:"source_info"`
	IdentificatedAt  time.Time `gorm:"default:now()" json:"identificated_at"`
	Timezone         int16     `gorm:"not null" json:"timezone"`

	// 関連
	User         User         `gorm:"foreignKey:UserID" json:"user"`
	Occurrence   Occurrence   `gorm:"foreignKey:OccurrenceID" json:"occurrence"`
}
