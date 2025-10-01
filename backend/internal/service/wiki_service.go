
// backend/internal/service/wiki_service.go
package service

import (
	"github.com/saku-730/specimen-web/backend/internal/model"
	"github.com/saku-730/specimen-web/backend/internal/repository"
	"gorm.io/gorm"
)

// CreateWikiPageRequest はWikiページ作成時のリクエストボディを表すのだ
type CreateWikiPageRequest struct {
	Title       string `json:"title"`
	UserID      uint   `json:"user_id"`
	ContentPath string `json:"content_path"`
}

// UpdateWikiPageRequest はWikiページ更新時のリクエストボディを表すのだ
type UpdateWikiPageRequest struct {
	Title       string `json:"title"`
	UserID      uint   `json:"user_id"`
	ContentPath string `json:"content_path"`
}

// WikiService はWikiページ関連のビジネスロジックのインターフェースなのだ
type WikiService interface {
	GetWikiPageByID(id uint) (*model.WikiPage, error)
	GetAllWikiPages() ([]model.WikiPage, error)
	CreateWikiPage(req CreateWikiPageRequest) (*model.WikiPage, error)
	UpdateWikiPage(id uint, req UpdateWikiPageRequest) (*model.WikiPage, error)
	DeleteWikiPage(id uint) error
}

type wikiService struct {
	db   *gorm.DB
	repo repository.WikiRepository
}

// NewWikiService は新しいサービスを生成するのだ
func NewWikiService(db *gorm.DB, repo repository.WikiRepository) WikiService {
	return &wikiService{db: db, repo: repo}
}

func (s *wikiService) GetWikiPageByID(id uint) (*model.WikiPage, error) {
	return s.repo.FindByID(id)
}

func (s *wikiService) GetAllWikiPages() ([]model.WikiPage, error) {
	return s.repo.FindAll()
}

func (s *wikiService) CreateWikiPage(req CreateWikiPageRequest) (*model.WikiPage, error) {
	newPage := &model.WikiPage{
		Title:       req.Title,
		UserID:      req.UserID,
		ContentPath: req.ContentPath,
	}

	var createdPage *model.WikiPage
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		createdPage, err = s.repo.Create(tx, newPage)
		return err
	})

	if err != nil {
		return nil, err
	}
	return createdPage, nil
}

func (s *wikiService) UpdateWikiPage(id uint, req UpdateWikiPageRequest) (*model.WikiPage, error) {
	var updatedPage *model.WikiPage
	err := s.db.Transaction(func(tx *gorm.DB) error {
		target, err := s.repo.FindByID(id)
		if err != nil {
			return err
		}

		target.Title = req.Title
		target.UserID = req.UserID
		target.ContentPath = req.ContentPath

		updatedPage, err = s.repo.Update(tx, target)
		return err
	})

	if err != nil {
		return nil, err
	}
	return updatedPage, nil
}

func (s *wikiService) DeleteWikiPage(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Delete(tx, id)
	})
}
