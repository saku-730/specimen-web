// internal/model/specimen_models.go
package model

import "time"

// InstitutionIDCode は "institution_ID_code" テーブルに対応するのだ
type InstitutionIDCode struct {
	InstitutionID   uint   `gorm:"primaryKey" json:"institution_id"`
	InstitutionCode string `json:"institution_code"`
}

func (InstitutionIDCode) TableName() string {
	return "institution_id_code"
}

// CollectionIDCode は "collection_ID_code" テーブルに対応するのだ
type CollectionIDCode struct {
	CollectionID   uint   `gorm:"primaryKey" json:"collection_id"`
	CollectionCode string `json:"collection_code"`
}

func (CollectionIDCode) TableName() string {
	return "collection_id_code"
}

// SpecimenMethod は "specimen_methods" テーブルに対応するのだ
type SpecimenMethod struct {
	SpecimenMethodsID  uint   `gorm:"primaryKey" json:"specimen_methods_id"`
	MethodCommonName string `json:"method_common_name"`
	PageID             uint   `json:"page_id"`

	// 関連
	WikiPage WikiPage `gorm:"foreignKey:PageID" json:"wiki_page"`
}

// Specimen は "specimen" テーブルに対応するのだ
type Specimen struct {
	SpecimenID        uint `gorm:"primaryKey" json:"specimen_id"`
	OccurrenceID      uint `json:"occurrence_id"`
	SpecimenMethodID  *uint `json:"specimen_method_id"`
	InstitutionID     *uint `json:"institution_id"`
	CollectionID      *uint `gorm:"column:collectionid" json:"collection_id"` // SQLのカラム名が小文字なので合わせる

	// 関連
	Occurrence       Occurrence        `gorm:"foreignKey:OccurrenceID" json:"occurrence"`
	SpecimenMethod   SpecimenMethod    `gorm:"foreignKey:SpecimenMethodID" json:"specimen_method"`
	InstitutionIDCode InstitutionIDCode `gorm:"foreignKey:InstitutionID" json:"institution_id_code"`
	CollectionIDCode  CollectionIDCode  `gorm:"foreignKey:CollectionID" json:"collection_id_code"`
}

func (Specimen) TableName() string {
	return "specimen"
}

// MakeSpecimen は "make_specimen" テーブルに対応するのだ
type MakeSpecimen struct {
	MakeSpecimenID   uint      `gorm:"primaryKey" json:"make_specimen_id"`
	OccurrenceID     uint      `json:"occurrence_id"`
	UserID           uint      `json:"user_id"`
	SpecimenID       uint      `json:"specimen_id"`
	Date             *time.Time `json:"date"`
	SpecimenMethodID *uint      `json:"specimen_method_id"`
	CreatedAt        time.Time `gorm:"default:now()" json:"created_at"`
	Timezone         int16     `gorm:"not null" json:"timezone"`

	// 関連
	Occurrence     Occurrence     `gorm:"foreignKey:OccurrenceID" json:"occurrence"`
	User           User           `gorm:"foreignKey:UserID" json:"user"`
	Specimen       Specimen       `gorm:"foreignKey:SpecimenID" json:"specimen"`
	SpecimenMethod SpecimenMethod `gorm:"foreignKey:SpecimenMethodID" json:"specimen_method"`
}

func (MakeSpecimen) TableName() string {
	return "make_specimen"
}
