package payout

import (
	"github.com/go-playground/validator/v10"
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

func (s *Service) FindAll(status string) []*Payout {
	var payouts []*Payout
	s.DB.Preload("PayoutTransactions").Order("created_at DESC").Find(&payouts, "status = ?", status)
	return payouts
}

func (s *Service) FindByID(id string) (*Payout, error) {
	var payout *Payout
	result := s.DB.Preload("PayoutTransactions").First(&payout, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return payout, nil
}

func (s *Service) Create(payout *Payout) (*Payout, error) {
	err := s.Validate.Struct(payout)
	if err != nil {
		return nil, err
	}

	result := s.DB.Create(&payout)
	if result.Error != nil {
		return nil, result.Error
	}

	return payout, nil
}

func (s *Service) Update(id string, payout *Payout) (*Payout, error) {
	err := s.Validate.Struct(payout)
	if err != nil {
		return nil, err
	}

	result := s.DB.Where("id = ?", id).Updates(&payout)
	if result.Error != nil {
		return nil, result.Error
	}
	return payout, nil
}

func (s *Service) Delete(id string) (*Payout, error) {
	var payout *Payout
	result := s.DB.First(&payout, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	result = s.DB.Where("id = ?", id).Delete(&payout)
	if result.Error != nil {
		return nil, result.Error
	}
	return payout, nil
}
