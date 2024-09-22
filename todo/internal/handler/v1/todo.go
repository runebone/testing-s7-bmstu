package v1

import (
	"net/http"
	"todo/internal/usecase"
)

type TodoHandler struct {
	todoUseCase usecase.TodoUseCase
}

func NewTodoHandler(todoUseCase usecase.TodoUseCase) *TodoHandler {
	return &TodoHandler{todoUseCase: todoUseCase}

}

func (h *TodoHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
}

func (h *TodoHandler) GetBoardByID(w http.ResponseWriter, r *http.Request) {
}

func (h *TodoHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
}

func (h *TodoHandler) DeleteBoard(w http.ResponseWriter, r *http.Request) {
}

func (h *TodoHandler) CreateColumn(w http.ResponseWriter, r *http.Request) {
}

func (h *TodoHandler) GetColumnByID(w http.ResponseWriter, r *http.Request) {
}

func (h *TodoHandler) GetColumnsByBoard(w http.ResponseWriter, r *http.Request) {
}

func (h *TodoHandler) UpdateColumn(w http.ResponseWriter, r *http.Request) {
}

func (h *TodoHandler) DeleteColumn(w http.ResponseWriter, r *http.Request) {
}

func (h *TodoHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
}

func (h *TodoHandler) GetCardByID(w http.ResponseWriter, r *http.Request) {
}

func (h *TodoHandler) GetCardsByColumn(w http.ResponseWriter, r *http.Request) {
}

func (h *TodoHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
}

func (h *TodoHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
}
