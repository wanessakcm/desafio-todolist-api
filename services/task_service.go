package services

import (
	"errors"
	"strings"
	"time"

	"desafio-todolist-api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskRepository interface {
	Create(task *models.Task) error
	FindAll(status, priority string) ([]models.Task, error)
	FindByID(id string) (*models.Task, error)
	Update(id string, update bson.M) error
	Delete(id string) error
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

var (
	ErrInvalidTitle      = errors.New("title deve ter entre 3 e 100 caracteres")
	ErrInvalidStatus     = errors.New("status inválido")
	ErrInvalidPriority   = errors.New("priority inválida")
	ErrDueDateInPast     = errors.New("due_date não pode ser no passado")
	ErrCompletedCantEdit = errors.New("tarefa completed não pode ser editada")
	ErrInvalidDueDate    = errors.New("due_date inválido (use YYYY-MM-DD)")
)

// Entradas do JSON

type CreateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
}

type UpdateTaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
}

// Métodos

func (s *TaskService) Create(input CreateTaskInput) (*models.Task, error) {
	title := strings.TrimSpace(input.Title)
	if len(title) < 3 || len(title) > 100 {
		return nil, ErrInvalidTitle
	}

	priority := models.PriorityMedium
	if input.Priority != "" {
		if !isValidPriority(input.Priority) {
			return nil, ErrInvalidPriority
		}
		priority = models.Priority(input.Priority)
	}

	due, err := time.Parse("2006-01-02", input.DueDate)
	if err != nil {
		return nil, ErrInvalidDueDate
	}
	if isPastDate(due) {
		return nil, ErrDueDateInPast
	}

	task := models.NewTask(title, input.Description, priority, due)
	if err := s.repo.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) List(status, priority string) ([]models.Task, error) {
	if status != "" && !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}
	if priority != "" && !isValidPriority(priority) {
		return nil, ErrInvalidPriority
	}
	return s.repo.FindAll(status, priority)
}

func (s *TaskService) GetByID(id string) (*models.Task, error) {
	return s.repo.FindByID(id)
}

func (s *TaskService) Update(id string, input UpdateTaskInput) error {
	current, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if current.Status == models.StatusCompleted {
		return ErrCompletedCantEdit
	}

	update := bson.M{}

	// só atualiza campos que vierem preenchidos
	if strings.TrimSpace(input.Title) != "" {
		title := strings.TrimSpace(input.Title)
		if len(title) < 3 || len(title) > 100 {
			return ErrInvalidTitle
		}
		update["title"] = title
	}

	if input.Description != "" {
		update["description"] = input.Description
	}

	if input.Status != "" {
		if !isValidStatus(input.Status) {
			return ErrInvalidStatus
		}
		update["status"] = input.Status
	}

	if input.Priority != "" {
		if !isValidPriority(input.Priority) {
			return ErrInvalidPriority
		}
		update["priority"] = input.Priority
	}

	if input.DueDate != "" {
		due, err := time.Parse("2006-01-02", input.DueDate)
		if err != nil {
			return ErrInvalidDueDate
		}
		if isPastDate(due) {
			return ErrDueDateInPast
		}
		update["due_date"] = due
	}

	if len(update) == 0 {
		return nil
	}

	update["updated_at"] = time.Now()
	return s.repo.Update(id, update)
}

func (s *TaskService) Delete(id string) error {
	return s.repo.Delete(id)
}

// Validações

func isValidStatus(v string) bool {
	switch v {
	case "pending", "in_progress", "completed", "cancelled":
		return true
	default:
		return false
	}
}

func isValidPriority(v string) bool {
	switch v {
	case "low", "medium", "high":
		return true
	default:
		return false
	}
}

func isPastDate(d time.Time) bool {
	today := time.Now().Truncate(24 * time.Hour)
	date := d.Truncate(24 * time.Hour)
	return date.Before(today)
}
