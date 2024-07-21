package minecraftserver

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

func (s *Service) FindAll() []*MinecraftServer {
	var minecraftservers []*MinecraftServer
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

func (s *Service) FindByIP(ip string, detail bool) (*MinecraftServer, error) {
	var minecraftserver *MinecraftServer
	var result *gorm.DB
	if detail {
		result = s.DB.Preload("Author").Preload("MinecraftServerRcon").First(&minecraftserver, "ip = ?", ip)
		minecraftserver.Author = minecraftserver.Author.ToResponse()
	} else {
		result = s.DB.First(&minecraftserver, "ip = ?", ip)
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

	result := s.DB.Create(&minecraftserver)
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
	minecraftserver.Author = minecraftserver.Author.ToResponse()
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
