package service

import (
	"denettest/internal/domain"
	repository "denettest/internal/repository/interfaces"
	"log/slog"
)

type UserService interface {
	GetUser(id int64) (domain.User, error)
	GetLeaderboard() ([]domain.User, error)
	SetReferrer(referrefID, userID int64) error
}

type TaskService interface {
	CompleteTask(userID int64, taskID int) (int, error)
}
type AuthService interface {
	Register(email, password string) error
	Login(email, password string) (string, error)
}
type Service struct {
	log            *slog.Logger
	authRepository repository.AuthRepository
	userRepository repository.UserRepository
	taskRepository repository.TaskRepository
}
