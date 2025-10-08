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
	Create(tx *gorm.DB, occurrence *model.Occurrence) (*model.Occurrence, error)
	Search(params SearchRepository) ([]handler.SerachResult, error)
}

type occurrenceRepository struct {
	db *gorm.DB
}

// NewOccurrenceRepository make new Occurence data
func NewOccurrenceRepository(db *gorm.DB) OccurrenceRepository {
	return &occurrenceRepository{db: db}
}

func (r *occurrenceRepository) Create(tx *gorm.DB, occurrence *model.Occurrence) (*model.Occurrence, error) {
	if err := tx.Create(occurrence).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *occurrenceRepository) Search(params SearchParams) ([]handler.SearchResult, error){
	var occurrences []model.Occurrence

	query := r.db.Model(&model.Occurrence{})

}
