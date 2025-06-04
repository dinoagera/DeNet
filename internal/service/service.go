package service

import (
	"denettest/internal/config"
	"denettest/internal/domain"
	jwtour "denettest/internal/middleware/jwt"
	repository "denettest/internal/repository/interfaces"
	storageerrors "denettest/internal/repository/storageErrors"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	log            *slog.Logger
	authRepository repository.AuthRepository
	userRepository repository.UserRepository
	taskRepository repository.TaskRepository
}

func New(log *slog.Logger,
	authRepository repository.AuthRepository,
	userRepository repository.UserRepository,
	taskRepository repository.TaskRepository,
) *Service {
	return &Service{
		log:            log,
		authRepository: authRepository,
		userRepository: userRepository,
		taskRepository: taskRepository,
	}
}
func (s *Service) Register(email, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to generate password hash")
		return err
	}
	err = s.authRepository.CreateUser(email, passHash)
	if err != nil {
		s.log.Error("failed to create user", "error:", err)
		return err
	}
	return nil
}
func (s *Service) Login(email, password string) (string, error) {
	user, err := s.authRepository.LoginUser(email)
	if err != nil {
		if errors.Is(err, storageerrors.ErrUserNotFound) {
			s.log.Warn("user not found")
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("internal server error")
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Info("invalid password", "email", email)
		return "", fmt.Errorf("invalid credentials")
	}
	token, err := jwtour.NewToken(user, config.GetConfig().TTL)
	if err != nil {
		s.log.Error("failed to generate token", "email", email, "error", err)
		return "", fmt.Errorf("failed to generate token")
	}
	return token, nil
}
func (s *Service) GetUser(id int64) (domain.User, error) {
	user, err := s.userRepository.GetUser(id)
	if err != nil {
		s.log.Error("failed to get user", "error:", err)
		return domain.User{}, err
	}
	return user, nil
}
func (s *Service) GetLeaderboard() ([]domain.User, error) {
	users, err := s.userRepository.GetLeaderboard()
	if err != nil {
		s.log.Error("failed to get leaderboard", "error", err)
		return nil, err
	}
	if users == nil {
		return nil, fmt.Errorf("leaderboard is empty")
	}
	return users, nil
}
func (s *Service) CompleteTask(userID int64, taskID int) (int, error) {
	reward, err := s.taskRepository.CompleteTask(userID, taskID)
	if err != nil {
		s.log.Error("failed to completed", "error", err)
		return 0, err
	}
	return reward, nil
}
func (s *Service) SetReferrer(referrefID, userID int64) error {
	err := s.userRepository.SetReferrer(referrefID, userID)
	if err != nil {
		s.log.Error("failed to set referred id", "error", err)
		return err
	}
	return nil
}
