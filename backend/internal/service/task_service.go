package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/taskflow/backend/internal/domain"
	"github.com/taskflow/backend/internal/errors"
	"github.com/taskflow/backend/internal/models"
)

// TaskService maneja la lógica de negocio de tareas
type TaskService struct {
	taskRepo domain.TaskRepository
}

// NewTaskService crea una nueva instancia de TaskService
func NewTaskService(taskRepo domain.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

// CreateTask crea una nueva tarea
func (s *TaskService) CreateTask(ctx context.Context, req *models.CreateTaskRequest, userID string) (*models.Task, error) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("🔴 PANIC en CreateTask: %v\n", rec)
		}
	}()

	log.Printf("📝 CreateTask Service - Input: title=%s, userID=%s\n", req.Title, userID)

	// Convertir string de fecha a *time.Time
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		formats := []string{
			time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
			"2006-01-02T15:04:05", // Sin timezone
			"2006-01-02 15:04:05", // Con espacio
			"2006-01-02",          // Solo fecha
		}

		var parseErr error
		for _, format := range formats {
			if t, err := time.Parse(format, *req.DueDate); err == nil {
				dueDate = &t
				log.Printf("📅 Due date parsed: %v (format: %s)\n", dueDate, format)
				break
			} else {
				parseErr = err
			}
		}

		if dueDate == nil && parseErr != nil {
			return nil, errors.NewBadRequest(fmt.Sprintf("formato de fecha inválido: %s", *req.DueDate))
		}
	}

	taskID, err := s.taskRepo.Create(ctx, req.Title, req.Description, req.Priority, dueDate, userID)
	if err != nil {
		log.Printf("🔴 ERROR en CreateTask Service - Repository Error: %v (type: %T)\n", err, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al crear tarea: %v", err))
	}

	log.Printf("✅ Task created in repository: %s\n", taskID)

	// Obtener la tarea creada
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		log.Printf("🔴 ERROR en CreateTask Service - GetByID Error: %v (type: %T)\n", err, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al obtener tarea: %v", err))
	}

	log.Printf("✅ Task retrieved successfully: %+v\n", task)
	return task, nil
}

// GetTasks obtiene tareas del usuario con paginación
func (s *TaskService) GetTasks(ctx context.Context, userID string, status string, page, pageSize int) (*models.TasksListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	tasks, total, err := s.taskRepo.GetAll(ctx, userID, status, page, pageSize)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al obtener tareas: %v", err))
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &models.TasksListResponse{
		Tasks:      tasks,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// GetTaskByID obtiene una tarea específica
func (s *TaskService) GetTaskByID(ctx context.Context, taskID string) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, errors.ErrTaskNotFound
	}
	return task, nil
}

// UpdateTask actualiza una tarea
func (s *TaskService) UpdateTask(ctx context.Context, taskID string, req *models.UpdateTaskRequest) (*models.Task, error) {
	// Obtener la tarea actual
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, errors.ErrTaskNotFound
	}

	// Preparar valores para actualizar (usar valores actuales si no se proporcionan nuevos)
	title := task.Title
	if req.Title != nil {
		title = *req.Title
	}

	description := task.Description
	if req.Description != nil {
		description = *req.Description
	}

	priority := task.Priority
	if req.Priority != nil {
		priority = *req.Priority
	}

	dueDate := task.DueDate
	if req.DueDate != nil && *req.DueDate != "" {
		formats := []string{
			time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
			"2006-01-02T15:04:05", // Sin timezone
			"2006-01-02 15:04:05", // Con espacio
			"2006-01-02",          // Solo fecha
		}

		var parseErr error
		for _, format := range formats {
			if t, err := time.Parse(format, *req.DueDate); err == nil {
				dueDate = &t
				break
			} else {
				parseErr = err
			}
		}

		if dueDate == nil && parseErr != nil {
			return nil, errors.NewBadRequest(fmt.Sprintf("formato de fecha inválido: %s", *req.DueDate))
		}
	}

	// Actualizar tarea
	err = s.taskRepo.Update(ctx, taskID, title, description, priority, dueDate)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al actualizar tarea: %v", err))
	}

	// Obtener tarea actualizada
	return s.taskRepo.GetByID(ctx, taskID)
}

// DeleteTask elimina una tarea
func (s *TaskService) DeleteTask(ctx context.Context, taskID string) error {
	// Verificar que la tarea existe
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return errors.ErrTaskNotFound
	}

	err = s.taskRepo.Delete(ctx, taskID)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error al eliminar tarea: %v", err))
	}

	return nil
}

// UpdateTaskStatus actualiza el estado de una tarea
func (s *TaskService) UpdateTaskStatus(ctx context.Context, taskID string, req *models.UpdateTaskStatusRequest) (*models.Task, error) {
	// Verificar que la tarea existe
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, errors.ErrTaskNotFound
	}

	// Actualizar estado
	err = s.taskRepo.UpdateStatus(ctx, taskID, req.Status)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al actualizar estado: %v", err))
	}

	// Obtener tarea actualizada
	return s.taskRepo.GetByID(ctx, taskID)
}

// AssignTask asigna una tarea a un usuario
func (s *TaskService) AssignTask(ctx context.Context, taskID string, req *models.AssignTaskRequest) (*models.Task, error) {
	// Verificar que la tarea existe
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, errors.ErrTaskNotFound
	}

	// Asignar tarea
	userID := ""
	if req.AssignedTo != nil {
		userID = *req.AssignedTo
	}

	err = s.taskRepo.AssignTask(ctx, taskID, userID)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al asignar tarea: %v", err))
	}

	// Obtener tarea actualizada
	return s.taskRepo.GetByID(ctx, taskID)
}

// GetTaskStats obtiene estadísticas de tareas del usuario
func (s *TaskService) GetTaskStats(ctx context.Context, userID string) (*models.TaskStats, error) {
	stats, err := s.taskRepo.GetStats(ctx, userID)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al obtener estadísticas: %v", err))
	}
	return stats, nil
}
