package transaction

import (
	"context"
	"fmt"
	"math"
	"time"

	_ "image/jpeg"
	_ "image/png"

	"github.com/go-playground/validator/v10"
	"github.com/safatanc/blockstuff-api/pkg/util"
	"github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/payment_request"
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

func (s *Service) FindAll() []*Transaction {
	var transactions = make([]*Transaction, 0)
	s.DB.Preload("TransactionItems").Preload("TransactionItems.Item").Order("created_at DESC").Find(&transactions)
	return transactions
}

func (s *Service) FindByID(id string) (*Transaction, error) {
	var transaction *Transaction
	result := s.DB.Preload("TransactionItems").Preload("TransactionItems.Item").First(&transaction, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

func (s *Service) FindByCode(code string) (*Transaction, error) {
	var transaction *Transaction
	result := s.DB.Preload("TransactionItems").Preload("TransactionItems.Item").First(&transaction, "code = ?", code)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

func (s *Service) FindByCodeTx(tx *gorm.DB, code string) (*Transaction, error) {
	var transaction *Transaction
	result := tx.Preload("TransactionItems").Preload("TransactionItems.Item").First(&transaction, "code = ?", code)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

func (s *Service) Create(transaction *Transaction) (*Transaction, error) {
	err := s.Validate.Struct(transaction)
	if err != nil {
		return nil, err
	}

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		var findTransaction *Transaction
		result := s.DB.Order("created_at DESC").First(&findTransaction, "minecraft_username = ?", transaction.MinecraftUsername)
		if result.Error == nil {
			difference := time.Until(findTransaction.CreatedAt)
			if math.Abs(difference.Seconds()) < 60 {
				return fmt.Errorf("request limit reached. cooldown %.2f seconds", 60-math.Abs(difference.Seconds()))
			}
		}

		transactionCode := fmt.Sprintf("BS-%v", util.RandomString(10))
		transaction.Code = transactionCode

		result = tx.Create(&transaction)
		if result.Error != nil {
			return result.Error
		}

		transaction, err = s.FindByCodeTx(tx, transaction.Code)
		if err != nil {
			return err
		}

		for _, transactionItem := range transaction.TransactionItems {
			transactionItem.Subtotal = transactionItem.Item.Price * int64(transactionItem.Quantity)
			transactionItem, err := s.UpdateItemTx(tx, transactionItem.ID.String(), transactionItem)
			if err != nil {
				return err
			}
			transaction.Subtotal += transactionItem.Subtotal
		}

		paymentResponse, err := s.CreatePayment(transaction)
		if err != nil {
			return err
		}

		transaction.PaymentID = &paymentResponse.Id
		transaction.QrisString = paymentResponse.PaymentMethod.GetQrCode().ChannelProperties.GetQrString()

		transaction, err = s.UpdateTx(tx, transaction.ID.String(), transaction)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *Service) AddItem(transactionItem *TransactionItem) (*TransactionItem, error) {
	err := s.Validate.Struct(transactionItem)
	if err != nil {
		return nil, err
	}

	result := s.DB.Create(&transactionItem)
	if result.Error != nil {
		return nil, result.Error
	}

	result = s.DB.Preload("TransactionItems.Item").First(&transactionItem, "id = ?", transactionItem.ID)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactionItem, nil
}

func (s *Service) Update(id string, transaction *Transaction) (*Transaction, error) {
	err := s.Validate.Struct(transaction)
	if err != nil {
		return nil, err
	}

	result := s.DB.Where("id = ?", id).Updates(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

func (s *Service) UpdateTx(tx *gorm.DB, id string, transaction *Transaction) (*Transaction, error) {
	err := s.Validate.Struct(transaction)
	if err != nil {
		return nil, err
	}

	result := tx.Where("id = ?", id).Updates(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

func (s *Service) UpdateItem(id string, transactionItem *TransactionItem) (*TransactionItem, error) {
	err := s.Validate.Struct(transactionItem)
	if err != nil {
		return nil, err
	}

	result := s.DB.Where("id = ?", id).Updates(&transactionItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactionItem, nil
}

func (s *Service) UpdateItemTx(tx *gorm.DB, id string, transactionItem *TransactionItem) (*TransactionItem, error) {
	err := s.Validate.Struct(transactionItem)
	if err != nil {
		return nil, err
	}

	result := tx.Where("id = ?", id).Updates(&transactionItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return transactionItem, nil
}

func (s *Service) Delete(id string) (*Transaction, error) {
	var transaction *Transaction
	result := s.DB.First(&transaction, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	result = s.DB.Where("id = ?", id).Delete(&transaction)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

func (s *Service) CreatePayment(transaction *Transaction) (*payment_request.PaymentRequest, error) {
	var xenditItems []payment_request.PaymentRequestBasketItem
	for _, transactionItem := range transaction.TransactionItems {
		xenditItems = append(xenditItems, payment_request.PaymentRequestBasketItem{
			ReferenceId: &transactionItem.ItemID,
			Name:        transactionItem.Item.Name,
			Price:       float64(transactionItem.Item.Price),
			Currency:    "IDR",
			Quantity:    float64(transactionItem.Quantity),
			Type:        &transactionItem.Item.Category,
		})
	}

	transactionAmount := float64(transaction.Subtotal)

	payload := payment_request.PaymentRequestParameters{
		ReferenceId: &transaction.Code,
		Amount:      &transactionAmount,
		Currency:    payment_request.PAYMENTREQUESTCURRENCY_IDR,
		PaymentMethod: payment_request.NewPaymentMethodParameters(
			payment_request.PAYMENTMETHODTYPE_QR_CODE,
			payment_request.PAYMENTMETHODREUSABILITY_ONE_TIME_USE,
		),
		Items: xenditItems,
		Metadata: map[string]interface{}{
			"minecraft_username": transaction.MinecraftUsername,
			"email":              transaction.Email,
			"phone":              fmt.Sprintf("%v", transaction.Phone),
		},
	}

	paymentResponse, _, err := s.XenditClient.PaymentRequestApi.CreatePaymentRequest(context.Background()).PaymentRequestParameters(payload).Execute()
	if err != nil {
		return nil, err
	}
	return paymentResponse, nil
}
