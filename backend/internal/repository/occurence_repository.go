// internal/repository/occurrence_repository.go
package repository

import (
	"github.com/saku-730/specimen-web/backend/internal/model"

	"gorm.io/gorm"
)

type SearchParams struct{
	user
}

// OccurrenceRepository は発生情報関連のデータ操作の契約書なのだ
type OccurrenceRepository interface {
	FindByID(id uint) (*model.Occurrence, error)
	FindAll() ([]model.Occurrence, error)
	FindByConditions(conditions *model.Occurrence) ([]model.Occurrence, error)
	Create(tx *gorm.DB, occurrence *model.Occurrence) (*model.Occurrence, error)
	Update(tx *gorm.DB, occurrence *model.Occurrence) (*model.Occurrence, error)
	Delete(tx *gorm.DB, id uint) error
	Search(params SearchRepository) ([]handler.SerachResult, error)
}

type occurrenceRepository struct {
	db *gorm.DB
}

// NewOccurrenceRepository make new Occurence data
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

func (r *occurrenceRepository) Search(params SearchParams) ([]handler.SearchResult, error){
	var occurrences []model.Occurrence

	query := r.db.Model(&model.Occurrence{})

}
