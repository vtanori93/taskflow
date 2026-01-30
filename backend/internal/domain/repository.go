package domain

import (
	"context"

	"github.com/taskflow/backend/internal/models"
)

// UserRepository define los métodos para acceder a datos de usuarios
type UserRepository interface {
	// GetByEmail obtiene un usuario por email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// GetByID obtiene un usuario por ID
	GetByID(ctx context.Context, id string) (*models.User, error)

	// Create crea un nuevo usuario
	Create(ctx context.Context, email, passwordHash, name string) (string, error)

	// Update actualiza un usuario
	Update(ctx context.Context, user *models.User) error

	// GetAllUsers obtiene todos los usuarios registrados
	GetAllUsers(ctx context.Context) ([]*models.User, error)
}

// TaskRepository define los métodos para acceder a datos de tareas
type TaskRepository interface {
	// GetAll obtiene todas las tareas con filtros y paginación
	GetAll(ctx context.Context, userID, status string, page, pageSize int) ([]models.Task, int, error)

	// GetByID obtiene una tarea por ID
	GetByID(ctx context.Context, id string) (*models.Task, error)

	// Create crea una nueva tarea
	Create(ctx context.Context, title, description, priority string, dueDate interface{}, createdBy string) (string, error)

	// Update actualiza una tarea
	Update(ctx context.Context, id, title, description, priority string, dueDate interface{}) error

	// Delete elimina una tarea
	Delete(ctx context.Context, id string) error

	// UpdateStatus actualiza el estado de una tarea
	UpdateStatus(ctx context.Context, id, status string) error

	// AssignTask asigna una tarea a un usuario
	AssignTask(ctx context.Context, taskID, userID string) error

	// GetStats obtiene estadísticas de tareas de un usuario
	GetStats(ctx context.Context, userID string) (*models.TaskStats, error)
}
