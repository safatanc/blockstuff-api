package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/safatanc/blockstuff-api/internal/domain/auth"
	"github.com/safatanc/blockstuff-api/internal/domain/item"
	"github.com/safatanc/blockstuff-api/internal/domain/minecraftserver"
	"github.com/safatanc/blockstuff-api/internal/domain/transaction"
	"github.com/safatanc/blockstuff-api/internal/domain/user"
	"github.com/safatanc/blockstuff-api/internal/middleware"
	"github.com/safatanc/blockstuff-api/internal/server"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()
	PORT := 5000
	mux := http.NewServeMux()

	// Database
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")))
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&user.User{}, &minecraftserver.MinecraftServer{}, &item.Item{}, &item.ItemImage{}, &item.ItemAction{}, &transaction.Transaction{})

	// Validate
	validate := validator.New()

	// Midtrans
	midtransCore := coreapi.Client{}
	midtransServerKey := os.Getenv("MIDTRANS_SERVER_KEY")
	midtransEnvironment := midtrans.Production
	if strings.Contains(midtransServerKey, "SB") {
		midtransEnvironment = midtrans.Sandbox
	}
	midtransCore.New(midtransServerKey, midtransEnvironment)

	// Middleware
	mw := middleware.New()

	// Domain Auth
	authService := auth.NewService(db, validate)
	authController := auth.NewController(authService)
	authRoutes := auth.NewRoutes(mux, authController)
	authRoutes.Init()

	// Domain User
	userService := user.NewService(db, validate)
	userController := user.NewController(userService)
	userRoutes := user.NewRoutes(mux, userController, mw)
	userRoutes.Init()

	// Domain MinecraftServer
	minecraftServerService := minecraftserver.NewService(db, validate)
	minecraftServerController := minecraftserver.NewController(minecraftServerService, userService)
	minecraftServerRoutes := minecraftserver.NewRoutes(mux, minecraftServerController, mw)
	minecraftServerRoutes.Init()

	// Domain Item
	itemService := item.NewService(db, validate)
	itemController := item.NewController(itemService, userService, minecraftServerService)
	itemRoutes := item.NewRoutes(mux, itemController, mw)
	itemRoutes.Init()

	// Domain Transaction
	transactionService := transaction.NewService(db, validate)
	transactionController := transaction.NewController(transactionService, userService, minecraftServerService)
	transactionRoutes := transaction.NewRoutes(mux, transactionController, mw)
	transactionRoutes.Init()

	// Server
	log.Printf("Running on http://localhost:%v", PORT)
	srv := server.New(mux, PORT)
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
