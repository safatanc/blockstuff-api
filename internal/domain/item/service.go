package item

import (
	"fmt"
	"mime/multipart"
	"strings"

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
	s.DB.Preload("ItemImages").Order("price ASC").Find(&items, "minecraft_server_id = ? AND visible = ?", minecraftServerID, true)
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
	result := s.DB.Preload("ItemActions").Preload("ItemImages").First(&item, "minecraft_server_id = ? AND slug = ? AND visible = ?", minecraftServerID, slug, true)
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

	var findItem *Item
	err = s.DB.First(&findItem, "minecraft_server_id = ? AND slug = ?", *item.MinecraftServerID, item.Slug).Error
	if err == nil {
		return nil, fmt.Errorf("slug %v already exists", item.Slug)
	}

	result := s.DB.Create(&item)
	if result.Error != nil {
		return nil, result.Error
	}

	return item, nil
}

func (s *Service) FindActionByID(itemActionID string) (*ItemAction, error) {
	var itemAction *ItemAction
	result := s.DB.First(&itemAction, "id = ?", itemActionID)
	if result.Error != nil {
		return nil, result.Error
	}
	return itemAction, nil
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

func (s *Service) FindImageByID(itemImageID string) (*ItemImage, error) {
	var itemImage *ItemImage
	result := s.DB.First(&itemImage, "id = ?", itemImageID)
	if result.Error != nil {
		return nil, result.Error
	}
	return itemImage, nil
}

func (s *Service) AddImage(itemID string, image multipart.File, fileHeader *multipart.FileHeader) (*ItemImage, error) {
	objectName := fmt.Sprintf("%s--%s.png", itemID, util.RandomString(12))
	_, err := s.Storage.Upload(objectName, image, "image/*")
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://cdn.safatanc.com/blockstuff/%s", objectName)

	itemImage := &ItemImage{
		URL:        url,
		ItemID:     itemID,
		ObjectName: objectName,
	}

	result := s.DB.Create(&itemImage)
	if result.Error != nil {
		return nil, result.Error
	}

	return itemImage, nil
}

func (s *Service) DeleteImage(itemID string, itemImageID string) (*ItemImage, error) {
	itemImage, err := s.FindImageByID(itemImageID)
	if err != nil {
		return nil, err
	}

	if itemImage.ItemID != itemID {
		return nil, fmt.Errorf("unauthorized")
	}

	if itemImage.ObjectName == "" {
		splittedUrl := strings.Split(itemImage.URL, "/")
		itemImage.ObjectName = splittedUrl[len(splittedUrl)-1]
	}

	err = s.Storage.Delete(itemImage.ObjectName)
	if err != nil {
		return nil, err
	}
	result := s.DB.Where("id = ?", itemImageID).Delete(&itemImage)
	if result.Error != nil {
		return nil, result.Error
	}
	return itemImage, nil
}

func (s *Service) DeleteAction(itemID string, itemActionID string) (*ItemAction, error) {
	itemAction, err := s.FindActionByID(itemActionID)
	if err != nil {
		return nil, err
	}

	if itemAction.ItemID != itemID {
		return nil, fmt.Errorf("unauthorized")
	}

	result := s.DB.Where("id = ?", itemActionID).Delete(&itemAction)
	if result.Error != nil {
		return nil, result.Error
	}
	return itemAction, nil
}
