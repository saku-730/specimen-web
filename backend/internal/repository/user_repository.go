// internal/repository/user_repository.go
package repository

import (
	"github.com/saku-730/specimen-web/backend/internal/model"

	"gorm.io/gorm"
)

// 1. UserRepository は、ユーザーデータに関する操作の「契約書」(インターフェース)なのだ
type UserRepository interface {
	FindByID(id uint) (*model.User, error)
	FindAll() ([]model.User, error)
	Create(tx *gorm.DB, user *model.User) (*model.User, error)
	Update(tx *gorm.DB, user *model.User) (*model.User, error)
	Delete(tx *gorm.DB, id uint) error
	FindByEmail(email string)(*model.User, error)
}

// 2. userRepository は、UserRepositoryインターフェースの「実装」なのだ
type userRepository struct {
	db *gorm.DB
}

// 3. NewUserRepository は、新しいuserRepositoryを生成するための「コンストラクタ」なのだ
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// 4. FindByID は、IDを元にユーザーを1件取得するメソッドなのだ
func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	// GORMのFirstメソッドを使って、主キーで検索するのだ
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 5. FindAll は、全てのユーザーを取得するメソッドなのだ
func (r *userRepository) FindAll() ([]model.User, error) {
	var users []model.User
	// GORMのFindメソッドを使って、全件検索するのだ
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}



// 6. Create は、新しいユーザーをデータベースに作成するメソッドなのだ
func (r *userRepository) Create(tx *gorm.DB, user *model.User) (*model.User, error) {
	// GORMのCreateメソッドを使って、レコードをINSERTするのだ
	if err := tx.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// 7. Update は、ユーザー情報を更新するメソッドなのだ
func (r *userRepository) Update(tx *gorm.DB, user *model.User) (*model.User, error) {
	// GORMのSaveメソッドを使って、レコードをUPDATEするのだ
	// Saveは全フィールドを更新する。一部だけ更新したい場合はUpdateを使うのだ。
	if err := tx.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// 8. Delete は、IDを元にユーザーを削除するメソッドなのだ
func (r *userRepository) Delete(tx *gorm.DB, id uint) error {
	// GORMのDeleteメソッドを使って、レコードをDELETEするのだ
	if err := tx.Delete(&model.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

// Find user by email address
func (r *userRepository) FindByEmail(email string)(*model.User, error){
	var user model.User
	if err := r.db.Where("mail_address = ?", email).First(&user).Error; err != nil{
		return nil, err
	}
	return &user, nil
}

