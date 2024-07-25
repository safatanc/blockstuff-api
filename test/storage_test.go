package test

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/safatanc/blockstuff-api/internal/domain/storage"
)

func TestFindObjects(t *testing.T) {
	godotenv.Load()
	s := storage.NewService()
	objects := s.FindAll()
	for obj := range objects {
		log.Println(obj.Key)
	}
}
