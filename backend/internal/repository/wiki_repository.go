
// backend/internal/repository/wiki_repository.go
package repository

import (
	"github.com/saku-730/specimen-web/backend/internal/model"
	"gorm.io/gorm"
)

// WikiRepository はWikiページ関連のデータ操作の契約書なのだ
type WikiRepository interface {
	FindByID(id uint) (*model.WikiPage, error)
	FindAll() ([]model.WikiPage, error)
	Create(tx *gorm.DB, page *model.WikiPage) (*model.WikiPage, error)
	Update(tx *gorm.DB, page *model.WikiPage) (*model.WikiPage, error)
	Delete(tx *gorm.DB, id uint) error
}

type wikiRepository struct {
	db *gorm.DB
}

// NewWikiRepository は新しいリポジトリを生成するのだ
func NewWikiRepository(db *gorm.DB) WikiRepository {
	return &wikiRepository{db: db}
}

// FindByID はIDでWikiページを1件取得するのだ
func (r *wikiRepository) FindByID(id uint) (*model.WikiPage, error) {
	var page model.WikiPage
	if err := r.db.First(&page, id).Error; err != nil {
		return nil, err
	}
	return &page, nil
}

// FindAll は全てのWikiページを取得するのだ
func (r *wikiRepository) FindAll() ([]model.WikiPage, error) {
	var pages []model.WikiPage
	if err := r.db.Find(&pages).Error; err != nil {
		return nil, err
	}
	return pages, nil
}

// Create は新しいWikiページを作成するのだ
func (r *wikiRepository) Create(tx *gorm.DB, page *model.WikiPage) (*model.WikiPage, error) {
	if err := tx.Create(page).Error; err != nil {
		return nil, err
	}
	return page, nil
}

// Update はWikiページを更新するのだ
func (r *wikiRepository) Update(tx *gorm.DB, page *model.WikiPage) (*model.WikiPage, error) {
	if err := tx.Save(page).Error; err != nil {
		return nil, err
	}
	return page, nil
}

// Delete はIDを元にWikiページを削除するのだ
func (r *wikiRepository) Delete(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&model.WikiPage{}, id).Error; err != nil {
		return err
	}
	return nil
}
