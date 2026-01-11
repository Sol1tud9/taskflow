package domain

import "time"

type UserCreatedEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type UserUpdatedEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TeamUpdatedEvent struct {
	TeamID    string    `json:"team_id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TaskCreatedEvent struct {
	TaskID     string    `json:"task_id"`
	Title      string    `json:"title"`
	CreatorID  string    `json:"creator_id"`
	AssigneeID string    `json:"assignee_id"`
	TeamID     string    `json:"team_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type TaskUpdatedEvent struct {
	TaskID    string    `json:"task_id"`
	UserID    string    `json:"user_id"`
	Field     string    `json:"field"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	UpdatedAt time.Time `json:"updated_at"`
}

