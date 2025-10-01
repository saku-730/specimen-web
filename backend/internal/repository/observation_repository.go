
// backend/internal/repository/observation_repository.go
package repository

import (
	"github.com/saku-730/specimen-web/backend/internal/model"
	"gorm.io/gorm"
)

// ObservationRepository は観察情報関連のデータ操作の契約書なのだ
type ObservationRepository interface {
	FindByID(id uint) (*model.Observation, error)
	FindAll() ([]model.Observation, error)
	Create(tx *gorm.DB, observation *model.Observation) (*model.Observation, error)
	Update(tx *gorm.DB, observation *model.Observation) (*model.Observation, error)
	Delete(tx *gorm.DB, id uint) error
}

type observationRepository struct {
	db *gorm.DB
}

// NewObservationRepository は新しいリポジトリを生成するのだ
func NewObservationRepository(db *gorm.DB) ObservationRepository {
	return &observationRepository{db: db}
}

// FindByID はIDで観察情報を1件取得するのだ
func (r *observationRepository) FindByID(id uint) (*model.Observation, error) {
	var observation model.Observation
	if err := r.db.First(&observation, id).Error; err != nil {
		return nil, err
	}
	return &observation, nil
}

// FindAll は全ての観察情報を取得するのだ
func (r *observationRepository) FindAll() ([]model.Observation, error) {
	var observations []model.Observation
	if err := r.db.Find(&observations).Error; err != nil {
		return nil, err
	}
	return observations, nil
}

// Create は新しい観察情報を作成するのだ
func (r *observationRepository) Create(tx *gorm.DB, observation *model.Observation) (*model.Observation, error) {
	if err := tx.Create(observation).Error; err != nil {
		return nil, err
	}
	return observation, nil
}

// Update は観察情報を更新するのだ
func (r *observationRepository) Update(tx *gorm.DB, observation *model.Observation) (*model.Observation, error) {
	if err := tx.Save(observation).Error; err != nil {
		return nil, err
	}
	return observation, nil
}

// Delete はIDを元に観察情報を削除するのだ
func (r *observationRepository) Delete(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&model.Observation{}, id).Error; err != nil {
		return err
	}
	return nil
}
