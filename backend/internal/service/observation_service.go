
// backend/internal/service/observation_service.go
package service

import (
	"github.com/saku-730/specimen-web/backend/internal/model"
	"github.com/saku-730/specimen-web/backend/internal/repository"
	"gorm.io/gorm"
)

// CreateObservationRequest は観察情報作成時のリクエストボディを表すのだ
type CreateObservationRequest struct {
	UserID              uint   `json:"user_id"`
	OccurrenceID        uint   `json:"occurrence_id"`
	ObservationMethodID uint   `json:"observation_method_id"`
	Behavior            string `json:"behavior"`
	Timezone            int16  `json:"timezone"`
}

// UpdateObservationRequest は観察情報更新時のリクエストボディを表すのだ
type UpdateObservationRequest struct {
	UserID              uint   `json:"user_id"`
	OccurrenceID        uint   `json:"occurrence_id"`
	ObservationMethodID uint   `json:"observation_method_id"`
	Behavior            string `json:"behavior"`
	Timezone            int16  `json:"timezone"`
}

// ObservationService は観察情報関連のビジネスロジックのインターフェースなのだ
type ObservationService interface {
	GetObservationByID(id uint) (*model.Observation, error)
	GetAllObservations() ([]model.Observation, error)
	CreateObservation(req CreateObservationRequest) (*model.Observation, error)
	UpdateObservation(id uint, req UpdateObservationRequest) (*model.Observation, error)
	DeleteObservation(id uint) error
}

type observationService struct {
	db   *gorm.DB
	repo repository.ObservationRepository
}

// NewObservationService は新しいサービスを生成するのだ
func NewObservationService(db *gorm.DB, repo repository.ObservationRepository) ObservationService {
	return &observationService{db: db, repo: repo}
}

func (s *observationService) GetObservationByID(id uint) (*model.Observation, error) {
	return s.repo.FindByID(id)
}

func (s *observationService) GetAllObservations() ([]model.Observation, error) {
	return s.repo.FindAll()
}

func (s *observationService) CreateObservation(req CreateObservationRequest) (*model.Observation, error) {
	newObservation := &model.Observation{
		UserID:              req.UserID,
		OccurrenceID:        req.OccurrenceID,
		ObservationMethodID: req.ObservationMethodID,
		Behavior:            req.Behavior,
		Timezone:            req.Timezone,
	}

	var createdObservation *model.Observation
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		createdObservation, err = s.repo.Create(tx, newObservation)
		return err
	})

	if err != nil {
		return nil, err
	}
	return createdObservation, nil
}

func (s *observationService) UpdateObservation(id uint, req UpdateObservationRequest) (*model.Observation, error) {
	var updatedObservation *model.Observation
	err := s.db.Transaction(func(tx *gorm.DB) error {
		target, err := s.repo.FindByID(id)
		if err != nil {
			return err
		}

		target.UserID = req.UserID
		target.OccurrenceID = req.OccurrenceID
		target.ObservationMethodID = req.ObservationMethodID
		target.Behavior = req.Behavior
		target.Timezone = req.Timezone

		updatedObservation, err = s.repo.Update(tx, target)
		return err
	})

	if err != nil {
		return nil, err
	}
	return updatedObservation, nil
}

func (s *observationService) DeleteObservation(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Delete(tx, id)
	})
}
