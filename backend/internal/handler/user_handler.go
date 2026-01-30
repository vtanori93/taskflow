package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taskflow/backend/internal/infrastructure/response"
	"github.com/taskflow/backend/internal/service"
)

// UserHandler maneja los endpoints de usuarios
type UserHandler struct {
	userService    *service.UserService
	responseWriter response.ResponseWriter
}

// NewUserHandler crea una nueva instancia de UserHandler
func NewUserHandler(userService *service.UserService, rw response.ResponseWriter) *UserHandler {
	return &UserHandler{
		userService:    userService,
		responseWriter: rw,
	}
}

// GetAllUsers godoc
// @Summary Listar todos los usuarios
// @Description Obtiene el listado de todos los usuarios registrados en el sistema
// @Tags Users
// @Security Bearer
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.User}
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("🔴 PANIC en GetAllUsers Handler: %v\n", rec)
			h.responseWriter.InternalError(c, "Critical error fetching users")
		}
	}()

	fmt.Printf("📝 GetAllUsers Handler - Fetching all users\n")

	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		fmt.Printf("🔴 ERROR en GetAllUsers Handler - Service Error: %v (type: %T)\n", err, err)
		h.responseWriter.InternalError(c, "Error al obtener usuarios")
		return
	}

	fmt.Printf("✅ GetAllUsers Handler - Success, returning %d users\n", len(users))
	h.responseWriter.Success(c, http.StatusOK, "Usuarios obtenidos exitosamente", gin.H{
		"users": users,
		"count": len(users),
	})
}
