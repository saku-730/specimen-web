package service

import (
	"errors"
	"github.com/saku-730/specimen-web/backend/internal/model"
	"github.com/saku-730/specimen-web/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-super-secret-key")

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
	Login(email, password string) (string, error)
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

//create new user
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

//Login method
func (s *userService) Login(email, password string) (string, error) {
	// 1. メールアドレスでユーザーを探す
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("メールアドレスまたはパスワードが違います")
	}

	// 2. パスワードを比較する
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// パスワードが一致しない場合
		return "", errors.New("メールアドレスまたはパスワードが違います")
	}

	// 3. パスワードが一致したら、JWTトークンを生成する
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // トークンの有効期限 単位h
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
