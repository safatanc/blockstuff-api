package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/safatanc/blockstuff-api/internal/domain/auth"
	"github.com/safatanc/blockstuff-api/internal/domain/callback"
	"github.com/safatanc/blockstuff-api/internal/domain/item"
	"github.com/safatanc/blockstuff-api/internal/domain/minecraftserver"
	"github.com/safatanc/blockstuff-api/internal/domain/payout"
	"github.com/safatanc/blockstuff-api/internal/domain/storage"
	"github.com/safatanc/blockstuff-api/internal/domain/transaction"
	"github.com/safatanc/blockstuff-api/internal/domain/user"
	"github.com/safatanc/blockstuff-api/internal/middleware"
	"github.com/safatanc/blockstuff-api/internal/server"
	"github.com/xendit/xendit-go/v6"
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

	// db.AutoMigrate(
	// 	&user.User{}, &minecraftserver.MinecraftServer{}, &minecraftserver.MinecraftServerRcon{},
	// 	&item.Item{}, &item.ItemImage{}, &item.ItemAction{},
	// 	&transaction.Transaction{}, &transaction.TransactionItem{},
	// 	&payout.Payout{}, &payout.PayoutTransaction{},
	// )
	db.AutoMigrate(
		&transaction.Transaction{},
	)

	// Validate
	validate := validator.New()

	// Xendit
	xenditClient := xendit.NewClient(os.Getenv("XENDIT_SECRET_KEY"))

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
	transactionService := transaction.NewService(db, validate, xenditClient)
	transactionController := transaction.NewController(transactionService, userService, minecraftServerService)
	transactionRoutes := transaction.NewRoutes(mux, transactionController, mw)
	transactionRoutes.Init()

	// Domain Payout
	payoutService := payout.NewService(db, validate)
	payoutController := payout.NewController(payoutService, userService, itemService, transactionService)
	payoutRoutes := payout.NewRoutes(mux, payoutController, mw)
	payoutRoutes.Init()

	// Domain Callback
	callbackService := callback.NewService(db, minecraftServerService, itemService)
	callbackController := callback.NewController(callbackService)
	callbackRoutes := callback.NewRoutes(mux, callbackController)
	callbackRoutes.Init()

	// Server
	log.Printf("Running on http://localhost:%v", PORT)
	srv := server.New(mux, PORT)
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
