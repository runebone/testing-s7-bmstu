package main

import (
	"aggregator/internal/adapter/logger"
	"aggregator/internal/middleware"
	"fmt"
	"time"

	httpAuth "aggregator/internal/adapter/service/auth/http"
	httpTodo "aggregator/internal/adapter/service/todo/http"
	httpUser "aggregator/internal/adapter/service/user/http"
	api "aggregator/internal/api/v1"
	"aggregator/internal/config"
	h "aggregator/internal/handler/v1"
	"aggregator/internal/service/auth"
	"aggregator/internal/service/todo"
	"aggregator/internal/service/user"
	v1 "aggregator/internal/usecase/v1"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config, err := config.LoadConfig("config.toml")
	if err != nil {
		log.Println("Error reading config (config.toml)")
	}

	var userSvc user.UserService
	{
		baseURL := fmt.Sprintf("http://%s:%d/%s", config.User.ContainerName, config.User.LocalPort, config.User.BaseURL)
		logger := logger.NewZapLogger(config.User.Log)
		userSvc = httpUser.NewUserService(baseURL, 2*time.Second, logger)
	}

	var authSvc auth.AuthService
	{
		baseURL := fmt.Sprintf("http://%s:%d/%s", config.Auth.ContainerName, config.Auth.LocalPort, config.Auth.BaseURL)
		logger := logger.NewZapLogger(config.Auth.Log)
		authSvc = httpAuth.NewAuthService(baseURL, 2*time.Second, logger)
	}

	var todoSvc todo.TodoService
	{
		baseURL := fmt.Sprintf("http://%s:%d/%s", config.Todo.ContainerName, config.Todo.LocalPort, config.Todo.BaseURL)
		logger := logger.NewZapLogger(config.Todo.Log)
		todoSvc = httpTodo.NewTodoService(baseURL, 2*time.Second, logger)
	}

	logger := logger.NewZapLogger(config.Aggregator.Log)
	uc := v1.NewAggregatorUseCase(userSvc, authSvc, todoSvc)
	handler := h.NewAggregatorHandler(uc)

	router := mux.NewRouter()
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)
	// authMiddleware := middleware.NewAuthMiddleware(authSvc)
	router.Use(loggingMiddleware.Middleware)
	// router.Use(authMiddleware.Middleware)
	api.InitializeV1Routes(router, handler)

	localPort := fmt.Sprintf("%d", config.Aggregator.LocalPort)
	exposedPort := fmt.Sprintf("%d", config.Aggregator.ExposedPort)

	log.Printf("Starting server on :%s\n", exposedPort)
	http.ListenAndServe(":"+localPort, router)
}
