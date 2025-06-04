package repository

import "denettest/internal/domain"

type UserRepository interface {
	GetUser(id int64) (domain.User, error)
	GetLeaderboard() ([]domain.User, error)
	SetReferrer(referrefID, userID int64) error
}

type TaskRepository interface {
	CompleteTask(userID int64, taskID int) (int, error)
}
type AuthRepository interface {
	CreateUser(email string, passHash []byte) error
	LoginUser(email string) (domain.User, error)
}
