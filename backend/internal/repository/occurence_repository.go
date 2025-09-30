// internal/repository/occurrence_repository.go
package repository

import (
	"github.com/saku-730/specimen-web/backend/internal/model"

	"gorm.io/gorm"
)

// OccurrenceRepository は発生情報関連のデータ操作の契約書なのだ
type OccurrenceRepository interface {
	FindByID(id uint) (*model.Occurrence, error)
	FindAll() ([]model.Occurrence, error)
	FindByConditions(conditions *model.Occurrence) ([]model.Occurrence, error)
	Create(occurrence *model.Occurrence) (*model.Occurrence, error)
	Update(occurrence *model.Occurrence) (*model.Occurrence, error)
	Delete(id uint) error
}

type occurrenceRepository struct {
	db *gorm.DB
}

// NewOccurrenceRepository は新しいリポジトリを生成するのだ
func NewOccurrenceRepository(db *gorm.DB) OccurrenceRepository {
	return &occurrenceRepository{db: db}
}

// FindByID はIDで発生情報を1件取得する。関連する添付ファイルも一緒に読み込むのだ
func (r *occurrenceRepository) FindByID(id uint) (*model.Occurrence, error) {
	var occurrence model.Occurrence
	if err := r.db.Preload("Attachments").First(&occurrence, id).Error; err != nil {
		return nil, err
	}
	return &occurrence, nil
}

// FindAll は全ての発生情報を取得するのだ
func (r *occurrenceRepository) FindAll() ([]model.Occurrence, error) {
	var occurrences []model.Occurrence
	if err := r.db.Find(&occurrences).Error; err != nil {
		return nil, err
	}
	return occurrences, nil
}

// FindByConditions は与えられた条件で発生情報を検索するのだ
func (r *occurrenceRepository) FindByConditions(conditions *model.Occurrence) ([]model.Occurrence, error) {
	var occurrences []model.Occurrence
	// 構造体のゼロ値でないフィールドを元に、自動でWHERE句を組み立ててくれるのだ
	if err := r.db.Where(conditions).Find(&occurrences).Error; err != nil {
		return nil, err
	}
	return occurrences, nil
}

// Create は新しい発生情報を作成するのだ
func (r *occurrenceRepository) Create(occurrence *model.Occurrence) (*model.Occurrence, error) {
	if err := r.db.Create(occurrence).Error; err != nil {
		return nil, err
	}
	return occurrence, nil
}

// Update は発生情報を更新するのだ
func (r *occurrenceRepository) Update(occurrence *model.Occurrence) (*model.Occurrence, error) {
	// ゼロ値でないフィールドだけを安全に更新するのだ
	if err := r.db.Model(&model.Occurrence{OccurrenceID: occurrence.OccurrenceID}).Updates(occurrence).Error; err != nil {
		return nil, err
	}
	// 更新後の最新の情報を取得して返すのだ
	return r.FindByID(occurrence.OccurrenceID)
}

// Delete はIDを元に発生情報を削除するのだ
func (r *occurrenceRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.Occurrence{}, id).Error; err != nil {
		return err
	}
	return nil
}

