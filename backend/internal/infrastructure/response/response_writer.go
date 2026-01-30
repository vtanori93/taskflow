package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taskflow/backend/internal/models"
)

// ResponseWriter interfaz para escribir respuestas
type ResponseWriter interface {
	Success(c *gin.Context, statusCode int, message string, data interface{})
	Error(c *gin.Context, statusCode int, message string)
	ValidationError(c *gin.Context, message string)
	Unauthorized(c *gin.Context, message string)
	Forbidden(c *gin.Context, message string)
	NotFound(c *gin.Context, message string)
	InternalError(c *gin.Context, message string)
}

// StandardResponseWriter implementación de ResponseWriter
type StandardResponseWriter struct{}

// NewResponseWriter crea una nueva instancia de ResponseWriter
func NewResponseWriter() ResponseWriter {
	return &StandardResponseWriter{}
}

// Success envía una respuesta exitosa
func (w *StandardResponseWriter) Success(c *gin.Context, statusCode int, message string, data interface{}) {
	response := models.APIResponse{
		Data:       data,
		StatusCode: statusCode,
		Message:    message,
		Error:      false,
	}
	c.JSON(statusCode, response)
}

// Error envía una respuesta de error
func (w *StandardResponseWriter) Error(c *gin.Context, statusCode int, message string) {
	response := models.APIResponse{
		Data:       nil,
		StatusCode: statusCode,
		Message:    message,
		Error:      true,
	}
	c.JSON(statusCode, response)
}

// ValidationError envía error de validación (400)
func (w *StandardResponseWriter) ValidationError(c *gin.Context, message string) {
	w.Error(c, http.StatusBadRequest, message)
}

// Unauthorized envía error no autorizado (401)
func (w *StandardResponseWriter) Unauthorized(c *gin.Context, message string) {
	w.Error(c, http.StatusUnauthorized, message)
}

// Forbidden envía error prohibido (403)
func (w *StandardResponseWriter) Forbidden(c *gin.Context, message string) {
	w.Error(c, http.StatusForbidden, message)
}

// NotFound envía error no encontrado (404)
func (w *StandardResponseWriter) NotFound(c *gin.Context, message string) {
	w.Error(c, http.StatusNotFound, message)
}

// InternalError envía error interno (500)
func (w *StandardResponseWriter) InternalError(c *gin.Context, message string) {
	w.Error(c, http.StatusInternalServerError, message)
}
