package domain

type Task struct {
	ID          int    `json:"task_id"`
	Description string `json:"task_description"`
	Reward      int    `json:"task_reward"`
}
