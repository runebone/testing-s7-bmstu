package main

import (
	"log"
	"net/http"
	"user/internal/adapter/database"
	"user/internal/adapter/logger"
	sqlxRepo "user/internal/adapter/repository/sqlx"
	api "user/internal/api/v1"
	"user/internal/config"
	handler "user/internal/handler/v1"
	"user/internal/middleware"
	usecase "user/internal/usecase/v1"

	"github.com/gorilla/mux"
)

func main() {
	logger := logger.NewZapLogger()

	config, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Println("Error reading config (config.toml)")
	}

	db, err := database.NewPostgresDB(config.Database)
	if err != nil {
		log.Println("Couldn't connect to database, exiting")
		return
	}

	repo := sqlxRepo.NewSQLXUserRepository(db)
	uc := usecase.NewUserUseCase(repo)

	userHandler := handler.NewUserHandler(uc)
	router := mux.NewRouter()
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	router.Use(loggingMiddleware.Middleware)
	api.InitializeV1Routes(router, userHandler)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}
