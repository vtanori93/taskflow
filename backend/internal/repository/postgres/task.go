package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/taskflow/backend/internal/domain"
	"github.com/taskflow/backend/internal/models"
)

// TaskRepository implementa domain.TaskRepository usando PostgreSQL
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository crea una nueva instancia de TaskRepository
func NewTaskRepository(db *sql.DB) domain.TaskRepository {
	return &TaskRepository{db: db}
}

// GetAll obtiene todas las tareas con filtros y paginación
func (r *TaskRepository) GetAll(ctx context.Context, userID, status string, page, pageSize int) ([]models.Task, int, error) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("🔴 PANIC en GetAll: %v\n", rec)
		}
	}()

	offset := (page - 1) * pageSize

	// Construir query para contar tareas
	countQuery := "SELECT COUNT(*) FROM tasks WHERE 1=1"
	var countArgs []interface{}
	argIndex := 1

	// Construir query para obtener tareas
	query := `SELECT id, title, description, status, priority, due_date, created_by, assigned_to, 
	          created_at, updated_at FROM tasks WHERE 1=1`
	var args []interface{}

	// Filtro por usuario (si existe)
	if userID != "" {
		countQuery += " AND (created_by = $" + fmt.Sprintf("%d", argIndex) + " OR assigned_to = $" + fmt.Sprintf("%d", argIndex) + ")"
		query += " AND (created_by = $" + fmt.Sprintf("%d", argIndex) + " OR assigned_to = $" + fmt.Sprintf("%d", argIndex) + ")"
		countArgs = append(countArgs, userID)
		args = append(args, userID)
		argIndex++
	}

	// Filtro por estado (si existe)
	if status != "" {
		countQuery += " AND status = $" + fmt.Sprintf("%d", argIndex)
		query += " AND status = $" + fmt.Sprintf("%d", argIndex)
		countArgs = append(countArgs, status)
		args = append(args, status)
		argIndex++
	}

	// Obtener total
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		log.Printf("🔴 ERROR en GetAll - Count Error: %v\n", err)
		return nil, 0, fmt.Errorf("error al contar tareas: %w", err)
	}

	// Agregar ordering, limit y offset
	query += " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", argIndex) + " OFFSET $" + fmt.Sprintf("%d", argIndex+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("🔴 ERROR en GetAll - QueryContext Error: %v (type: %T)\n", err, err)
		return nil, 0, fmt.Errorf("error al obtener tareas: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		var dueDate sql.NullTime
		var assignedTo sql.NullString

		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status,
			&task.Priority, &dueDate, &task.CreatedBy, &assignedTo,
			&task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			log.Printf("🔴 ERROR en GetAll - Scan Error: %v (type: %T)\n", err, err)
			return nil, 0, fmt.Errorf("error al escanear tarea: %w", err)
		}

		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}
		if assignedTo.Valid {
			task.AssignedTo = &assignedTo.String
		}

		tasks = append(tasks, task)
	}

	log.Printf("✅ GetAll - Success: returned %d tasks, totalCount=%d\n", len(tasks), totalCount)
	return tasks, totalCount, rows.Err()
}

// GetByID obtiene una tarea específica
func (r *TaskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	var task models.Task
	var dueDate sql.NullTime
	var assignedTo sql.NullString

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, title, description, status, priority, due_date, created_by, assigned_to, created_at, updated_at FROM tasks WHERE id = $1::UUID",
		id,
	).Scan(
		&task.ID, &task.Title, &task.Description, &task.Status,
		&task.Priority, &dueDate, &task.CreatedBy, &assignedTo,
		&task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tarea no encontrada")
		}
		return nil, fmt.Errorf("error al obtener tarea: %w", err)
	}

	if dueDate.Valid {
		task.DueDate = &dueDate.Time
	}
	if assignedTo.Valid {
		task.AssignedTo = &assignedTo.String
	}

	return &task, nil
}

// Create crea una nueva tarea
func (r *TaskRepository) Create(ctx context.Context, title, description, priority string, dueDate interface{}, createdBy string) (string, error) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("🔴 PANIC en Create: %v\n", rec)
		}
	}()

	var taskID string

	log.Printf("📝 Create - Input: title=%s, priority=%s, dueDate=%v, createdBy=%s\n", title, priority, dueDate, createdBy)

	// Insertar directamente en la tabla tasks
	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO tasks (title, description, priority, due_date, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		title, description, priority, dueDate, createdBy,
	).Scan(&taskID)

	if err != nil {
		log.Printf("🔴 ERROR en Create - Database Error: %v (type: %T)\n", err, err)
		return "", fmt.Errorf("error al crear tarea: %w", err)
	}

	log.Printf("✅ Task created successfully: %s\n", taskID)
	return taskID, nil
}

// Update actualiza una tarea existente
func (r *TaskRepository) Update(ctx context.Context, id, title, description, priority string, dueDate interface{}) error {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("🔴 PANIC en Update: %v\n", rec)
		}
	}()

	var dueDatePtr *time.Time
	if d, ok := dueDate.(*time.Time); ok {
		dueDatePtr = d
	}

	log.Printf("📝 Update - Input: id=%s, title=%s, priority=%s, dueDate=%v\n", id, title, priority, dueDatePtr)

	result, err := r.db.ExecContext(
		ctx,
		"UPDATE tasks SET title=$2, description=$3, priority=$4, due_date=$5, updated_at=CURRENT_TIMESTAMP WHERE id=$1::UUID",
		id, title, description, priority, dueDatePtr,
	)

	if err != nil {
		log.Printf("🔴 ERROR en Update - ExecContext Error: %v (type: %T)\n", err, err)
		return fmt.Errorf("error al actualizar tarea: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("🔴 ERROR en Update - RowsAffected Error: %v (type: %T)\n", err, err)
		return fmt.Errorf("error al actualizar tarea: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("🔴 ERROR en Update - Task not found\n")
		return fmt.Errorf("tarea no encontrada")
	}

	log.Printf("✅ Update - Success, rows affected: %d\n", rowsAffected)
	return nil
}

// Delete elimina una tarea
func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("🔴 PANIC en Delete: %v\n", rec)
		}
	}()

	log.Printf("📝 Delete - Input: id=%s\n", id)

	result, err := r.db.ExecContext(
		ctx,
		"DELETE FROM tasks WHERE id = $1::UUID",
		id,
	)

	if err != nil {
		log.Printf("🔴 ERROR en Delete - ExecContext Error: %v (type: %T)\n", err, err)
		return fmt.Errorf("error al eliminar tarea: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("🔴 ERROR en Delete - RowsAffected Error: %v (type: %T)\n", err, err)
		return fmt.Errorf("error al eliminar tarea: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("🔴 ERROR en Delete - Task not found\n")
		return fmt.Errorf("tarea no encontrada")
	}

	log.Printf("✅ Delete - Success, rows affected: %d\n", rowsAffected)
	return nil
}

// UpdateStatus actualiza el estado de una tarea
func (r *TaskRepository) UpdateStatus(ctx context.Context, id, status string) error {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("🔴 PANIC en UpdateStatus: %v\n", rec)
		}
	}()

	log.Printf("📝 UpdateStatus - Input: id=%s, status=%s\n", id, status)

	result, err := r.db.ExecContext(
		ctx,
		"UPDATE tasks SET status=$2, updated_at=CURRENT_TIMESTAMP WHERE id=$1::UUID",
		id, status,
	)

	if err != nil {
		log.Printf("🔴 ERROR en UpdateStatus - ExecContext Error: %v (type: %T)\n", err, err)
		return fmt.Errorf("error al actualizar estado de tarea: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("🔴 ERROR en UpdateStatus - RowsAffected Error: %v (type: %T)\n", err, err)
		return fmt.Errorf("error al actualizar estado de tarea: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("🔴 ERROR en UpdateStatus - Task not found\n")
		return fmt.Errorf("tarea no encontrada")
	}

	log.Printf("✅ UpdateStatus - Success, rows affected: %d\n", rowsAffected)
	return nil
}

// AssignTask asigna una tarea a un usuario
func (r *TaskRepository) AssignTask(ctx context.Context, taskID, userID string) error {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("🔴 PANIC en AssignTask: %v\n", rec)
		}
	}()

	log.Printf("📝 AssignTask - Input: taskID=%s, userID=%s\n", taskID, userID)

	// Permitir userID nil/empty para desasignar
	var assignedTo interface{} = userID
	if userID == "" {
		assignedTo = nil
	}

	result, err := r.db.ExecContext(
		ctx,
		"UPDATE tasks SET assigned_to=$2, updated_at=CURRENT_TIMESTAMP WHERE id=$1::UUID",
		taskID, assignedTo,
	)

	if err != nil {
		log.Printf("🔴 ERROR en AssignTask - ExecContext Error: %v (type: %T)\n", err, err)
		return fmt.Errorf("error al asignar tarea: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("🔴 ERROR en AssignTask - RowsAffected Error: %v (type: %T)\n", err, err)
		return fmt.Errorf("error al asignar tarea: %w", err)
	}

	if rowsAffected == 0 {
		log.Printf("🔴 ERROR en AssignTask - Task not found\n")
		return fmt.Errorf("tarea no encontrada")
	}

	log.Printf("✅ AssignTask - Success, rows affected: %d\n", rowsAffected)
	return nil
}

// GetStats obtiene estadísticas de tareas de un usuario
func (r *TaskRepository) GetStats(ctx context.Context, userID string) (*models.TaskStats, error) {
	var stats models.TaskStats

	err := r.db.QueryRowContext(
		ctx,
		`SELECT 
			COUNT(*),
			COUNT(*) FILTER (WHERE status = 'pending'),
			COUNT(*) FILTER (WHERE status = 'in_progress'),
			COUNT(*) FILTER (WHERE status = 'completed'),
			COUNT(*) FILTER (WHERE status = 'cancelled'),
			COUNT(*) FILTER (WHERE priority IN ('high', 'urgent')),
			COUNT(*) FILTER (WHERE due_date < CURRENT_TIMESTAMP AND status != 'completed')
		 FROM tasks
		 WHERE created_by = $1 OR assigned_to = $1`,
		userID,
	).Scan(
		&stats.TotalTasks, &stats.PendingCount, &stats.InProgressCount,
		&stats.CompletedCount, &stats.CancelledCount, &stats.HighPriorityCount,
		&stats.OverdueCount,
	)

	if err != nil {
		return nil, fmt.Errorf("error al obtener estadísticas: %w", err)
	}

	return &stats, nil
}
