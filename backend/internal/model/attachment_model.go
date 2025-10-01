// internal/model/attachment_models.go
package model

// FileType は "file_types" テーブルに対応するのだ
type FileType struct {
	FileTypeID uint   `gorm:"primaryKey" json:"file_type_id"`
	TypeName   string `json:"type_name"`
}

// FileExtension は "file_extensions" テーブルに対応するのだ
type FileExtension struct {
	ExtensionID  uint   `gorm:"primaryKey" json:"extension_id"`
	ExtensionText string `gorm:"size:255" json:"extension_text"`
	FileTypeID   uint   `json:"file_type_id"`
}

// Attachment は "attachments" テーブルに対応するのだ
type Attachment struct {
	AttachmentID uint   `gorm:"primaryKey" json:"attachment_id"`
	FilePath     string `gorm:"not null" json:"file_path"`
	ExtensionID  uint   `json:"extension_id"`
	UserID       uint   `json:"user_id"`
}

// AttachmentGroup は "attachment_goup" (groupのtypo?) 中間テーブルに対応するのだ
type AttachmentGroup struct {
	OccurrenceID uint `gorm:"primaryKey" json:"occurrence_id"`
	AttachmentID uint `gorm:"primaryKey" json:"attachment_id"`
	Priority     *int `json:"priority"`
}
