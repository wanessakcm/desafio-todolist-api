package models

import (
	"time"

	"github.com/google/uuid"
)

type Status string
type Priority string

const (
	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusCancelled  Status = "cancelled"

	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	ID          string    `json:"id" bson:"id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Status      Status    `json:"status" bson:"status"`
	Priority    Priority  `json:"priority" bson:"priority"`
	DueDate     time.Time `json:"due_date" bson:"due_date"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// Função auxiliar para criar nova task com valores padrão
func NewTask(title, description string, priority Priority, dueDate time.Time) *Task {
	now := time.Now()

	return &Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Status:      StatusPending,
		Priority:    priority,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
