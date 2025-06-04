package handlers

import (
	"denettest/internal/domain"
	service "denettest/internal/service/interfaces"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	log         *slog.Logger
	userService service.UserService
	taskService service.TaskService
	authService service.AuthService
}

func New(
	log *slog.Logger,
	userService service.UserService,
	taskService service.TaskService,
	authService service.AuthService,
) *Handler {
	return &Handler{
		log:         log,
		userService: userService,
		taskService: taskService,
		authService: authService,
	}
}
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.log.Info("invalid method")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are empty", http.StatusBadRequest)
		return
	}
	err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		h.log.Info("failed to register user")
		http.Error(w, "failed to register", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode("user created successfully")
}
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.log.Info("invalid method")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are empty", http.StatusBadRequest)
		return
	}
	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, "failed to login", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.log.Info("invalid method")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		h.log.Info("empty user id")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Info("failed to parse to int64")
		http.Error(w, "failed to get status user", http.StatusBadRequest)
		return
	}
	user, err := h.userService.GetUser(id)
	if err != nil {
		h.log.Info("failed to get user status")
		http.Error(w, "failed to get status user", http.StatusBadRequest)
		return
	}
	if user.ID == 0 && user.Email == "" {
		h.log.Info("user is not founded")
		http.Error(w, "user is not founded", http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		h.log.Error("failed to encode response", "error", err)
		http.Error(w, "Failed to prepare response", http.StatusInternalServerError)
		return
	}
}
func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.log.Info("invalid method")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	users, err := h.userService.GetLeaderboard()
	if err != nil {
		h.log.Info("failed to get leaderboard")
		http.Error(w, "failed to get leaderboard", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		h.log.Error("failed to encode response", "error", err)
		http.Error(w, "Failed to prepare response", http.StatusInternalServerError)
		return
	}
}
func (h *Handler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.log.Info("invalid method")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		h.log.Info("empty user id")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	var task domain.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		h.log.Error("failed to decode response", "error", err)
		http.Error(w, "Failed to prepare response", http.StatusInternalServerError)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Info("failed to parse to int64")
		http.Error(w, "failed to complete task", http.StatusInternalServerError)
		return
	}
	reward, err := h.taskService.CompleteTask(id, task.ID)
	if err != nil {
		h.log.Info("failed to complete task")
		http.Error(w, "failed to complete task", http.StatusInternalServerError)
	}
	response := map[string]interface{}{
		"message": fmt.Sprintf("Task completed, you get reward: %d", reward),
	}

	json.NewEncoder(w).Encode(response)
}
func (h *Handler) SetReferrer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		h.log.Info("empty user id")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	userid, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Info("failed to parse")
		http.Error(w, "failed to set referrer", http.StatusInternalServerError)
		return
	}
	var req struct {
		ReferrerID int64 `json:"referrer_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Info("failed to decode")
		http.Error(w, "Failed to prepare response", http.StatusInternalServerError)
		return
	}
	err = h.userService.SetReferrer(req.ReferrerID, userid)
	if err != nil {
		h.log.Info("failed to set referret in service")
		http.Error(w, "failed to set reffere", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode("Referre set is successufully")
}
