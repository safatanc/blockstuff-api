package item

import (
	"fmt"
	"mime/multipart"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/safatanc/blockstuff-api/internal/domain/storage"
	"github.com/safatanc/blockstuff-api/pkg/util"
	"gorm.io/gorm"
)

type Service struct {
	DB       *gorm.DB
	Validate *validator.Validate
	Storage  *storage.Service
}

func NewService(db *gorm.DB, validate *validator.Validate, storage *storage.Service) *Service {
	return &Service{
		DB:       db,
		Validate: validate,
		Storage:  storage,
	}
}

func (s *Service) FindAll(minecraftServerID string) []*Item {
	var items = make([]*Item, 0)
	s.DB.Preload("ItemImages").Order("price ASC").Find(&items, "minecraft_server_id = ?", minecraftServerID)
	return items
}

func (s *Service) FindByID(id string) (*Item, error) {
	var item *Item
	result := s.DB.Preload("ItemActions").Preload("ItemImages").Preload("MinecraftServer").First(&item, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return item, nil
}

func (s *Service) FindBySlug(minecraftServerID string, slug string) (*Item, error) {
	var item *Item
	result := s.DB.Preload("ItemActions").Preload("ItemImages").First(&item, "minecraft_server_id = ? AND slug = ?", minecraftServerID, slug)
	if result.Error != nil {
		return nil, result.Error
	}
	return item, nil
}

func (s *Service) Create(item *Item) (*Item, error) {
	err := s.Validate.Struct(item)
	if err != nil {
		return nil, err
	}

	_, err = s.FindBySlug(*item.MinecraftServerID, item.Slug)
	if err == nil {
		return nil, fmt.Errorf("slug %v already exists", item.Slug)
	}

	result := s.DB.Create(&item)
	if result.Error != nil {
		return nil, result.Error
	}

	return item, nil
}

func (s *Service) AddImage(itemID string, image multipart.File, fileHeader *multipart.FileHeader) (*ItemImage, error) {
	objectName := fmt.Sprintf("%s--%s.png", itemID, util.RandomString(12))
	_, err := s.Storage.Upload(objectName, image, "image/*")
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/storage/%s", os.Getenv("BASE_URL"), objectName)

	itemImage := &ItemImage{
		URL:    url,
		ItemID: itemID,
	}

	result := s.DB.Create(&itemImage)
	if result.Error != nil {
		return nil, result.Error
	}

	return itemImage, nil
}

func (s *Service) AddAction(itemAction *ItemAction) (*ItemAction, error) {
	err := s.Validate.Struct(itemAction)
	if err != nil {
		return nil, err
	}

	result := s.DB.Create(&itemAction)
	if result.Error != nil {
		return nil, result.Error
	}

	return itemAction, nil
}

func (s *Service) Update(id string, item *Item) (*Item, error) {
	err := s.Validate.Struct(item)
	if err != nil {
		return nil, err
	}

	result := s.DB.Where("id = ?", id).Updates(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return item, nil
}

func (s *Service) Delete(id string) (*Item, error) {
	var item *Item
	result := s.DB.First(&item, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	result = s.DB.Where("id = ?", id).Delete(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return item, nil
}
