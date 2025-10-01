
// backend/internal/repository/identification_repository.go
package repository

import (
	"github.com/saku-730/specimen-web/backend/internal/model"
	"gorm.io/gorm"
)

// IdentificationRepository は同定情報関連のデータ操作の契約書なのだ
type IdentificationRepository interface {
	FindByID(id uint) (*model.Identification, error)
	FindAll() ([]model.Identification, error)
	Create(tx *gorm.DB, identification *model.Identification) (*model.Identification, error)
	Update(tx *gorm.DB, identification *model.Identification) (*model.Identification, error)
	Delete(tx *gorm.DB, id uint) error
}

type identificationRepository struct {
	db *gorm.DB
}

// NewIdentificationRepository は新しいリポジトリを生成するのだ
func NewIdentificationRepository(db *gorm.DB) IdentificationRepository {
	return &identificationRepository{db: db}
}

// FindByID はIDで同定情報を1件取得するのだ
func (r *identificationRepository) FindByID(id uint) (*model.Identification, error) {
	var identification model.Identification
	if err := r.db.First(&identification, id).Error; err != nil {
		return nil, err
	}
	return &identification, nil
}

// FindAll は全ての同定情報を取得するのだ
func (r *identificationRepository) FindAll() ([]model.Identification, error) {
	var identifications []model.Identification
	if err := r.db.Find(&identifications).Error; err != nil {
		return nil, err
	}
	return identifications, nil
}

// Create は新しい同定情報を作成するのだ
func (r *identificationRepository) Create(tx *gorm.DB, identification *model.Identification) (*model.Identification, error) {
	if err := tx.Create(identification).Error; err != nil {
		return nil, err
	}
	return identification, nil
}

// Update は同定情報を更新するのだ
func (r *identificationRepository) Update(tx *gorm.DB, identification *model.Identification) (*model.Identification, error) {
	if err := tx.Save(identification).Error; err != nil {
		return nil, err
	}
	return identification, nil
}

// Delete はIDを元に同定情報を削除するのだ
func (r *identificationRepository) Delete(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&model.Identification{}, id).Error; err != nil {
		return err
	}
	return nil
}
