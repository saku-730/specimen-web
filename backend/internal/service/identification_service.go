
// backend/internal/service/identification_service.go
package service

import (
	"github.com/saku-730/specimen-web/backend/internal/model"
	"github.com/saku-730/specimen-web/backend/internal/repository"
	"gorm.io/gorm"
)

// CreateIdentificationRequest は同定情報作成時のリクエストボディを表すのだ
type CreateIdentificationRequest struct {
	UserID       uint   `json:"user_id"`
	OccurrenceID uint   `json:"occurrence_id"`
	SourceInfo   string `json:"source_info"`
	Timezone     int16  `json:"timezone"`
}

// UpdateIdentificationRequest は同定情報更新時のリクエストボディを表すのだ
type UpdateIdentificationRequest struct {
	UserID       uint   `json:"user_id"`
	OccurrenceID uint   `json:"occurrence_id"`
	SourceInfo   string `json:"source_info"`
	Timezone     int16  `json:"timezone"`
}

// IdentificationService は同定情報関連のビジネスロジックのインターフェースなのだ
type IdentificationService interface {
	GetIdentificationByID(id uint) (*model.Identification, error)
	GetAllIdentifications() ([]model.Identification, error)
	CreateIdentification(req CreateIdentificationRequest) (*model.Identification, error)
	UpdateIdentification(id uint, req UpdateIdentificationRequest) (*model.Identification, error)
	DeleteIdentification(id uint) error
}

type identificationService struct {
	db   *gorm.DB
	repo repository.IdentificationRepository
}

// NewIdentificationService は新しいサービスを生成するのだ
func NewIdentificationService(db *gorm.DB, repo repository.IdentificationRepository) IdentificationService {
	return &identificationService{db: db, repo: repo}
}

func (s *identificationService) GetIdentificationByID(id uint) (*model.Identification, error) {
	return s.repo.FindByID(id)
}

func (s *identificationService) GetAllIdentifications() ([]model.Identification, error) {
	return s.repo.FindAll()
}

func (s *identificationService) CreateIdentification(req CreateIdentificationRequest) (*model.Identification, error) {
	newIdentification := &model.Identification{
		UserID:       req.UserID,
		OccurrenceID: req.OccurrenceID,
		SourceInfo:   req.SourceInfo,
		Timezone:     req.Timezone,
	}

	var createdIdentification *model.Identification
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		createdIdentification, err = s.repo.Create(tx, newIdentification)
		return err
	})

	if err != nil {
		return nil, err
	}
	return createdIdentification, nil
}

func (s *identificationService) UpdateIdentification(id uint, req UpdateIdentificationRequest) (*model.Identification, error) {
	var updatedIdentification *model.Identification
	err := s.db.Transaction(func(tx *gorm.DB) error {
		target, err := s.repo.FindByID(id)
		if err != nil {
			return err
		}

		target.UserID = req.UserID
		target.OccurrenceID = req.OccurrenceID
		target.SourceInfo = req.SourceInfo
		target.Timezone = req.Timezone

		updatedIdentification, err = s.repo.Update(tx, target)
		return err
	})

	if err != nil {
		return nil, err
	}
	return updatedIdentification, nil
}

func (s *identificationService) DeleteIdentification(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Delete(tx, id)
	})
}
