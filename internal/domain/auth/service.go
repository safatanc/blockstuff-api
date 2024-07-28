package auth

import (
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/safatanc/blockstuff-api/internal/domain/user"
	"github.com/safatanc/blockstuff-api/pkg/converter"
	"github.com/safatanc/blockstuff-api/pkg/jwthelper"
	"gorm.io/gorm"
)

type Service struct {
	DB          *gorm.DB
	Validate    *validator.Validate
	UserService *user.Service
}

func NewService(db *gorm.DB, validate *validator.Validate, userService *user.Service) *Service {
	return &Service{
		DB:          db,
		Validate:    validate,
		UserService: userService,
	}
}

func (s *Service) NewToken(user *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.ToResponse(),
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})

	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (s *Service) VerifyToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwthelper.VerifyToken(tokenString)
	claims := token.Claims.(jwt.MapClaims)
	return &claims, err
}

func (s *Service) VerifyUser(username string, password string) (*Auth, error) {
	var user *user.User
	result := s.DB.First(&user, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}
	err := converter.VerifyPassword(password, user.Password)
	if err != nil {
		return nil, err
	}

	tokenString, err := s.NewToken(user)
	if err != nil {
		return nil, err
	}
	return &Auth{
		Token: tokenString,
	}, nil
}

func (s *Service) Register(user *user.User) (*user.User, error) {
	user, err := s.UserService.Create(user)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}
