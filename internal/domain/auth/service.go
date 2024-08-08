package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/safatanc/blockstuff-api/internal/domain/mail"
	"github.com/safatanc/blockstuff-api/internal/domain/user"
	"github.com/safatanc/blockstuff-api/pkg/converter"
	"github.com/safatanc/blockstuff-api/pkg/jwthelper"
	"github.com/safatanc/blockstuff-api/pkg/util"
	"gorm.io/gorm"
)

type Service struct {
	DB          *gorm.DB
	Validate    *validator.Validate
	UserService *user.Service
	MailService *mail.Service
}

func NewService(db *gorm.DB, validate *validator.Validate, userService *user.Service, mailService *mail.Service) *Service {
	return &Service{
		DB:          db,
		Validate:    validate,
		UserService: userService,
		MailService: mailService,
	}
}

func (s *Service) NewToken(user *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":     user.ToResponse(),
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
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

	if !user.EmailVerified {
		return nil, fmt.Errorf("email is not verified")
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
	emailVerifyCode := util.RandomString(5)
	user.EmailVerifyCode = &emailVerifyCode

	user, err := s.UserService.Create(user)
	if err != nil {
		return nil, err
	}

	err = s.MailService.Send([]string{user.Email}, "Verify Email", fmt.Sprintf("Kode Verifikasi: %v", *user.EmailVerifyCode))
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

func (s *Service) VerifyEmail(user *user.User, code string) error {
	if code == "" {
		return fmt.Errorf("invalid code")
	}

	if *user.EmailVerifyCode != code {
		return fmt.Errorf("invalid code")
	}
	return nil
}

func (s *Service) SendVerifyCode(email string, subject string) error {
	user, err := s.UserService.FindByEmail(email)
	if err != nil {
		return err
	}

	emailVerifyCode := util.RandomString(5)
	user.EmailVerifyCode = &emailVerifyCode

	_, err = s.UserService.Update(user.ID.String(), user)
	if err != nil {
		return err
	}

	err = s.MailService.Send([]string{user.Email}, subject, fmt.Sprintf("Kode Verifikasi: %v", *user.EmailVerifyCode))
	return err
}

func (s *Service) ResetPasswordVerify(email string, code string, newPassword string) error {
	validatePassword := user.User{
		Password: newPassword,
	}
	err := s.Validate.StructPartial(validatePassword, "Password")
	if err != nil {
		return err
	}

	user, err := s.UserService.FindByEmail(email)
	if err != nil {
		return err
	}

	if user.EmailVerifyCode == nil {
		return fmt.Errorf("invalid code")
	}

	if *user.EmailVerifyCode != code {
		return fmt.Errorf("invalid code")
	}

	hashPassword, err := converter.PasswordToHash(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashPassword
	user.EmailVerifyCode = nil

	_, err = s.UserService.Update(user.ID.String(), user)
	return err
}
