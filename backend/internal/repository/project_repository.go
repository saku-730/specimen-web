// internal/repository/project_repository.go
package repository

import (
	"github.com/saku-730/specimen-web/backend/internal/model"

	"gorm.io/gorm"
)

// ProjectRepository はプロジェクト関連のデータ操作の契約書なのだ
type ProjectRepository interface {
	FindByID(id uint) (*model.Project, error)
	FindAll() ([]model.Project, error)
	Create(project *model.Project) (*model.Project, error)
	AddMember(member *model.ProjectMember) (*model.ProjectMember, error)
}

type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository は新しいリポジトリを生成するのだ
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

// FindByID はIDでプロジェクトを1件取得する。メンバー情報も一緒に読み込むのだ
func (r *projectRepository) FindByID(id uint) (*model.Project, error) {
	var project model.Project
	// Preloadを使うと、関連するProjectMembersと、さらにその中のUser情報も一緒に取得できるのだ
	if err := r.db.Preload("ProjectMembers.User").First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// FindAll は全てのプロジェクトを取得するのだ
func (r *projectRepository) FindAll() ([]model.Project, error) {
	var projects []model.Project
	if err := r.db.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// Create は新しいプロジェクトを作成するのだ
func (r *projectRepository) Create(project *model.Project) (*model.Project, error) {
	if err := r.db.Create(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

// AddMember はプロジェクトに新しいメンバーを追加するのだ
func (r *projectRepository) AddMember(member *model.ProjectMember) (*model.ProjectMember, error) {
	if err := r.db.Create(member).Error; err != nil {
		return nil, err
	}
	return member, nil
}
