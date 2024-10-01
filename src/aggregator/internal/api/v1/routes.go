package v1

import (
	v1 "aggregator/internal/handler/v1"

	"github.com/gorilla/mux"
)

func InitializeV1Routes(router *mux.Router, aggHandler *v1.AggregatorHandler) {
	router.HandleFunc("/api/v1/register", aggHandler.Register).Methods("POST")
	router.HandleFunc("/api/v1/login", aggHandler.Login).Methods("POST")
	router.HandleFunc("/api/v1/logout", aggHandler.Logout).Methods("POST")

	router.HandleFunc("/api/v1/boards", aggHandler.GetBoards).Methods("GET")      // Boards
	router.HandleFunc("/api/v1/board/{id}", aggHandler.GetBoard).Methods("GET")   // Columns + cards
	router.HandleFunc("/api/v1/column/{id}", aggHandler.GetColumn).Methods("GET") // Cards
	router.HandleFunc("/api/v1/card/{id}", aggHandler.GetCard).Methods("GET")     // Card + description

	router.HandleFunc("/api/v1/board", aggHandler.CreateBoard).Methods("POST")
	router.HandleFunc("/api/v1/column", aggHandler.CreateColumn).Methods("POST")
	router.HandleFunc("/api/v1/card", aggHandler.CreateCard).Methods("POST")

	router.HandleFunc("/api/v1/board", aggHandler.UpdateBoard).Methods("PUT")
	router.HandleFunc("/api/v1/column", aggHandler.UpdateColumn).Methods("PUT")
	router.HandleFunc("/api/v1/card", aggHandler.UpdateCard).Methods("PUT")

	router.HandleFunc("/api/v1/board/{id}", aggHandler.DeleteBoard).Methods("DELETE")
	router.HandleFunc("/api/v1/column/{id}", aggHandler.DeleteColumn).Methods("DELETE")
	router.HandleFunc("/api/v1/card/{id}", aggHandler.DeleteCard).Methods("DELETE")

	router.HandleFunc("/api/v1/stats/{from}/{to}", aggHandler.GetStats).Methods("GET")
	router.HandleFunc("/api/v1/stats/{from}", aggHandler.GetStats).Methods("GET")
	router.HandleFunc("/api/v1/stats", aggHandler.GetStats).Methods("GET")
}
