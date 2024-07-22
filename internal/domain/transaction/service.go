package transaction

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

func (s *Service) FindAll() []*Transaction {
	var transactions []*Transaction
	s.DB.Preload("TransactionItems").Order("created_at DESC").Find(&transactions)
	return transactions
}

func (s *Service) FindByID(id string) (*Transaction, error) {
	var transaction *Transaction
	result := s.DB.Preload("ItemActions").Preload("ItemImages").First(&transaction, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return transaction, nil
}

func (s *Service) FindByCode(code string) (*Transaction, error) {
	var transaction *Transaction
	result := s.DB.Preload("TransactionItems").First(&transaction, "code = ?", code)
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

	result := s.DB.Create(&transaction)
	if result.Error != nil {
		return nil, result.Error
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

	return transactionItem, nil
}

func (s *Service) CreateWithItems(transaction *Transaction, transactionItems []*TransactionItem) (*Transaction, error) {
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		transaction, err := s.Create(transaction)
		if err != nil {
			return err
		}

		for _, transactionItem := range transactionItems {
			transactionItem.TransactionID = transaction.ID.String()
			_, err := s.AddItem(transactionItem)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	transaction, err = s.FindByCode(transaction.Code)
	if err != nil {
		return nil, err
	}

	return transaction, nil
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
