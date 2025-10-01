
// backend/internal/service/project_service.go
package service

import (
	"time"

	"github.com/saku-730/specimen-web/backend/internal/model"
	"github.com/saku-730/specimen-web/backend/internal/repository"
	"gorm.io/gorm"
)

// CreateProjectRequest はプロジェクト作成時のリクエストボディを表すのだ
type CreateProjectRequest struct {
	ProjectName string     `json:"project_name"`
	Description string     `json:"description"`
	StartDay     *time.Time `json:"start_day"`
	FinishedDay  *time.Time `json:"finished_day"`
	Note        string     `json:"note"`
}

// UpdateProjectRequest はプロジェクト更新時のリクエストボディを表すのだ
type UpdateProjectRequest struct {
	ProjectName string     `json:"project_name"`
	Description string     `json:"description"`
	StartDay     *time.Time `json:"start_day"`
	FinishedDay  *time.Time `json:"finished_day"`
	Note        string     `json:"note"`
}

// ProjectService はプロジェクト関連のビジネスロジックのインターフェースなのだ
type ProjectService interface {
	GetProjectByID(id uint) (*model.Project, error)
	GetAllProjects() ([]model.Project, error)
	CreateProject(req CreateProjectRequest) (*model.Project, error)
	UpdateProject(id uint, req UpdateProjectRequest) (*model.Project, error)
	DeleteProject(id uint) error
}

type projectService struct {
	db   *gorm.DB
	repo repository.ProjectRepository
}

// NewProjectService は新しいサービスを生成するのだ
func NewProjectService(db *gorm.DB, repo repository.ProjectRepository) ProjectService {
	return &projectService{db: db, repo: repo}
}

// GetProjectByID はIDでプロジェクトを1件取得するのだ
func (s *projectService) GetProjectByID(id uint) (*model.Project, error) {
	return s.repo.FindByID(id)
}

// GetAllProjects は全てのプロジェクトを取得するのだ
func (s *projectService) GetAllProjects() ([]model.Project, error) {
	return s.repo.FindAll()
}

// CreateProject は新しいプロジェクトを作成するのだ
func (s *projectService) CreateProject(req CreateProjectRequest) (*model.Project, error) {
	newProject := &model.Project{
		ProjectName: req.ProjectName,
		Description: req.Description,
		StartDay:    req.StartDay,
		FinishedDay: req.FinishedDay,
		Note:        req.Note,
	}

	var createdProject *model.Project
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		createdProject, err = s.repo.Create(tx, newProject)
		return err
	})

	if err != nil {
		return nil, err
	}
	return createdProject, nil
}

// UpdateProject はプロジェクトを更新するのだ
func (s *projectService) UpdateProject(id uint, req UpdateProjectRequest) (*model.Project, error) {
	var updatedProject *model.Project
	err := s.db.Transaction(func(tx *gorm.DB) error {
		target, err := s.repo.FindByID(id)
		if err != nil {
			return err
		}

		target.ProjectName = req.ProjectName
		target.Description = req.Description
		target.StartDay = req.StartDay
		target.FinishedDay = req.FinishedDay
		target.Note = req.Note

		updatedProject, err = s.repo.Update(tx, target)
		return err
	})

	if err != nil {
		return nil, err
	}
	return updatedProject, nil
}

// DeleteProject はIDを元にプロジェクトを削除するのだ
func (s *projectService) DeleteProject(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.repo.Delete(tx, id)
	})
}
