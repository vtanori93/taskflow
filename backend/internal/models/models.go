package models

import (
	"time"
)

// ...existing code...

// User representa un usuario del sistema
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Task representa una tarea del sistema
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	CreatedBy   string     `json:"created_by"`
	AssignedTo  *string    `json:"assigned_to"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TaskStats contiene estadísticas de tareas de un usuario
type TaskStats struct {
	TotalTasks        int `json:"total_tasks"`
	PendingCount      int `json:"pending_count"`
	InProgressCount   int `json:"in_progress_count"`
	CompletedCount    int `json:"completed_count"`
	CancelledCount    int `json:"cancelled_count"`
	HighPriorityCount int `json:"high_priority_count"`
	OverdueCount      int `json:"overdue_count"`
}

// ============================================================================
// REQUEST/RESPONSE MODELS
// ============================================================================

// RegisterRequest modelo de request para registrar usuario
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// LoginRequest modelo de request para login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse respuesta de login con tokens
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	User         User   `json:"user"`
}

// RefreshTokenRequest modelo para refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// CreateTaskRequest modelo para crear tarea
type CreateTaskRequest struct {
	Title       string  `json:"title" binding:"required,max=100"`
	Description string  `json:"description" binding:"max=500"`
	Priority    string  `json:"priority" binding:"required,oneof=low medium high urgent"`
	DueDate     *string `json:"due_date,omitempty"` // Fecha como string plano
}

// UpdateTaskRequest modelo para actualizar tarea
type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty" binding:"omitempty,max=100"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=500"`
	Priority    *string `json:"priority,omitempty" binding:"omitempty,oneof=low medium high urgent"`
	DueDate     *string `json:"due_date,omitempty"` // Fecha como string plano
}

// UpdateTaskStatusRequest modelo para cambiar estado
type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending in_progress completed cancelled"`
}

// AssignTaskRequest modelo para asignar tarea
type AssignTaskRequest struct {
	AssignedTo *string `json:"assigned_to"`
}

// TasksListResponse respuesta con lista de tareas
type TasksListResponse struct {
	Tasks      []Task `json:"tasks"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	TotalPages int    `json:"total_pages"`
}

// APIResponse respuesta genérica de API
type APIResponse struct {
	Data       interface{} `json:"data"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Error      bool        `json:"error"`
}

// PaginationParams parámetros de paginación
type PaginationParams struct {
	Page     int
	PageSize int
}
