// backend/internal/service/occurrence_service.go
package service

import (
	"time"

	"github.com/saku-730/specimen-web/backend/internal/model"
	"github.com/saku-730/specimen-web/backend/internal/repository"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// --- Structs for Occurrence Search ---

type SearchRequest struct {
	UserID      *uint   `form:"user_id"`
	
}


type SearchResponse struct {
	UserID      *uint
}

type ClassificationJSONB struct {
	Kingdom string `json;"kingdom"`
	Phylum string `json;"phylum"`
	Class string `json;"class"`
	Order string `json;"order"`
	Family string `json;"family"`
	Genus string `json;"genus"`
	Species string `json:"species"`
	//others string `json:"others"`
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


//---

type OccurrenceService interface {
	GetAllLanguages() ([]model.Language, error)
	CreateFullOccurrence(req FullOccurrenceRequest) error
	Search(req SearchRequest)([]SearchResponse, error)
}

type occurrenceService struct {
	db   *gorm.DB
	repo repository.OccurrenceRepository
}


// NewOccurrenceService は新しいサービスを生成するのだ
func NewOccurrenceService(db *gorm.DB, repo repository.OccurrenceRepository) OccurrenceService {
	return &occurrenceService{db: db, repo: repo}
}

// Search

func(s *occurrenceService) Search(req SeqrchRequest)([]SearchResponse, error){

	repoParams := repository.SearchParams{
		UserID:		req.UserID,

	}

	raw_results, err := s.occurrenceRepo.Search(repoParams)
	if err != nil{
		return nil, err
	}
	responses := make([]SearchOccurrencesResponse, 0, len(occurrences))
	for _, search_results := range raw_results {
		// JSONBデータからの値の取り出し（安全なチェック）
		var classificationData ClassificationJSONB

		if search_results.ClassificationJSON != nil && search_results.ClassificationJSON.ClassClassification.Valid {
			var classificationData ClassificationJSONB

			json.Unmarshal(search_results.ClassificationJSON.ClassClassification.RawMessage, &classificationData)

		dto := SearchOccurrencesResponse{
			OccurrenceID: search_results.ID,
			Note:         search_results.Note,
			CreatedAt:    search_results.CreatedAt,
			UserName:     search_results.User.UserName,       // Preloadしたデータを使う
			ProjectName:  search_results.Project.ProjectName,   // Preloadしたデータを使う
			Species:      classificationData.Species,
		}
		responses = append(responses, dto)
	}
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

		return nil 
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
