package domain

type User struct {
	ID         int64  `json:"user_id"`
	Email      string `json:"email"`
	PassHash   []byte
	Balance    int   `json:"balance"`
	ReferrerID int64 `json:"referrer_id"`
}

type UserTask struct {
	ID        uint `json:"user_task_id"`
	UserID    uint `json:"user_id"`
	TaskID    uint `json:"task_id"`
	Completed bool `json:"task_completed"`
}
