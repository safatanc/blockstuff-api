package auth

import (
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/safatanc/blockstuff-api/internal/domain/user"
	"github.com/safatanc/blockstuff-api/pkg/converter"
	"gorm.io/gorm"
)

type Service struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewService(db *gorm.DB, validate *validator.Validate) *Service {
	return &Service{
		DB:       db,
		Validate: validate,
	}
}

func (s *Service) NewToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (s *Service) VerifyUser(username string, password string) (*Auth, error) {
	var u *user.User
	result := s.DB.First(&u, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}
	err := converter.VerifyPassword(password, u.Password)
	if err != nil {
		return nil, err
	}

	tokenString, err := s.NewToken(u.Username)
	if err != nil {
		return nil, err
	}
	return &Auth{
		Token: tokenString,
	}, nil
}
