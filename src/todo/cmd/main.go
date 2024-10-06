package main

import (
	"todo/internal/adapter/database"
	"todo/internal/adapter/logger"

	"log"
	"net/http"
	sqlxRepo "todo/internal/adapter/repository/sqlx"
	api "todo/internal/api/v1"
	"todo/internal/config"
	handler "todo/internal/handler/v1"
	"todo/internal/middleware"
	usecase "todo/internal/usecase/v1"

	"github.com/gorilla/mux"
)

func main() {
	config, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Println("Error reading config (config.toml)")
	}

	db, err := database.NewPostgresDB(config.Todo.Postgres)
	if err != nil {
		log.Println("Couldn't connect to database, exiting")
		return
	}

	logger := logger.NewZapLogger(config.Todo.Log)

	boardRepo := sqlxRepo.NewSQLXBoardRepository(db)
	columnRepo := sqlxRepo.NewSQLXColumnRepository(db)
	cardRepo := sqlxRepo.NewSQLXCardRepository(db)

	uc := usecase.NewTodoUseCase(boardRepo, columnRepo, cardRepo)

	userHandler := handler.NewTodoHandler(uc, config.Pagination)
	router := mux.NewRouter()
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	router.Use(loggingMiddleware.Middleware)
	api.InitializeV1Routes(router, userHandler)

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", router)
}
