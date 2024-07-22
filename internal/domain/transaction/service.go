package transaction

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/safatanc/blockstuff-api/pkg/util"
	"gorm.io/gorm"
)

type Service struct {
	DB           *gorm.DB
	Validate     *validator.Validate
	MidtransCore *coreapi.Client
}

func NewService(db *gorm.DB, validate *validator.Validate, midtransCore *coreapi.Client) *Service {
	return &Service{
		DB:           db,
		Validate:     validate,
		MidtransCore: midtransCore,
	}
}

func (s *Service) FindAll() []*Transaction {
	var transactions []*Transaction
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
				return fmt.Errorf("request limit reached. cooldown %.2f seconds", math.Abs(difference.Seconds()))
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

		chargeResponse, err := s.CreatePayment(transaction)
		if err != nil {
			return err
		}
		transaction.QrisString = chargeResponse.QRString

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

func (s *Service) CreatePayment(transaction *Transaction) (*coreapi.ChargeResponse, error) {
	var midtransItems []midtrans.ItemDetails
	for _, transactionItem := range transaction.TransactionItems {
		midtransItems = append(midtransItems, midtrans.ItemDetails{
			ID:           transactionItem.ItemID,
			Name:         transactionItem.Item.Name,
			Price:        transactionItem.Item.Price,
			Qty:          int32(transactionItem.Quantity),
			Brand:        *transactionItem.Item.MinecraftServerID,
			Category:     transactionItem.Item.Category,
			MerchantName: "BLOCKSTUFF",
		})
	}

	payload := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeQris,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transaction.Code,
			GrossAmt: transaction.Subtotal,
		},
		Gopay: &coreapi.GopayDetails{
			EnableCallback: true,
			CallbackUrl:    os.Getenv("CALLBACK_URL"),
		},
		Items: &midtransItems,
		CustomerDetails: &midtrans.CustomerDetails{
			FName: transaction.MinecraftUsername,
			Email: transaction.Email,
			Phone: fmt.Sprintf("%v", transaction.Phone),
		},
	}

	chargeResponse, err := s.MidtransCore.ChargeTransaction(payload)
	if err != nil {
		return nil, err
	}
	return chargeResponse, nil
}
