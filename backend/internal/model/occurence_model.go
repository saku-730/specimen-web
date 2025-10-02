// internal/model/occurrence_models.go
package model

import (
	"time"
	"gorm.io/datatypes" // JSON型のためにインポートするのだ
)

// Language は "language" テーブルに対応するのだ
type Language struct {
	LanguageID      uint   `gorm:"primaryKey" json:"language_id"`
	LanguageShort   string `json:"language_short"`
	LanguageCommon  string `json:"language_common"`
}

// ClassificationJSON は "classification_json" テーブルに対応するのだ
type ClassificationJSON struct {
	ClassificationID   uint           `gorm:"primaryKey" json:"classification_id"`
	ClassClassification datatypes.JSON `gorm:"type:jsonb" json:"class_classification"`
}

func (ClassificationJSON) TableName() string {
	return "classification_json"
}

// Occurrence は "occurrence" テーブルに対応するのだ
type Occurrence struct {
	OccurrenceID      uint      `gorm:"primaryKey" json:"occurrence_id"`
	ProjectID         uint      `json:"project_id"`
	UserID            uint      `json:"user_id"`
	IndividualID      *int      `json:"individual_id"`
	Lifestage         string    `json:"lifestage"`
	Sex               string    `json:"sex"`
	ClassificationID  uint      `json:"classification_id"`
	PlaceID           uint      `json:"place_id"`
	AttachmentGroupID *int      `json:"attachment_group_id"`
	BodyLength        float64   `gorm:"type:numeric" json:"body_length"`
	LanguageID        *uint     `json:"language_id"`
	Note              string    `json:"note"`
	CreatedAt         time.Time `gorm:"default:now()" json:"created_at"`
	Timezone          int16     `gorm:"not null" json:"timezone"`

	// 関連
	Attachments []Attachment `gorm:"many2many:attachment_group;" json:"attachments"` // 多対多
}

func (Occurrence) TableName() string {
	return "occurrence"
}
