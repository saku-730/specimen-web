// internal/model/user_models.go
package model

import (
	"time"
)

// UserRole は "user_roles" テーブルに対応するのだ
type UserRole struct {
	RoleID   uint   `gorm:"primaryKey" json:"role_id"`
	RoleName string `gorm:"not null;unique" json:"role_name"`
}

// User は "users" テーブルに対応するのだ
type User struct {
	UserID       uint       `gorm:"primaryKey" json:"user_id"`
	UserName     string     `gorm:"size:255;not null" json:"user_name"`
	DisplayName  string     `gorm:"size:255;not null" json:"display_name"`
	MailAddress  string     `gorm:"size:255;unique" json:"mail_address"`
	Password     string     `json:"-"` // JSONには含めないようにするのだ
	RoleID       uint       `json:"role_id"`
	CreatedAt    time.Time  `gorm:"default:now()" json:"created_at"`
	Timezone     int16      `gorm:"not null" json:"timezone"`

	// 関連 (Associations)
	Role         UserRole      `gorm:"foreignKey:RoleID" json:"role"`
	UserDefault  UserDefault   `gorm:"foreignKey:UserID" json:"user_default"`
	WikiPages    []WikiPage    `gorm:"foreignKey:UserID" json:"wiki_pages"`
	ProjectMembers []ProjectMember `gorm:"foreignKey:UserID" json:"project_members"`
}

type UserDefault struct {
	UserID uint   `gorm:"primaryKey" json:"user_id"`
	Theme  string `gorm:"not null;default:'white'" json:"theme"`
}
