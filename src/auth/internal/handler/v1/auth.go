package handler

import (
	"auth/internal/dto"
	"auth/internal/usecase"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
	}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginReq dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := h.authUsecase.Login(r.Context(), loginReq.Email, loginReq.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

func (h *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var refreshReq dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := h.authUsecase.Refresh(r.Context(), refreshReq.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var logoutReq dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&logoutReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.authUsecase.Logout(r.Context(), logoutReq.RefreshToken)
	if err != nil {
		http.Error(w, "Failed to log out", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}
