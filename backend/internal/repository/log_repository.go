
// backend/internal/repository/log_repository.go
package repository

import (
	"github.com/saku-730/specimen-web/backend/internal/model"
	"gorm.io/gorm"
)

// LogRepository は変更ログ関連のデータ操作の契約書なのだ
type LogRepository interface {
	FindByID(id uint) (*model.ChangeLog, error)
	FindAll() ([]model.ChangeLog, error)
	Create(tx *gorm.DB, log *model.ChangeLog) (*model.ChangeLog, error)
	Update(tx *gorm.DB, log *model.ChangeLog) (*model.ChangeLog, error)
	Delete(tx *gorm.DB, id uint) error
}

type logRepository struct {
	db *gorm.DB
}

// NewLogRepository は新しいリポジトリを生成するのだ
func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

// FindByID はIDで変更ログを1件取得するのだ
func (r *logRepository) FindByID(id uint) (*model.ChangeLog, error) {
	var log model.ChangeLog
	if err := r.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// FindAll は全ての変更ログを取得するのだ
func (r *logRepository) FindAll() ([]model.ChangeLog, error) {
	var logs []model.ChangeLog
	if err := r.db.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// Create は新しい変更ログを作成するのだ
func (r *logRepository) Create(tx *gorm.DB, log *model.ChangeLog) (*model.ChangeLog, error) {
	if err := tx.Create(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

// Update は変更ログを更新するのだ
func (r *logRepository) Update(tx *gorm.DB, log *model.ChangeLog) (*model.ChangeLog, error) {
	if err := tx.Save(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

// Delete はIDを元に変更ログを削除するのだ
func (r *logRepository) Delete(tx *gorm.DB, id uint) error {
	if err := tx.Delete(&model.ChangeLog{}, id).Error; err != nil {
		return err
	}
	return nil
}
