package main

import (
	"auth/internal/adapter/database"
	"auth/internal/adapter/logger"

	// loggingRepo "auth/internal/adapter/repository/logging"
	sqlxRepo "auth/internal/adapter/repository/sqlx"
	// loggingUseCase "auth/internal/adapter/usecase/logging"
	"auth/internal/adapter/service/tokengen/jwt"
	user "auth/internal/adapter/service/user/http"
	api "auth/internal/api/v1"
	"auth/internal/config"
	handler "auth/internal/handler/v1"
	"auth/internal/middleware"
	usecase "auth/internal/usecase/v1"
	"log"
	"net/http"

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

	repo := sqlxRepo.NewSQLXTokenRepository(db)
	// repo := loggingRepo.NewLoggingAuthRepository(baseRepo, logger)

	userService := user.NewHTTPUserService("http://userservice:8080/api/v1", 2) // XXX:
	tokenService := jwt.NewJWTService("nigger", 15*60, 7*60*60*24)              // XXX:

	uc := usecase.NewAuthUseCase(repo, userService, tokenService)
	// uc := loggingUseCase.NewLoggingAuthUseCase(userUC, logger)

	userHandler := handler.NewAuthHandler(uc)
	router := mux.NewRouter()
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	router.Use(loggingMiddleware.Middleware)
	api.InitializeV1Routes(router, userHandler)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}
