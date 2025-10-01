
// backend/internal/repository/place_repository.go
package repository

import (
	"github.com/saku-730/specimen-web/backend/internal/model"
	"gorm.io/gorm"
)

// PlaceRepository は場所情報関連のデータ操作の契約書なのだ
type PlaceRepository interface {
	FindByID(id uint) (*model.Place, error)
	FindAll() ([]model.Place, error)
	Create(tx *gorm.DB, place *model.Place) (*model.Place, error)
	Update(tx *gorm.DB, place *model.Place) (*model.Place, error)
	Delete(tx *gorm.DB, id uint) error
}

type placeRepository struct {
	db *gorm.DB
}

// NewPlaceRepository は新しいリポジトリを生成するのだ
func NewPlaceRepository(db *gorm.DB) PlaceRepository {
	return &placeRepository{db: db}
}

// FindByID はIDで場所情報を1件取得するのだ
func (r *placeRepository) FindByID(id uint) (*model.Place, error) {
	var place model.Place
	if err := r.db.First(&place, id).Error; err != nil {
		return nil, err
	}
	return &place, nil
}

// FindAll は全ての場所情報を取得するのだ
func (r *placeRepository) FindAll() ([]model.Place, error) {
	var places []model.Place
	if err := r.db.Find(&places).Error; err != nil {
		return nil, err
	}
	return places, nil
}

// Create は新しい場所情報を作成するのだ
func (r *placeRepository) Create(tx *gorm.DB, place *model.Place) (*model.Place, error) {
	if err := tx.Create(place).Error; err != nil {
		return nil, err
	}
	return place, nil
}

// Update は場所情報を更新するのだ
func (r *placeRepository) Update(tx *gorm.DB, place *model.Place) (*model.Place, error) {
	if err := tx.Save(place).Error; err != nil {
		return nil, err
	}
	return place, nil
}

// Delete はIDを元に場所情報を削除するのだ
func (r *placeRepository) Delete(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&model.Place{}, id).Error; err != nil {
		return err
	}
	return nil
}
