
// backend/internal/service/specimen_service.go
package service

import (
	"github.com/saku-730/specimen-web/backend/internal/model"
	"github.com/saku-730/specimen-web/backend/internal/repository"
	"gorm.io/gorm"
)

// CreateSpecimenRequest は標本作成時のリクエストボディを表すのだ
type CreateSpecimenRequest struct {
	OccurrenceID     uint `json:"occurrence_id"`
	SpecimenMethodID uint `json:"specimen_method_id"`
	InstitutionID    uint `json:"institution_id"`
	CollectionID     uint `json:"collection_id"`
}

// UpdateSpecimenRequest は標本更新時のリクエストボディを表すのだ
type UpdateSpecimenRequest struct {
	OccurrenceID     uint `json:"occurrence_id"`
	SpecimenMethodID uint `json:"specimen_method_id"`
	InstitutionID    uint `json:"institution_id"`
	CollectionID     uint `json:"collection_id"`
}

// SpecimenService は標本関連のビジネスロジックのインターフェースなのだ
type SpecimenService interface {
	GetSpecimenByID(id uint) (*model.Specimen, error)
	GetAllSpecimens() ([]model.Specimen, error)
	CreateSpecimen(req CreateSpecimenRequest) (*model.Specimen, error)
	UpdateSpecimen(id uint, req UpdateSpecimenRequest) (*model.Specimen, error)
	DeleteSpecimen(id uint) error
}

type specimenService struct {
	db   *gorm.DB
	repo repository.SpecimenRepository
}

// NewSpecimenService は新しいサービスを生成するのだ
func NewSpecimenService(db *gorm.DB, repo repository.SpecimenRepository) SpecimenService {
	return &specimenService{db: db, repo: repo}
}

// GetSpecimenByID はIDで標本を1件取得するのだ
func (s *specimenService) GetSpecimenByID(id uint) (*model.Specimen, error) {
	return s.repo.FindByID(id)
}

// GetAllSpecimens は全ての標本を取得するのだ
func (s *specimenService) GetAllSpecimens() ([]model.Specimen, error) {
	return s.repo.FindAll()
}

// CreateSpecimen は新しい標本を作成するのだ
func (s *specimenService) CreateSpecimen(req CreateSpecimenRequest) (*model.Specimen, error) {
	newSpecimen := &model.Specimen{
		OccurrenceID:     req.OccurrenceID,
		SpecimenMethodID: req.SpecimenMethodID,
		InstitutionID:    req.InstitutionID,
		CollectionID:     req.CollectionID,
	}

	var createdSpecimen *model.Specimen
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		createdSpecimen, err = s.repo.Create(tx, newSpecimen)
		return err
	})

	if err != nil {
		return nil, err
	}
	return createdSpecimen, nil
}

// UpdateSpecimen は標本を更新するのだ
func (s *specimenService) UpdateSpecimen(id uint, req UpdateSpecimenRequest) (*model.Specimen, error) {
	var updatedSpecimen *model.Specimen
	err := s.db.Transaction(func(tx *gorm.DB) error {
		target, err := s.repo.FindByID(id)
		if err != nil {
			return err
		}

		target.OccurrenceID = req.OccurrenceID
		target.SpecimenMethodID = req.SpecimenMethodID
		target.InstitutionID = req.InstitutionID
		target.CollectionID = req.CollectionID

		updatedSpecimen, err = s.repo.Update(tx, target)
		return err
	})

	if err != nil {
		return nil, err
	}
	return updatedSpecimen, nil
}

// DeleteSpecimen はIDを元に標本を削除するのだ
func (s *specimenService) DeleteSpecimen(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Delete(tx, id)
	})
}
