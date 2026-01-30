package service

import (
	"context"

	"github.com/taskflow/backend/internal/errors"
	"github.com/taskflow/backend/internal/models"
)

// MockTaskRepository es un mock para TaskRepository
type MockTaskRepository struct {
	CreateFunc       func(ctx context.Context, title, description, priority string, dueDate interface{}, createdBy string) (string, error)
	GetByIDFunc      func(ctx context.Context, id string) (*models.Task, error)
	UpdateFunc       func(ctx context.Context, id, title, description, priority string, dueDate interface{}) error
	DeleteFunc       func(ctx context.Context, id string) error
	GetAllFunc       func(ctx context.Context, userID, status string, page, pageSize int) ([]models.Task, int, error)
	GetStatsFunc     func(ctx context.Context, userID string) (*models.TaskStats, error)
	UpdateStatusFunc func(ctx context.Context, id, status string) error
	AssignTaskFunc   func(ctx context.Context, taskID, userID string) error
}

func (m *MockTaskRepository) Create(ctx context.Context, title, description, priority string, dueDate interface{}, createdBy string) (string, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, title, description, priority, dueDate, createdBy)
	}
	return "task-123", nil
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, errors.ErrTaskNotFound
}

func (m *MockTaskRepository) Update(ctx context.Context, id, title, description, priority string, dueDate interface{}) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, title, description, priority, dueDate)
	}
	return nil
}

func (m *MockTaskRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockTaskRepository) GetAll(ctx context.Context, userID, status string, page, pageSize int) ([]models.Task, int, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, userID, status, page, pageSize)
	}
	return nil, 0, nil
}

func (m *MockTaskRepository) GetStats(ctx context.Context, userID string) (*models.TaskStats, error) {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockTaskRepository) UpdateStatus(ctx context.Context, id, status string) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, id, status)
	}
	return nil
}

func (m *MockTaskRepository) AssignTask(ctx context.Context, taskID, userID string) error {
	if m.AssignTaskFunc != nil {
		return m.AssignTaskFunc(ctx, taskID, userID)
	}
	return nil
}
