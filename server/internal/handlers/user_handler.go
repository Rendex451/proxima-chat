package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"server/internal/models"
)

type UserService interface {
	CreateUser(c context.Context, req *models.CreateUserReq) (*models.CreateUserRes, error)
	Login(c context.Context, req *models.LoginUserReq) (*models.LoginUserRes, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

// POST /register
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	res, err := h.service.CreateUser(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// POST /login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	res, err := h.service.Login(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Устанавливаем JWT в cookie (пример: res.AccessToken)
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    res.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true для HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   60 * 60 * 24, // 1 день
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// POST /logout
func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Удаляем JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // удалить cookie
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "logout successful",
	})
}
