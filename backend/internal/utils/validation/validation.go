package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// ValidateUUID valida que una cadena sea un UUID válido
func ValidateUUID(id string) error {
	if id == "" {
		return fmt.Errorf("ID no puede estar vacío")
	}
	_, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("ID inválido: no es un UUID válido")
	}
	return nil
}

// ValidateStatus valida que el estado sea uno de los permitidos
func ValidateStatus(status string) error {
	validStatuses := map[string]bool{
		"pending":     true,
		"in_progress": true,
		"completed":   true,
		"cancelled":   true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("estado inválido: debe ser uno de pending, in_progress, completed, cancelled")
	}
	return nil
}

// ValidatePriority valida que la prioridad sea una de las permitidas
func ValidatePriority(priority string) error {
	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
		"urgent": true,
	}

	if !validPriorities[priority] {
		return fmt.Errorf("prioridad inválida: debe ser una de low, medium, high, urgent")
	}
	return nil
}

// ValidateEmail valida que sea un email válido
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email no puede estar vacío")
	}

	// Regex básico para email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("email inválido")
	}

	if len(email) > 255 {
		return fmt.Errorf("email muy largo (máximo 255 caracteres)")
	}

	return nil
}

// ValidatePassword valida que la contraseña cumpla los requisitos
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("contraseña debe tener al menos 8 caracteres")
	}

	if len(password) > 128 {
		return fmt.Errorf("contraseña muy larga (máximo 128 caracteres)")
	}

	return nil
}

// ValidateString valida una cadena de texto
func ValidateString(value string, minLen, maxLen int, fieldName string) error {
	if minLen > 0 && len(value) < minLen {
		return fmt.Errorf("%s debe tener al menos %d caracteres", fieldName, minLen)
	}

	if maxLen > 0 && len(value) > maxLen {
		return fmt.Errorf("%s debe tener máximo %d caracteres", fieldName, maxLen)
	}

	return nil
}

// SanitizeString elimina espacios y caracteres de control
func SanitizeString(value string) string {
	// Eliminar espacios al inicio y final
	value = strings.TrimSpace(value)

	// Eliminar caracteres de control
	value = strings.Map(func(r rune) rune {
		if r < 32 && r != '\n' && r != '\t' {
			return -1 // Eliminar
		}
		return r
	}, value)

	return value
}

// ValidateSQLIdentifier valida que sea un identificador SQL seguro (para column names, etc)
func ValidateSQLIdentifier(identifier string) error {
	if identifier == "" {
		return fmt.Errorf("identificador no puede estar vacío")
	}

	// Solo permitir alphanumericos y underscore
	sqlIdentifierRegex := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	if !sqlIdentifierRegex.MatchString(identifier) {
		return fmt.Errorf("identificador SQL inválido: solo se permiten letras, números y underscore")
	}

	return nil
}

// ValidatePaginationParams valida los parámetros de paginación
func ValidatePaginationParams(page, pageSize int) (int, int, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 20
	}

	if pageSize > 100 {
		pageSize = 100 // Máximo 100 registros por página
	}

	return page, pageSize, nil
}

// ValidateTaskTitle valida el título de una tarea
func ValidateTaskTitle(title string) error {
	title = strings.TrimSpace(title)

	if len(title) == 0 {
		return fmt.Errorf("título no puede estar vacío")
	}

	if len(title) > 200 {
		return fmt.Errorf("título debe tener máximo 200 caracteres")
	}

	return nil
}

// ValidateTaskDescription valida la descripción de una tarea
func ValidateTaskDescription(description string) error {
	if len(description) > 2000 {
		return fmt.Errorf("descripción debe tener máximo 2000 caracteres")
	}

	return nil
}
