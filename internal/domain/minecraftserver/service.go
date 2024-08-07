package minecraftserver

import (
	"fmt"
	"math"
	"mime/multipart"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/safatanc/blockstuff-api/internal/domain/storage"
	"github.com/safatanc/blockstuff-api/pkg/converter"
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

func (s *Service) FindAll() []*MinecraftServer {
	var minecraftservers = make([]*MinecraftServer, 0)
	s.DB.Find(&minecraftservers)

	return minecraftservers
}

func (s *Service) FindByID(id string) (*MinecraftServer, error) {
	var minecraftserver *MinecraftServer
	result := s.DB.First(&minecraftserver, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return minecraftserver, nil
}

func (s *Service) FindBySlug(slug string, detail bool) (*MinecraftServer, error) {
	var minecraftserver *MinecraftServer
	var result *gorm.DB
	if detail {
		result = s.DB.Preload("Author").Preload("MinecraftServerRcon").First(&minecraftserver, "slug = ?", slug)
		minecraftserver.Author = minecraftserver.Author.ToResponse()
		minecraftserver.MinecraftServerRcon.Password = ""
	} else {
		result = s.DB.First(&minecraftserver, "slug = ?", slug)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return minecraftserver, nil
}

func (s *Service) Create(minecraftserver *MinecraftServer) (*MinecraftServer, error) {
	err := s.Validate.Struct(minecraftserver)
	if err != nil {
		return nil, err
	}

	var findMinecraftServer *MinecraftServer
	result := s.DB.Order("created_at DESC").First(&findMinecraftServer, "author_id = ?", minecraftserver.AuthorID)
	if result.Error == nil {
		difference := time.Until(findMinecraftServer.CreatedAt)
		if math.Abs(difference.Seconds()) < 60 {
			return nil, fmt.Errorf("request limit reached. cooldown %.2f seconds", 60-math.Abs(difference.Seconds()))
		}
	}

	result = s.DB.Create(&minecraftserver)
	if result.Error != nil {
		return nil, result.Error
	}

	s.DB.Preload("Author").Preload("MinecraftServerRcon").First(&minecraftserver, "id = ?", minecraftserver.ID)

	minecraftserver.Author = minecraftserver.Author.ToResponse()
	return minecraftserver, nil
}

func (s *Service) Update(id string, minecraftserver *MinecraftServer) (*MinecraftServer, error) {
	err := s.Validate.Struct(minecraftserver)
	if err != nil {
		return nil, err
	}

	result := s.DB.Where("id = ?", id).Updates(&minecraftserver)
	if result.Error != nil {
		return nil, result.Error
	}
	return minecraftserver, nil
}

func (s *Service) Delete(id string) (*MinecraftServer, error) {
	var minecraftserver *MinecraftServer
	result := s.DB.First(&minecraftserver, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	result = s.DB.Where("id = ?", id).Delete(&minecraftserver)
	if result.Error != nil {
		return nil, result.Error
	}
	return minecraftserver, nil
}

func (s *Service) UpdateRcon(rcon *MinecraftServerRcon) (*MinecraftServer, error) {
	err := s.Validate.Struct(rcon)
	if err != nil {
		return nil, err
	}

	encryptedPassword, err := converter.EncryptPassword(rcon.Password)
	if err != nil {
		return nil, err
	}

	rcon.Password = *encryptedPassword

	result := s.DB.Where("minecraft_server_id = ?", rcon.MinecraftServerID).Updates(&rcon)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		result = s.DB.Create(&rcon)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	var minecraftserver *MinecraftServer
	s.DB.Preload("Author").Preload("MinecraftServerRcon").First(&minecraftserver, "id = ?", rcon.MinecraftServerID)

	minecraftserver.Author = minecraftserver.Author.ToResponse()
	return minecraftserver, nil
}

func (s *Service) UpdateLogo(id string, image multipart.File, fileHeader *multipart.FileHeader) (*MinecraftServer, error) {
	minecraftserver, err := s.FindByID(id)
	if err != nil {
		return nil, err
	}

	if minecraftserver.Logo != nil {
		splittedUrl := strings.Split(*minecraftserver.Logo, "/")
		currentObjectName := splittedUrl[len(splittedUrl)-1]

		_, err = s.Storage.FindByObjectName(currentObjectName)
		if err == nil {
			err := s.Storage.Delete(currentObjectName)
			if err != nil {
				return nil, err
			}
		}
	}

	objectName := fmt.Sprintf("%s--%s.png", id, util.RandomString(4))

	_, err = s.Storage.Upload(objectName, image, "image/*")
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://cdn.safatanc.com/blockstuff/%s", objectName)
	minecraftserver = &MinecraftServer{
		Logo: &url,
	}
	result := s.DB.Where("id = ?", id).Updates(&minecraftserver)
	if result.Error != nil {
		return nil, result.Error
	}

	return minecraftserver, nil
}
