package user

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/safatanc/blockstuff-api/internal/domain/mail"
	"github.com/safatanc/blockstuff-api/pkg/converter"
	"github.com/safatanc/blockstuff-api/pkg/util"
	"gorm.io/gorm"
)

type Service struct {
	DB          *gorm.DB
	Validate    *validator.Validate
	MailService *mail.Service
}

func NewService(db *gorm.DB, validate *validator.Validate, mailService *mail.Service) *Service {
	return &Service{
		DB:          db,
		Validate:    validate,
		MailService: mailService,
	}
}

func (s *Service) FindAll() []*User {
	var users = make([]*User, 0)
	s.DB.Find(&users)

	var userResponses []*User
	for _, u := range users {
		userResponses = append(userResponses, u.ToResponse())
	}

	return userResponses
}

func (s *Service) FindByID(id string) (*User, error) {
	var user *User
	result := s.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.ToResponse(), nil
}

func (s *Service) FindByUsername(username string) (*User, error) {
	var user *User
	result := s.DB.First(&user, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.ToResponse(), nil
}

func (s *Service) Create(user *User) (*User, error) {
	err := s.Validate.Struct(user)
	if err != nil {
		return nil, err
	}

	hashPassword, err := converter.PasswordToHash(user.Password)
	if err != nil {
		return nil, err
	}

	emailVerifyCode := util.RandomString(5)

	user.Password = hashPassword
	user.EmailVerifyCode = &emailVerifyCode

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&user)
		if result.Error != nil {
			return result.Error
		}

		err := s.MailService.Send([]string{user.Email}, "Verify Email", fmt.Sprintf("Kode Verifikasi: %v", *user.EmailVerifyCode))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *Service) Update(id string, user *User) (*User, error) {
	err := s.Validate.Struct(user)
	if err != nil {
		return nil, err
	}

	result := s.DB.Where("id = ?", id).Updates(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.ToResponse(), nil
}

func (s *Service) Delete(id string) (*User, error) {
	var user *User
	result := s.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	result = s.DB.Where("id = ?", id).Delete(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user.ToResponse(), nil
}
