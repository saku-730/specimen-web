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
	GetAllSpecimenMethods() ([]model.SpecimenMethod, error)
	GetAllInstitutionCodes() ([]model.InstitutionIDCode, error)
	GetAllCollectionCodes() ([]model.CollectionIDCode, error)
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
		SpecimenMethodID: uintToPtr(req.SpecimenMethodID),
		InstitutionID:    uintToPtr(req.InstitutionID),
		CollectionID:     uintToPtr(req.CollectionID),
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
		target.SpecimenMethodID = uintToPtr(req.SpecimenMethodID)
		target.InstitutionID = uintToPtr(req.InstitutionID)
		target.CollectionID = uintToPtr(req.CollectionID)

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

// GetAllSpecimenMethods は全ての標本作成方法を取得するのだ
func (s *specimenService) GetAllSpecimenMethods() ([]model.SpecimenMethod, error) {
	var methods []model.SpecimenMethod
	if err := s.db.Find(&methods).Error; err != nil {
		return nil, err
	}
	return methods, nil
}

// GetAllInstitutionCodes は全ての機関コードを取得するのだ
func (s *specimenService) GetAllInstitutionCodes() ([]model.InstitutionIDCode, error) {
	var codes []model.InstitutionIDCode
	if err := s.db.Find(&codes).Error; err != nil {
		return nil, err
	}
	return codes, nil
}

// GetAllCollectionCodes は全てのコレクションコードを取得するのだ
func (s *specimenService) GetAllCollectionCodes() ([]model.CollectionIDCode, error) {
	var codes []model.CollectionIDCode
	if err := s.db.Find(&codes).Error; err != nil {
		return nil, err
	}
	return codes, nil
}