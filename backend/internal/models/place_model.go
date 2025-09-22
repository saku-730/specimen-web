// internal/model/occurrence_models.go
package model

import (
	"time"
	"gorm.io/datatypes" // JSON型のためにインポートするのだ
)

// PlaceNameJSON は "place_names_json" テーブルに対応するのだ
type PlaceNameJSON struct {
	PlaceNameID      uint           `gorm:"primaryKey" json:"place_name_id"`
	ClassPlaceName   datatypes.JSON `gorm:"type:jsonb" json:"class_place_name"`
}

// Place は "places" テーブルに対応するのだ
type Place struct {
	PlaceID       uint    `gorm:"primaryKey" json:"place_id"`
	Coordinates   string  `gorm:"type:geography(Point,4326)" json:"coordinates"` // PostGIS型はstringで受けて、別途ライブラリで処理するのが一般的なのだ
	PlaceNameID   uint    `json:"place_name_id"`
	Accuracy      float64 `gorm:"type:numeric" json:"accuracy"`
}

