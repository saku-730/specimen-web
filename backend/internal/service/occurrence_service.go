// backend/internal/service/occurrence_service.go
package service

import (
	"time"

	"github.com/saku-730/specimen-web/backend/internal/model"
	"github.com/saku-730/specimen-web/backend/internal/repository"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// --- Structs for simple CRUD ---

type CreateOccurrenceRequest struct {
	ProjectID         uint      `json:"project_id"`
	UserID            uint      `json:"user_id"`
	IndividualID      *int      `json:"individual_id"`
	Lifestage         string    `json:"lifestage"`
	Sex               string    `json:"sex"`
	ClassificationID  uint      `json:"classification_id"`
	PlaceID           uint      `json:"place_id"`
	AttachmentGroupID *int      `json:"attachment_group_id"`
	BodyLength        *float64  `json:"body_length"`
	LanguageID        uint      `json:"language_id"`
	Note              string    `json:"note"`
	Timezone          int16     `json:"timezone"`
}

type UpdateOccurrenceRequest struct {
	ProjectID         uint      `json:"project_id"`	
	UserID            uint      `json:"user_id"`
	IndividualID      *int      `json:"individual_id"`
	Lifestage         string    `json:"lifestage"`
	Sex               string    `json:"sex"`
	ClassificationID  uint      `json:"classification_id"`
	PlaceID           uint      `json:"place_id"`
	AttachmentGroupID *int      `json:"attachment_group_id"`
	BodyLength        *float64  `json:"body_length"`
	LanguageID        uint      `json:"language_id"`
	Note              string    `json:"note"`
	Timezone          int16     `json:"timezone"`
}

// --- Structs for Full Occurrence Form ---

type FullOccurrenceRequest struct {
	Occurrence     OccurrencePayload     `json:"occurrence"`
	Classification ClassificationPayload `json:"classification"`
	Place          PlacePayload          `json:"place"`
	Observation    ObservationPayload    `json:"observation"`
	Specimen       SpecimenPayload       `json:"specimen"`
	MakeSpecimen   MakeSpecimenPayload   `json:"make_specimen"`
	Identification IdentificationPayload `json:"identification"`
}

type OccurrencePayload struct {
	ProjectID    uint      `json:"project_id"`
	UserID       uint      `json:"user_id"`
	IndividualID *int      `json:"individual_id"`
	Lifestage    string    `json:"lifestage"`	
	Sex          string    `json:"sex"`
	BodyLength   *float64  `json:"body_length"`
	CreatedAt    string    `json:"created_at"`
	Timezone     int16     `json:"timezone"`
	LanguageID   uint      `json:"language_id"`
	Note         string    `json:"note"`
}

type ClassificationPayload struct {
	ClassClassification datatypes.JSON `json:"class_classification"`
}

type PlacePayload struct {
	Coordinates   string `json:"coordinates"`
	PlaceNameJSON struct {
		ClassPlaceName datatypes.JSON `json:"class_place_name"`
	} `json:"place_name_json"`
}

type ObservationPayload struct {
	UserID              uint   `json:"user_id"`
	ObservationMethodID uint   `json:"observation_method_id"`
	Behavior            string `json:"behavior"`
	ObservedAt          string `json:"observed_at"`
	Timezone            int16  `json:"timezone"`
}

type SpecimenPayload struct {
	SpecimenMethodID uint `json:"specimen_method_id"`
	InstitutionID    uint `json:"institution_id"`
	CollectionID     uint `json:"collection_id"`
}

type MakeSpecimenPayload struct {
	UserID    uint   `json:"user_id"`
	Date      string `json:"date"`
	CreatedAt string `json:"created_at"`
	Timezone  int16  `json:"timezone"`
}

type IdentificationPayload struct {
	UserID          uint   `json:"user_id"`
	SourceInfo      string `json:"source_info"`
	IdentificatedAt string `json:"identificated_at"`
	Timezone        int16  `json:"timezone"`
}


type OccurrenceService interface {
	GetOccurrenceByID(id uint) (*model.Occurrence, error)
	GetAllOccurrences() ([]model.Occurrence, error)
	FindOccurrencesByConditions(conditions *model.Occurrence) ([]model.Occurrence, error)
	CreateOccurrence(req CreateOccurrenceRequest) (*model.Occurrence, error)
	UpdateOccurrence(id uint, req UpdateOccurrenceRequest) (*model.Occurrence, error)
	DeleteOccurrence(id uint) error
	GetAllLanguages() ([]model.Language, error)
	CreateFullOccurrence(req FullOccurrenceRequest) error
}

type occurrenceService struct {
	db   *gorm.DB
	repo repository.OccurrenceRepository
}

// NewOccurrenceService は新しいサービスを生成するのだ
func NewOccurrenceService(db *gorm.DB, repo repository.OccurrenceRepository) OccurrenceService {
	return &occurrenceService{db: db, repo: repo}
}

// uintToPtr は uint が 0 でなければそのポインタを、0 なら nil を返すのだ
func uintToPtr(val uint) *uint {
	if val == 0 {
		return nil
	}
	return &val
}

// CreateFullOccurrence はフォームからの全データを受け取ってまとめて登録するのだ
func (s *occurrenceService) CreateFullOccurrence(req FullOccurrenceRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create Classification
		classification := model.ClassificationJSON{
			ClassClassification: req.Classification.ClassClassification,
		}
		if err := tx.Create(&classification).Error; err != nil {
			return err
		}

		// 2. Create Place
		placeName := model.PlaceNameJSON{
			ClassPlaceName: req.Place.PlaceNameJSON.ClassPlaceName,
		}
		if err := tx.Create(&placeName).Error; err != nil {
			return err
		}
		var coords *string
		if req.Place.Coordinates != "POINT( )" {
			coords = &req.Place.Coordinates
		}
		place := model.Place{
			Coordinates: coords,
			PlaceNameID: placeName.PlaceNameID,
		}
		if err := tx.Create(&place).Error; err != nil {
			return err
		}

		// 3. Create Occurrence
		const layout = "2006-01-02T15:04"
		createdAt, err := time.Parse(layout, req.Occurrence.CreatedAt)
		if err != nil {
			return err
		}
		occurrence := model.Occurrence{
			ProjectID:        uintToPtr(req.Occurrence.ProjectID),
			UserID:           req.Occurrence.UserID,
			IndividualID:     req.Occurrence.IndividualID,
			Lifestage:        req.Occurrence.Lifestage,
			Sex:              req.Occurrence.Sex,
			ClassificationID: classification.ClassificationID,
			PlaceID:          uintToPtr(place.PlaceID),
			BodyLength:       req.Occurrence.BodyLength,
			LanguageID:       uintToPtr(req.Occurrence.LanguageID),
			Note:             req.Occurrence.Note,
			CreatedAt:        createdAt,
			Timezone:         req.Occurrence.Timezone,
		}
		if err := tx.Create(&occurrence).Error; err != nil {
			return err
		}

		// 4. Create Observation
		observedAt, err := time.Parse(layout, req.Observation.ObservedAt)
		if err != nil {
			return err
		}
		observation := model.Observation{
			UserID:              req.Observation.UserID,
			OccurrenceID:        occurrence.OccurrenceID,
			ObservationMethodID: uintToPtr(req.Observation.ObservationMethodID),
			Behavior:            req.Observation.Behavior,
			ObservedAt:          observedAt,
			Timezone:            req.Observation.Timezone,
		}
		if err := tx.Create(&observation).Error; err != nil {
			return err
		}

		// 5. Create Specimen
		specimen := model.Specimen{
			OccurrenceID:     occurrence.OccurrenceID,
			SpecimenMethodID: uintToPtr(req.Specimen.SpecimenMethodID),
			InstitutionID:    uintToPtr(req.Specimen.InstitutionID),
			CollectionID:     uintToPtr(req.Specimen.CollectionID),
		}
		if err := tx.Create(&specimen).Error; err != nil {
			return err
		}

		// 6. Create MakeSpecimen
		makeDate, err := time.Parse("2006-01-02", req.MakeSpecimen.Date)
		if err != nil {
			return err
		}
		makeCreatedAt, err := time.Parse(layout, req.MakeSpecimen.CreatedAt)
		if err != nil {
			return err
		}
		makeSpecimen := model.MakeSpecimen{
			OccurrenceID:     occurrence.OccurrenceID,
			UserID:           req.MakeSpecimen.UserID,
			SpecimenID:       specimen.SpecimenID,
			Date:             &makeDate,
			SpecimenMethodID: uintToPtr(req.Specimen.SpecimenMethodID), // 標本テーブルのメソッドIDを流用
			CreatedAt:        makeCreatedAt,
			Timezone:         req.MakeSpecimen.Timezone,
		}
		if err := tx.Create(&makeSpecimen).Error; err != nil {
			return err
		}

		// 7. Create Identification
		identificatedAt, err := time.Parse(layout, req.Identification.IdentificatedAt)
		if err != nil {
			return err
		}
		identification := model.Identification{
			UserID:          req.Identification.UserID,
			OccurrenceID:    occurrence.OccurrenceID,
			SourceInfo:      req.Identification.SourceInfo,
			IdentificatedAt: identificatedAt,
			Timezone:        req.Identification.Timezone,
		}
		if err := tx.Create(&identification).Error; err != nil {
			return err
		}

		return nil // commit
	})
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
		ProjectID:         uintToPtr(req.ProjectID),
		UserID:            req.UserID,
		IndividualID:      req.IndividualID,
		Lifestage:         req.Lifestage,
		Sex:               req.Sex,
		ClassificationID:  req.ClassificationID,
		PlaceID:           uintToPtr(req.PlaceID),
		AttachmentGroupID: req.AttachmentGroupID,
		BodyLength:        req.BodyLength,
		LanguageID:        uintToPtr(req.LanguageID),
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
		target.ProjectID = uintToPtr(req.ProjectID)
		target.UserID = req.UserID
		target.IndividualID = req.IndividualID
		target.Lifestage = req.Lifestage
		target.Sex = req.Sex
		target.ClassificationID = req.ClassificationID
		target.PlaceID = uintToPtr(req.PlaceID)
		target.AttachmentGroupID = req.AttachmentGroupID
		target.BodyLength = req.BodyLength
		target.LanguageID = uintToPtr(req.LanguageID)
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

// GetAllLanguages は全ての言語を取得するのだ
func (s *occurrenceService) GetAllLanguages() ([]model.Language, error) {
	var languages []model.Language
	if err := s.db.Find(&languages).Error; err != nil {
		return nil, err
	}
	return languages, nil
}
