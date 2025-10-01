
// backend/internal/repository/attachment_repository.go
package repository

import (
	"github.com/saku-730/specimen-web/backend/internal/model"
	"gorm.io/gorm"
)

// AttachmentRepository は添付ファイル関連のデータ操作の契約書なのだ
type AttachmentRepository interface {
	FindByID(id uint) (*model.Attachment, error)
	FindAll() ([]model.Attachment, error)
	Create(tx *gorm.DB, attachment *model.Attachment) (*model.Attachment, error)
	Update(tx *gorm.DB, attachment *model.Attachment) (*model.Attachment, error)
	Delete(tx *gorm.DB, id uint) error
}

type attachmentRepository struct {
	db *gorm.DB
}

// NewAttachmentRepository は新しいリポジトリを生成するのだ
func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepository{db: db}
}

// FindByID はIDで添付ファイルを1件取得するのだ
func (r *attachmentRepository) FindByID(id uint) (*model.Attachment, error) {
	var attachment model.Attachment
	if err := r.db.First(&attachment, id).Error; err != nil {
		return nil, err
	}
	return &attachment, nil
}

// FindAll は全ての添付ファイルを取得するのだ
func (r *attachmentRepository) FindAll() ([]model.Attachment, error) {
	var attachments []model.Attachment
	if err := r.db.Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

// Create は新しい添付ファイルを作成するのだ
func (r *attachmentRepository) Create(tx *gorm.DB, attachment *model.Attachment) (*model.Attachment, error) {
	if err := tx.Create(attachment).Error; err != nil {
		return nil, err
	}
	return attachment, nil
}

// Update は添付ファイルを更新するのだ
func (r *attachmentRepository) Update(tx *gorm.DB, attachment *model.Attachment) (*model.Attachment, error) {
	if err := tx.Save(attachment).Error; err != nil {
		return nil, err
	}
	return attachment, nil
}

// Delete はIDを元に添付ファイルを削除するのだ
func (r *attachmentRepository) Delete(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&model.Attachment{}, id).Error; err != nil {
		return err
	}
	return nil
}
