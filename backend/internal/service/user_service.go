package service

import (
	"errors"
	"github.com/saku-730/specimen-web/backend/internal/model"
	"github.com/saku-730/specimen-web/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUserRequest struct {
	UserName    string `json:"user_name"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
	RoleID      uint   `json:"role_id"`
}

type UserService interface {
	GetUserByID(id uint) (*model.User, error)
	GetAllUsers() ([]model.User, error)
	CreateUser(req CreateUserRequest) (*model.User, error)
}

type userService struct {
	db       *gorm.DB // トランザクション用にdb接続を持つ
	userRepo repository.UserRepository
}

func NewUserService(db *gorm.DB, userRepo repository.UserRepository) UserService {
	return &userService{db: db, userRepo: userRepo}
}

func (s *userService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.userRepo.FindAll()
}

func (s *userService) CreateUser(req CreateUserRequest) (*model.User, error) {
	if req.UserName == "" || req.Password == "" {
		return nil, errors.New("ユーザー名とパスワードは必須です")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		UserName:    req.UserName,
		DisplayName: req.DisplayName,
		Password:    string(hashedPassword),
		RoleID:      req.RoleID,
	}

	var createdUser *model.User
	err = s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		createdUser, err = s.userRepo.Create(tx, newUser)
		if err != nil {
			return err
		}
		// 他にもユーザー作成時にやるべきことがあればここに追加
		return nil
	})

	if err != nil {
		return nil, err
	}
	return createdUser, nil
}
