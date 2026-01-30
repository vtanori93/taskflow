package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskflow/backend/internal/errors"
	"github.com/taskflow/backend/internal/infrastructure/response"
	"github.com/taskflow/backend/internal/models"
	"github.com/taskflow/backend/internal/service"
	"github.com/taskflow/backend/internal/utils/validation"
)

// TaskHandler maneja los endpoints de tareas
type TaskHandler struct {
	taskService    *service.TaskService
	responseWriter response.ResponseWriter
}

// NewTaskHandler crea una nueva instancia de TaskHandler
func NewTaskHandler(taskService *service.TaskService, rw response.ResponseWriter) *TaskHandler {
	return &TaskHandler{
		taskService:    taskService,
		responseWriter: rw,
	}
}

// CreateTask godoc
// @Summary Crear nueva tarea
// @Description Crea una nueva tarea para el usuario autenticado
// @Tags Tasks
// @Security Bearer
// @Accept json
// @Produce json
// @Param request body models.CreateTaskRequest true "Datos de la tarea"
// @Success 201 {object} models.APIResponse{data=models.Task}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /api/v1/tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("🔴 PANIC en CreateTask Handler: %v\n", rec)
			h.responseWriter.InternalError(c, "Critical error in task creation")
		}
	}()

	var req models.CreateTaskRequest

	fmt.Printf("📝 CreateTask Handler - Parsing JSON request\n")
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("🔴 ERROR en CreateTask Handler - JSON Parse Error: %v (type: %T)\n", err, err)
		h.responseWriter.ValidationError(c, err.Error())
		return
	}

	// Validar inputs
	if err := validation.ValidateTaskTitle(req.Title); err != nil {
		h.responseWriter.ValidationError(c, fmt.Sprintf("Título: %v", err))
		return
	}

	if req.Description != "" {
		if err := validation.ValidateTaskDescription(req.Description); err != nil {
			h.responseWriter.ValidationError(c, fmt.Sprintf("Descripción: %v", err))
			return
		}
	}

	if req.Priority != "" {
		if err := validation.ValidatePriority(req.Priority); err != nil {
			h.responseWriter.ValidationError(c, fmt.Sprintf("Prioridad: %v", err))
			return
		}
	}

	fmt.Printf("✅ JSON parsed: title=%s, priority=%s, dueDate=%v\n", req.Title, req.Priority, req.DueDate)

	userID, exists := c.Get("user_id")
	if !exists {
		fmt.Printf("🔴 ERROR en CreateTask Handler - No user_id in context\n")
		h.responseWriter.Unauthorized(c, "No autorizado")
		return
	}

	fmt.Printf("📝 CreateTask Handler - Calling service with userID=%v\n", userID)
	task, err := h.taskService.CreateTask(c.Request.Context(), &req, userID.(string))
	if err != nil {
		fmt.Printf("🔴 ERROR en CreateTask Handler - Service Error: %v (type: %T)\n", err, err)
		h.handleError(c, err)
		return
	}

	fmt.Printf("✅ CreateTask Handler - Success, returning task: %+v\n", task)
	h.responseWriter.Success(c, http.StatusCreated, "Tarea creada exitosamente", task)
}

// GetTasks godoc
// @Summary Listar todas las tareas
// @Description Obtiene todas las tareas del sistema con paginación (sin filtro por usuario)
// @Tags Tasks
// @Security Bearer
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param page_size query int false "Tamaño de página" default(20)
// @Param status query string false "Filtrar por estado" Enums(pending,in_progress,completed,cancelled)
// @Success 200 {object} models.APIResponse{data=models.TasksListResponse}
// @Failure 401 {object} models.APIResponse
// @Router /api/v1/tasks [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("🔴 PANIC en GetTasks Handler: %v\n", rec)
			h.responseWriter.InternalError(c, "Critical error in get tasks")
		}
	}()

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 20
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	status := c.Query("status")

	fmt.Printf("📝 GetTasks Handler - Calling service with page=%d, pageSize=%d, status=%s (ALL TASKS)\n", page, pageSize, status)

	// Pasar userID vacío para obtener todas las tareas
	resp, err := h.taskService.GetTasks(c.Request.Context(), "", status, page, pageSize)
	if err != nil {
		fmt.Printf("🔴 ERROR en GetTasks Handler - Service Error: %v (type: %T)\n", err, err)
		h.handleError(c, err)
		return
	}

	fmt.Printf("✅ GetTasks Handler - Success, returning %d tasks\n", len(resp.Tasks))
	h.responseWriter.Success(c, http.StatusOK, "Tareas obtenidas exitosamente", resp)
}

// GetMyTasks godoc
// @Summary Listar mis tareas
// @Description Obtiene las tareas creadas por o asignadas al usuario autenticado con paginación
// @Tags Tasks
// @Security Bearer
// @Produce json
// @Param page query int false "Número de página" default(1)
// @Param page_size query int false "Tamaño de página" default(20)
// @Param status query string false "Filtrar por estado" Enums(pending,in_progress,completed,cancelled)
// @Success 200 {object} models.APIResponse{data=models.TasksListResponse}
// @Failure 401 {object} models.APIResponse
// @Router /api/v1/tasks/my [get]
func (h *TaskHandler) GetMyTasks(c *gin.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("🔴 PANIC en GetMyTasks Handler: %v\n", rec)
			h.responseWriter.InternalError(c, "Critical error in get my tasks")
		}
	}()

	userID, exists := c.Get("user_id")
	if !exists {
		fmt.Printf("🔴 ERROR en GetMyTasks Handler - No user_id in context\n")
		h.responseWriter.Unauthorized(c, "No autorizado")
		return
	}

	fmt.Printf("📝 GetMyTasks Handler - userID=%v\n", userID)

	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 20
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	status := c.Query("status")

	fmt.Printf("📝 GetMyTasks Handler - Calling service with page=%d, pageSize=%d, status=%s\n", page, pageSize, status)

	resp, err := h.taskService.GetTasks(c.Request.Context(), userID.(string), status, page, pageSize)
	if err != nil {
		fmt.Printf("🔴 ERROR en GetMyTasks Handler - Service Error: %v (type: %T)\n", err, err)
		h.handleError(c, err)
		return
	}

	fmt.Printf("✅ GetMyTasks Handler - Success, returning %d tasks\n", len(resp.Tasks))
	h.responseWriter.Success(c, http.StatusOK, "Mis tareas obtenidas exitosamente", resp)
}

// GetTask godoc
// @Summary Obtener tarea por ID
// @Description Obtiene los detalles de una tarea específica
// @Tags Tasks
// @Security Bearer
// @Produce json
// @Param id path string true "ID de la tarea"
// @Success 200 {object} models.APIResponse{data=models.Task}
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID := c.Param("id")

	if taskID == "" {
		h.responseWriter.ValidationError(c, "ID de tarea requerido")
		return
	}

	// Validar que sea un UUID válido
	if err := validation.ValidateUUID(taskID); err != nil {
		h.responseWriter.ValidationError(c, fmt.Sprintf("ID inválido: %v", err))
		return
	}

	task, err := h.taskService.GetTaskByID(c.Request.Context(), taskID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.responseWriter.Success(c, http.StatusOK, "Tarea obtenida exitosamente", task)
}

// UpdateTask godoc
// @Summary Actualizar tarea
// @Description Actualiza los detalles de una tarea
// @Tags Tasks
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "ID de la tarea"
// @Param request body models.UpdateTaskRequest true "Datos a actualizar"
// @Success 200 {object} models.APIResponse{data=models.Task}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")

	if taskID == "" {
		h.responseWriter.ValidationError(c, "ID de tarea requerido")
		return
	}

	// Validar UUID
	if err := validation.ValidateUUID(taskID); err != nil {
		h.responseWriter.ValidationError(c, fmt.Sprintf("ID inválido: %v", err))
		return
	}

	var req models.UpdateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.responseWriter.ValidationError(c, err.Error())
		return
	}

	// Validar campos opcionales
	if req.Title != nil && *req.Title != "" {
		if err := validation.ValidateTaskTitle(*req.Title); err != nil {
			h.responseWriter.ValidationError(c, fmt.Sprintf("Título: %v", err))
			return
		}
	}

	if req.Description != nil && *req.Description != "" {
		if err := validation.ValidateTaskDescription(*req.Description); err != nil {
			h.responseWriter.ValidationError(c, fmt.Sprintf("Descripción: %v", err))
			return
		}
	}

	if req.Priority != nil && *req.Priority != "" {
		if err := validation.ValidatePriority(*req.Priority); err != nil {
			h.responseWriter.ValidationError(c, fmt.Sprintf("Prioridad: %v", err))
			return
		}
	}

	task, err := h.taskService.UpdateTask(c.Request.Context(), taskID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.responseWriter.Success(c, http.StatusOK, "Tarea actualizada exitosamente", task)
}

// DeleteTask godoc
// @Summary Eliminar tarea
// @Description Elimina una tarea del sistema
// @Tags Tasks
// @Security Bearer
// @Produce json
// @Param id path string true "ID de la tarea"
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("🔴 PANIC en DeleteTask Handler: %v\n", rec)
			h.responseWriter.InternalError(c, "Critical error in delete task")
		}
	}()

	taskID := c.Param("id")

	if taskID == "" {
		fmt.Printf("🔴 ERROR en DeleteTask Handler - ID vacío\n")
		h.responseWriter.ValidationError(c, "ID de tarea requerido")
		return
	}

	fmt.Printf("📝 DeleteTask Handler - taskID=%s\n", taskID)

	err := h.taskService.DeleteTask(c.Request.Context(), taskID)
	if err != nil {
		fmt.Printf("🔴 ERROR en DeleteTask Handler - Service Error: %v (type: %T)\n", err, err)
		h.handleError(c, err)
		return
	}

	fmt.Printf("✅ DeleteTask Handler - Success\n")
	h.responseWriter.Success(c, http.StatusOK, "Tarea eliminada exitosamente", nil)
}

// UpdateTaskStatus godoc
// @Summary Actualizar estado de tarea
// @Description Cambia el estado de una tarea
// @Tags Tasks
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "ID de la tarea"
// @Param request body models.UpdateTaskStatusRequest true "Nuevo estado"
// @Success 200 {object} models.APIResponse{data=models.Task}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/tasks/{id}/status [patch]
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	taskID := c.Param("id")

	if taskID == "" {
		h.responseWriter.ValidationError(c, "ID de tarea requerido")
		return
	}

	// Validar UUID
	if err := validation.ValidateUUID(taskID); err != nil {
		h.responseWriter.ValidationError(c, fmt.Sprintf("ID inválido: %v", err))
		return
	}

	var req models.UpdateTaskStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.responseWriter.ValidationError(c, err.Error())
		return
	}

	// Validar estado
	if err := validation.ValidateStatus(req.Status); err != nil {
		h.responseWriter.ValidationError(c, fmt.Sprintf("Estado: %v", err))
		return
	}

	task, err := h.taskService.UpdateTaskStatus(c.Request.Context(), taskID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.responseWriter.Success(c, http.StatusOK, "Estado de tarea actualizado exitosamente", task)
}

// AssignTask godoc
// @Summary Asignar tarea a usuario
// @Description Asigna una tarea a un usuario específico
// @Tags Tasks
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "ID de la tarea"
// @Param request body models.AssignTaskRequest true "Usuario a asignar"
// @Success 200 {object} models.APIResponse{data=models.Task}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/tasks/{id}/assign [post]
func (h *TaskHandler) AssignTask(c *gin.Context) {
	taskID := c.Param("id")

	if taskID == "" {
		h.responseWriter.ValidationError(c, "ID de tarea requerido")
		return
	}

	var req models.AssignTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.responseWriter.ValidationError(c, err.Error())
		return
	}

	task, err := h.taskService.AssignTask(c.Request.Context(), taskID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.responseWriter.Success(c, http.StatusOK, "Tarea asignada exitosamente", task)
}

// GetTaskStats godoc
// @Summary Obtener estadísticas de tareas
// @Description Obtiene estadísticas de tareas del usuario autenticado
// @Tags Tasks
// @Security Bearer
// @Produce json
// @Success 200 {object} models.APIResponse{data=models.TaskStats}
// @Failure 401 {object} models.APIResponse
// @Router /api/v1/tasks/stats [get]
func (h *TaskHandler) GetTaskStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.responseWriter.Unauthorized(c, "No autorizado")
		return
	}

	stats, err := h.taskService.GetTaskStats(c.Request.Context(), userID.(string))
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.responseWriter.Success(c, http.StatusOK, "Estadísticas obtenidas exitosamente", stats)
}

// handleError maneja los errores de la aplicación
func (h *TaskHandler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		h.responseWriter.Error(c, appErr.Code, appErr.Message)
		return
	}

	h.responseWriter.InternalError(c, err.Error())
}
