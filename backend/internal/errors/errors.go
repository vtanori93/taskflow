package errors

import "fmt"

// AppError representa un error de la aplicación
type AppError struct {
	Code    int
	Message string
	Details string
}

// Error implementa la interfaz error
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// Errores predefinidos
var (
	ErrInvalidCredentials = &AppError{
		Code:    401,
		Message: "Credenciales inválidas",
	}

	ErrUserNotFound = &AppError{
		Code:    404,
		Message: "Usuario no encontrado",
	}

	ErrTaskNotFound = &AppError{
		Code:    404,
		Message: "Tarea no encontrada",
	}

	ErrEmailAlreadyExists = &AppError{
		Code:    409,
		Message: "El email ya está registrado",
	}

	ErrUnauthorized = &AppError{
		Code:    401,
		Message: "No autorizado",
	}

	ErrInvalidToken = &AppError{
		Code:    401,
		Message: "Token inválido o expirado",
	}

	ErrInternalServer = &AppError{
		Code:    500,
		Message: "Error interno del servidor",
	}

	ErrBadRequest = &AppError{
		Code:    400,
		Message: "Solicitud inválida",
	}
)

// NewAppError crea un nuevo AppError
func NewAppError(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// NewInternalServerError crea un error de servidor interno
func NewInternalServerError(details string) *AppError {
	return &AppError{
		Code:    500,
		Message: "Error interno del servidor",
		Details: details,
	}
}

// NewBadRequest crea un error de solicitud inválida
func NewBadRequest(details string) *AppError {
	return &AppError{
		Code:    400,
		Message: "Solicitud inválida",
		Details: details,
	}
}
