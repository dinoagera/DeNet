package postgres

import (
	"context"
	"denettest/internal/domain"
	storageerrors "denettest/internal/repository/storageErrors"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthStorage struct {
	Pool *pgxpool.Pool
}
type UserStorage struct {
	Pool *pgxpool.Pool
}
type TaskStorage struct {
	Pool *pgxpool.Pool
}

func New(storagePath string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping to database: %w", err)
	}
	return pool, nil
}
func NewAuthStorage(pool *pgxpool.Pool) *AuthStorage {
	return &AuthStorage{Pool: pool}
}
func NewUserStorage(pool *pgxpool.Pool) *UserStorage {
	return &UserStorage{Pool: pool}
}
func NewTaskStorage(pool *pgxpool.Pool) *TaskStorage {
	return &TaskStorage{Pool: pool}
}
func (s *AuthStorage) CreateUser(email string, passHash []byte) error {
	_, err := s.Pool.Exec(context.Background(), "INSERT INTO users(email, pass_hash) VALUES ($1, $2)", email, passHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return storageerrors.ErrUserExists
		}
		return fmt.Errorf("failed to add user:%w", err)
	}
	return nil
}
func (s *AuthStorage) LoginUser(email string) (domain.User, error) {
	var user domain.User
	err := s.Pool.QueryRow(
		context.Background(),
		"SELECT id, email, pass_hash FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, storageerrors.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("failed to login user: %w", err)
	}
	return user, nil
}
func (s *UserStorage) GetUser(id int64) (domain.User, error) {
	var user domain.User
	err := s.Pool.QueryRow(
		context.Background(),
		"SELECT id, email, balance FROM users WHERE id=$1",
		id,
	).Scan(&user.ID, &user.Email, &user.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, storageerrors.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
func (s *UserStorage) GetLeaderboard() ([]domain.User, error) {
	var users []domain.User
	rows, err := s.Pool.Query(context.Background(),
		"SELECT email, balance FROM users ORDER BY balance")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.Email,
			&user.Balance,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (s *TaskStorage) CompleteTask(userID int64, taskID int) (int, error) {
	var reward int
	err := s.Pool.QueryRow(context.Background(),
		"SELECT reward FROM tasks WHERE id = $1",
		taskID).Scan(&reward)
	if err != nil {
		return 0, err
	}
	_, err = s.Pool.Exec(context.Background(),
		`INSERT INTO user_tasks (user_id, task_id) 
		VALUES ($1, $2)
		ON CONFLICT (user_id, task_id) DO NOTHING`,
		userID, taskID)
	if err != nil {
		return 0, err
	}
	_, err = s.Pool.Exec(context.Background(),
		"UPDATE users SET balance = balance + $1 WHERE id = $2",
		reward, userID)
	if err != nil {
		return 0, err
	}
	return reward, nil
}
func (s *UserStorage) SetReferrer(referrerID, userID int64) error {
	_, err := s.Pool.Exec(context.Background(),
		"UPDATE users SET referrer_id = $1 WHERE id = $2",
		referrerID, userID)
	if err != nil {
		return err
	}
	return nil
}
