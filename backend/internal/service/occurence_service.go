
// backend/internal/service/occurrence_service.go
package service

import (
	"github.com/saku-730/specimen-web/backend/internal/model"
	"github.com/saku-730/specimen-web/backend/internal/repository"
	"gorm.io/gorm"
)

// CreateOccurrenceRequest は発生情報作成時のリクエストボディを表すのだ
type CreateOccurrenceRequest struct {
	ProjectID         uint      `json:"project_id"`
	UserID            uint      `json:"user_id"`
	IndividualID      *int      `json:"individual_id"`
	Lifestage         string    `json:"lifestage"`
	Sex               string    `json:"sex"`
	ClassificationID  uint      `json:"classification_id"`
	PlaceID           uint      `json:"place_id"`
	AttachmentGroupID *int      `json:"attachment_group_id"`
	BodyLength        float64   `json:"body_length"`
	LanguageID        uint      `json:"language_id"`
	Note              string    `json:"note"`
	Timezone          int16     `json:"timezone"`
}

// UpdateOccurrenceRequest は発生情報更新時のリクエストボディを表すのだ
type UpdateOccurrenceRequest struct {
	ProjectID         uint      `json:"project_id"`
	UserID            uint      `json:"user_id"`
	IndividualID      *int      `json:"individual_id"`
	Lifestage         string    `json:"lifestage"`
	Sex               string    `json:"sex"`
	ClassificationID  uint      `json:"classification_id"`
	PlaceID           uint      `json:"place_id"`
	AttachmentGroupID *int      `json:"attachment_group_id"`
	BodyLength        float64   `json:"body_length"`
	LanguageID        uint      `json:"language_id"`
	Note              string    `json:"note"`
	Timezone          int16     `json:"timezone"`
}


// OccurrenceService は発生情報関連のビジネスロジックのインターフェースなのだ
type OccurrenceService interface {
	GetOccurrenceByID(id uint) (*model.Occurrence, error)
	GetAllOccurrences() ([]model.Occurrence, error)
	FindOccurrencesByConditions(conditions *model.Occurrence) ([]model.Occurrence, error)
	CreateOccurrence(req CreateOccurrenceRequest) (*model.Occurrence, error)
	UpdateOccurrence(id uint, req UpdateOccurrenceRequest) (*model.Occurrence, error)
	DeleteOccurrence(id uint) error
}

type occurrenceService struct {
	db       *gorm.DB
	repo repository.OccurrenceRepository
}

// NewOccurrenceService は新しいサービスを生成するのだ
func NewOccurrenceService(db *gorm.DB, repo repository.OccurrenceRepository) OccurrenceService {
	return &occurrenceService{db: db, repo: repo}
}

// GetOccurrenceByID はIDで発生情報を1件取得するのだ
func (s *occurrenceService) GetOccurrenceByID(id uint) (*model.Occurrence, error) {
	return s.repo.FindByID(id)
}

// GetAllOccurrences は全ての発生情報を取得するのだ
func (s *occurrenceService) GetAllOccurrences() ([]model.Occurrence, error) {
	return s.repo.FindAll()
}

// FindOccurrencesByConditions は与えられた条件で発生情報を検索するのだ
func (s *occurrenceService) FindOccurrencesByConditions(conditions *model.Occurrence) ([]model.Occurrence, error) {
	return s.repo.FindByConditions(conditions)
}

// CreateOccurrence は新しい発生情報を作成するのだ
func (s *occurrenceService) CreateOccurrence(req CreateOccurrenceRequest) (*model.Occurrence, error) {
	newOccurrence := &model.Occurrence{
		ProjectID:         req.ProjectID,
		UserID:            req.UserID,
		IndividualID:      req.IndividualID,
		Lifestage:         req.Lifestage,
		Sex:               req.Sex,
		ClassificationID:  req.ClassificationID,
		PlaceID:           req.PlaceID,
		AttachmentGroupID: req.AttachmentGroupID,
		BodyLength:        req.BodyLength,
		LanguageID:        req.LanguageID,
		Note:              req.Note,
		Timezone:          req.Timezone,
	}

	var createdOccurrence *model.Occurrence
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		createdOccurrence, err = s.repo.Create(tx, newOccurrence)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return createdOccurrence, nil
}

// UpdateOccurrence は発生情報を更新するのだ
func (s *occurrenceService) UpdateOccurrence(id uint, req UpdateOccurrenceRequest) (*model.Occurrence, error) {
	var updatedOccurrence *model.Occurrence
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 最初に、更新対象のレコードが存在するか確認するのだ
		target, err := s.repo.FindByID(id)
		if err != nil {
			return err // レコードが見つからなければエラー
		}

		// リクエストの情報を元に、更新するフィールドを設定するのだ
		target.ProjectID = req.ProjectID
		target.UserID = req.UserID
		target.IndividualID = req.IndividualID
		target.Lifestage = req.Lifestage
		target.Sex = req.Sex
		target.ClassificationID = req.ClassificationID
		target.PlaceID = req.PlaceID
		target.AttachmentGroupID = req.AttachmentGroupID
		target.BodyLength = req.BodyLength
		target.LanguageID = req.LanguageID
		target.Note = req.Note
		target.Timezone = req.Timezone
		
		updatedOccurrence, err = s.repo.Update(tx, target)
		return err
	})

	if err != nil {
		return nil, err
	}
	return updatedOccurrence, nil
}

// DeleteOccurrence はIDを元に発生情報を削除するのだ
func (s *occurrenceService) DeleteOccurrence(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Delete(tx, id)
	})
}
