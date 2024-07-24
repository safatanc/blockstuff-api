package payout

import (
	"os"
	"strconv"

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
	var payouts = make([]*Payout, 0)
	if status != "" {
		s.DB.Preload("PayoutTransactions").Order("created_at DESC").Find(&payouts, "status = ?", status)
	} else {
		s.DB.Preload("PayoutTransactions").Order("created_at DESC").Find(&payouts)
	}
	return payouts
}

func (s *Service) FindByID(id string) (*Payout, error) {
	var payout *Payout
	result := s.DB.Preload("PayoutTransactions.Transaction.TransactionItems.Item").First(&payout, "id = ?", id)
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

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&payout)
		if result.Error != nil {
			return result.Error
		}
		result = tx.Preload("PayoutTransactions.Transaction.TransactionItems.Item").First(&payout, "id = ?", payout.ID)
		if result.Error != nil {
			return result.Error
		}

		for _, payoutTransaction := range payout.PayoutTransactions {
			for _, transactionItem := range payoutTransaction.Transaction.TransactionItems {
				payout.TransactionSubtotal += transactionItem.Item.Price
			}
		}

		payoutFeePercent, err := strconv.Atoi(os.Getenv("PAYOUT_FEE_PERCENT"))
		if err != nil {
			return err
		}

		payout.Fee = int64((float64(payoutFeePercent) / float64(100)) * float64(payout.TransactionSubtotal))
		payout.Subtotal = payout.TransactionSubtotal - payout.Fee

		result = tx.Where("id = ?", payout.ID).Updates(&payout)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		return nil, err
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
