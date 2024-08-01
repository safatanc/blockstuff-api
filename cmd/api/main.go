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
	"github.com/safatanc/blockstuff-api/internal/domain/payout"
	"github.com/safatanc/blockstuff-api/internal/domain/storage"
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

	db.AutoMigrate(
		&user.User{}, &minecraftserver.MinecraftServer{}, &minecraftserver.MinecraftServerRcon{},
		&item.Item{}, &item.ItemImage{}, &item.ItemAction{},
		&transaction.Transaction{}, &transaction.TransactionItem{},
		&payout.Payout{}, &payout.PayoutTransaction{},
	)

	// Validate
	validate := validator.New()

	// Midtrans
	callbackUrl := os.Getenv("CALLBACK_URL")
	midtransCore := &coreapi.Client{
		Options: &midtrans.ConfigOptions{
			PaymentOverrideNotification: &callbackUrl,
		},
	}
	midtransServerKey := os.Getenv("MIDTRANS_SERVER_KEY")
	midtransEnvironment := midtrans.Production
	if strings.Contains(midtransServerKey, "SB") {
		midtransEnvironment = midtrans.Sandbox
	}
	midtransCore.New(midtransServerKey, midtransEnvironment)

	// Middleware
	mw := middleware.New()

	// Domain Storage
	storageService := storage.NewService()

	// Domain User
	userService := user.NewService(db, validate)
	userController := user.NewController(userService)
	userRoutes := user.NewRoutes(mux, userController, mw)
	userRoutes.Init()

	// Domain Auth
	authService := auth.NewService(db, validate, userService)
	authController := auth.NewController(authService)
	authRoutes := auth.NewRoutes(mux, authController)
	authRoutes.Init()

	// Domain MinecraftServer
	minecraftServerService := minecraftserver.NewService(db, validate, storageService)
	minecraftServerController := minecraftserver.NewController(minecraftServerService, userService)
	minecraftServerRoutes := minecraftserver.NewRoutes(mux, minecraftServerController, mw)
	minecraftServerRoutes.Init()

	// Domain Item
	itemService := item.NewService(db, validate, storageService)
	itemController := item.NewController(itemService, userService, minecraftServerService)
	itemRoutes := item.NewRoutes(mux, itemController, mw)
	itemRoutes.Init()

	// Domain Transaction
	transactionService := transaction.NewService(db, validate, midtransCore)
	transactionController := transaction.NewController(transactionService, userService, minecraftServerService)
	transactionRoutes := transaction.NewRoutes(mux, transactionController, mw)
	transactionRoutes.Init()

	// Domain Payout
	payoutService := payout.NewService(db, validate)
	payoutController := payout.NewController(payoutService, userService, itemService, transactionService)
	payoutRoutes := payout.NewRoutes(mux, payoutController, mw)
	payoutRoutes.Init()

	// Server
	log.Printf("Running on http://localhost:%v", PORT)
	srv := server.New(mux, PORT)
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
