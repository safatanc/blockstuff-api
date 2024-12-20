package payout

import (
	"context"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/safatanc/blockstuff-api/internal/domain/user"
	"github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/payout"
	"gorm.io/gorm"
)

type Service struct {
	DB           *gorm.DB
	Validate     *validator.Validate
	XenditClient *xendit.APIClient
}

func NewService(db *gorm.DB, validate *validator.Validate, xenditClient *xendit.APIClient) *Service {
	return &Service{
		DB:           db,
		Validate:     validate,
		XenditClient: xenditClient,
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

func (s *Service) FindPayoutChannels() ([]payout.Channel, error) {
	payoutChannels, _, err := s.XenditClient.PayoutApi.GetPayoutChannels(context.Background()).Currency("IDR").Execute()
	if err != nil {
		return nil, err
	}
	return payoutChannels, nil
}

func (s *Service) GetPayoutChannel(userID string) (*user.UserPayoutChannel, error) {
	var userPayoutChannel *user.UserPayoutChannel
	err := s.DB.Take(&userPayoutChannel, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	return userPayoutChannel, nil
}

func (s *Service) SetPayoutChannel(userID string, userPayoutChannel *user.UserPayoutChannel) (*user.UserPayoutChannel, error) {
	userPayoutChannel.UserID = userID
	err := s.Validate.Struct(userPayoutChannel)
	if err != nil {
		return nil, err
	}

	findUserPayoutChannel, err := s.GetPayoutChannel(userID)
	if err == nil {
		userPayoutChannel = findUserPayoutChannel
	}

	err = s.DB.Save(&userPayoutChannel).Error
	return userPayoutChannel, err
}
