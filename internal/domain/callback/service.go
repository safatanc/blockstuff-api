package callback

import (
	"fmt"

	"github.com/gorcon/rcon"
	"github.com/safatanc/blockstuff-api/internal/domain/item"
	"github.com/safatanc/blockstuff-api/internal/domain/minecraftserver"
	"gorm.io/gorm"
)

type Service struct {
	DB                     *gorm.DB
	MinecraftServerService *minecraftserver.Service
	ItemService            *item.Service
}

func NewService(db *gorm.DB, minecraftserverService *minecraftserver.Service, itemService *item.Service) *Service {
	return &Service{
		DB:                     db,
		MinecraftServerService: minecraftserverService,
		ItemService:            itemService,
	}
}

func (s *Service) XenditCallback(payload *XenditPayload) error {
	if payload.Event == "payment.succeeded" {
		for _, payloadItem := range payload.Data.Items {
			item, err := s.ItemService.FindByID(payloadItem.ReferenceID)
			if err != nil {
				return err
			}

			minecraftserver, err := s.MinecraftServerService.FindByID(*item.MinecraftServerID)
			if err != nil {
				return err
			}

			rconConnection, err := rcon.Dial(fmt.Sprintf("%v:%v", minecraftserver.MinecraftServerRcon.IP, minecraftserver.MinecraftServerRcon.Port), minecraftserver.MinecraftServerRcon.Password)
			if err != nil {
				return err
			}
			defer rconConnection.Close()

			for _, itemAction := range item.ItemActions {
				if itemAction.Type == "COMMAND" {
					_, err := rconConnection.Execute(itemAction.Action)
					if err != nil {
						return err
					}
				}
			}

		}
	}
	return nil
}
