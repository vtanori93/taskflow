package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taskflow/backend/internal/errors"
	"github.com/taskflow/backend/internal/infrastructure/response"
	"github.com/taskflow/backend/internal/models"
	"github.com/taskflow/backend/internal/service"
)

// AuthHandler maneja los endpoints de autenticación
type AuthHandler struct {
	authService    *service.AuthService
	responseWriter response.ResponseWriter
}

// NewAuthHandler crea una nueva instancia de AuthHandler
func NewAuthHandler(authService *service.AuthService, rw response.ResponseWriter) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		responseWriter: rw,
	}
}

// Register godoc
// @Summary Registrar nuevo usuario
// @Description Registra un nuevo usuario en el sistema
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Datos de registro"
// @Success 201 {object} models.APIResponse{data=models.User}
// @Failure 400 {object} models.APIResponse
// @Failure 409 {object} models.APIResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.responseWriter.ValidationError(c, err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.responseWriter.Success(c, http.StatusCreated, "Usuario registrado exitosamente", user)
}

// Login godoc
// @Summary Login de usuario
// @Description Autentica a un usuario y retorna tokens JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Credenciales"
// @Success 200 {object} models.APIResponse{data=models.LoginResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.responseWriter.ValidationError(c, err.Error())
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.responseWriter.Success(c, http.StatusOK, "Login exitoso", resp)
}

// RefreshToken godoc
// @Summary Refrescar token de acceso
// @Description Genera un nuevo access token usando un refresh token válido
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.APIResponse{data=models.LoginResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.responseWriter.ValidationError(c, err.Error())
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.responseWriter.Success(c, http.StatusOK, "Token refrescado exitosamente", resp)
}

// GetProfile godoc
// @Summary Obtener perfil del usuario actual
// @Description Obtiene la información del usuario autenticado
// @Tags Auth
// @Security Bearer
// @Produce json
// @Success 200 {object} models.APIResponse{data=models.User}
// @Failure 401 {object} models.APIResponse
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		h.responseWriter.Unauthorized(c, "No autorizado")
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.responseWriter.Success(c, http.StatusOK, "Perfil obtenido exitosamente", user)
}

// handleError maneja los errores de la aplicación
func (h *AuthHandler) handleError(c *gin.Context, err error) {
	// Importar el paquete errors
	appErr, ok := err.(*errors.AppError)
	if !ok {
		h.responseWriter.InternalError(c, "Error interno del servidor")
		return
	}

	h.responseWriter.Error(c, appErr.Code, appErr.Message)
}
