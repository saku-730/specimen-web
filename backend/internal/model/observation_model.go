// internal/model/observation_models.go
package model

import "time"

// ObservationMethod は "observation_methods" テーブルに対応するのだ
type ObservationMethod struct {
	ObservationMethodID uint   `gorm:"primaryKey" json:"observation_method_id"`
	MethodCommonName    string `json:"method_common_name"`
	PageID              uint   `gorm:"column:pageid" json:"page_id"` // SQLのカラム名が小文字なので合わせる

	// 関連
	WikiPage WikiPage `gorm:"foreignKey:PageID" json:"wiki_page"`
}

// Observation は "observations" テーブルに対応するのだ
type Observation struct {
	ObservationsID      uint      `gorm:"primaryKey" json:"observations_id"`
	UserID              uint      `json:"user_id"`
	OccurrenceID        uint      `json:"occurrence_id"`
	ObservationMethodID *uint     `json:"observation_method_id"`
	Behavior            string    `json:"behavior"`
	ObservedAt          time.Time `gorm:"default:now()" json:"observed_at"`
	Timezone            int16     `gorm:"not null" json:"timezone"`

	// 関連
	User              User              `gorm:"foreignKey:UserID" json:"user"`
	Occurrence        Occurrence        `gorm:"foreignKey:OccurrenceID" json:"occurrence"`
	ObservationMethod ObservationMethod `gorm:"foreignKey:ObservationMethodID" json:"observation_method"`
}
