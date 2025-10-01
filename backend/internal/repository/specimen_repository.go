// internal/repository/specimen_repository.go
package repository

import (
	"github.com/saku-730/specimen-web/backend/internal/model"

	"gorm.io/gorm"
)

// SpecimenRepository は標本関連のデータ操作の契約書なのだ
type SpecimenRepository interface {
	FindByID(id uint) (*model.Specimen, error)
	FindAll() ([]model.Specimen, error)
	FindByConditions(conditions *model.Specimen) ([]model.Specimen, error)
	Create(tx *gorm.DB, specimen *model.Specimen) (*model.Specimen, error)
	Update(tx *gorm.DB, specimen *model.Specimen) (*model.Specimen, error)
	Delete(tx *gorm.DB, id uint) error
}

type specimenRepository struct {
	db *gorm.DB
}

// NewSpecimenRepository は新しいリポジトリを生成するのだ
func NewSpecimenRepository(db *gorm.DB) SpecimenRepository {
	return &specimenRepository{db: db}
}

// FindByID はIDで標本を1件取得する。関連情報も一緒に取得するのだ
func (r *specimenRepository) FindByID(id uint) (*model.Specimen, error) {
	var specimen model.Specimen
	// 標本は多くの情報と紐づいているので、必要なものをPreloadで指定するのだ
	err := r.db.Preload("Occurrence").
		Preload("SpecimenMethod").
		Preload("InstitutionIDCode").
		Preload("CollectionIDCode").
		First(&specimen, id).Error
	if err != nil {
		return nil, err
	}
	return &specimen, nil
}

// FindAll は全ての標本を取得するのだ
func (r *specimenRepository) FindAll() ([]model.Specimen, error) {
	var specimens []model.Specimen
	if err := r.db.Find(&specimens).Error; err != nil {
		return nil, err
	}
	return specimens, nil
}

// FindByConditions は与えられた条件で標本を検索するのだ
func (r *specimenRepository) FindByConditions(conditions *model.Specimen) ([]model.Specimen, error) {
	var specimens []model.Specimen
	if err := r.db.Where(conditions).Find(&specimens).Error; err != nil {
		return nil, err
	}
	return specimens, nil
}

// Create は新しい標本を作成するのだ
func (r *specimenRepository) Create(tx *gorm.DB, specimen *model.Specimen) (*model.Specimen, error) {
	if err := tx.Create(specimen).Error; err != nil {
		return nil, err
	}
	return specimen, nil
}

// Update は標本情報を更新するのだ
func (r *specimenRepository) Update(tx *gorm.DB, specimen *model.Specimen) (*model.Specimen, error) {
	if err := tx.Model(&model.Specimen{SpecimenID: specimen.SpecimenID}).Updates(specimen).Error; err != nil {
		return nil, err
	}
	return specimen, nil
}

// Delete はIDを元に標本を削除するのだ
func (r *specimenRepository) Delete(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&model.Specimen{}, id).Error; err != nil {
		return err
	}
	return nil
}

