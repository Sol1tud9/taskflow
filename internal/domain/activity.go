package domain

import "time"

type EntityType string

const (
	EntityTypeUser EntityType = "user"
	EntityTypeTeam EntityType = "team"
	EntityTypeTask EntityType = "task"
)

type ActionType string

const (
	ActionTypeCreated ActionType = "created"
	ActionTypeUpdated ActionType = "updated"
	ActionTypeDeleted ActionType = "deleted"
)

type Activity struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	EntityType EntityType `json:"entity_type"`
	EntityID   string     `json:"entity_id"`
	Action     ActionType `json:"action"`
	Metadata   string     `json:"metadata"`
	CreatedAt  time.Time  `json:"created_at"`
}
