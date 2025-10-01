// internal/model/project_models.go
package model

import "time"

type Project struct {
	ProjectID    uint       `gorm:"primaryKey" json:"project_id"`
	ProjectName  string     `gorm:"not null" json:"project_name"`
	Description  string     `json:"description"` // SQLのdisscriptionはtypoだと思うので修正
	StartDay     *time.Time `json:"start_day"` // NULLを許容する日付はポインタ型にするのだ
	FinishedDay  *time.Time `json:"finished_day"`
	UpdatedDay   *time.Time `json:"updated_day"`
	Note         string     `json:"note"`
}

type ProjectMember struct {
	ProjectMemberID uint       `gorm:"primaryKey" json:"project_member_id"`
	ProjectID       uint       `json:"project_id"`
	UserID          uint       `json:"user_id"`
	JoinDay         *time.Time `json:"join_day"`
	FinishDay       *time.Time `json:"finish_day"`
	
	// 関連
	Project Project `gorm:"foreignKey:ProjectID" json:"project"`
	User    User    `gorm:"foreignKey:UserID" json:"user"`
}
