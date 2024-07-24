package user

import (
	"github.com/go-playground/validator/v10"
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

	user.Password = hashPassword

	result := s.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
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
