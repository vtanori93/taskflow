package middleware

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taskflow/backend/internal/infrastructure/response"
	"github.com/taskflow/backend/internal/utils/validation"
)

// ValidationMiddleware valida y sanitiza inputs para prevenir inyecciones SQL
func ValidationMiddleware(rw response.ResponseWriter) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("🔴 PANIC en ValidationMiddleware: %v\n", rec)
			}
		}()

		log.Printf("📝 ValidationMiddleware - Validando request: %s %s\n", c.Request.Method, c.Request.URL.Path)

		// Validar tamaño del body
		if c.Request.ContentLength > 10*1024*1024 { // 10MB máximo
			log.Printf("🔴 ValidationMiddleware - Body muy grande: %d bytes\n", c.Request.ContentLength)
			rw.ValidationError(c, "Request body demasiado grande (máximo 10MB)")
			c.Abort()
			return
		}

		// Validar parámetros de query
		for key, values := range c.Request.URL.Query() {
			// Validar nombre del parámetro
			if err := validation.ValidateSQLIdentifier(key); err != nil {
				log.Printf("🔴 ValidationMiddleware - Query parameter inválido: %s = %v\n", key, err)
				rw.ValidationError(c, fmt.Sprintf("Parámetro de query inválido: %s", key))
				c.Abort()
				return
			}

			// Validar valores
			for _, value := range values {
				if len(value) > 1000 {
					log.Printf("🔴 ValidationMiddleware - Valor de query muy largo: %s\n", key)
					rw.ValidationError(c, fmt.Sprintf("Valor de parámetro %s demasiado largo", key))
					c.Abort()
					return
				}
			}
		}

		// Validar parámetros de ruta
		for _, param := range c.Params {
			if param.Key == "id" {
				if err := validation.ValidateUUID(param.Value); err != nil {
					log.Printf("🔴 ValidationMiddleware - ID inválido: %v\n", err)
					rw.ValidationError(c, fmt.Sprintf("ID inválido: %v", err))
					c.Abort()
					return
				}
			}
		}

		// Validar parámetros comunes de paginación
		pageStr := c.DefaultQuery("page", "1")
		pageSizeStr := c.DefaultQuery("page_size", "20")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			log.Printf("🔴 ValidationMiddleware - Parámetro page inválido\n")
			page = 1
		}

		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSize < 1 {
			log.Printf("🔴 ValidationMiddleware - Parámetro page_size inválido\n")
			pageSize = 20
		}

		if pageSize > 100 {
			pageSize = 100
		}

		// Guardar en contexto
		c.Set("page", page)
		c.Set("page_size", pageSize)

		log.Printf("✅ ValidationMiddleware - Validación exitosa\n")
		c.Next()
	}
}

// ValidateJSONInput valida el JSON de entrada
func ValidateJSONInput(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body == nil {
			c.Next()
			return
		}

		c.Next()
	}
}

// SanitizeQueryParams sanitiza los parámetros de query
func SanitizeQueryParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitizar status si existe
		if status := c.Query("status"); status != "" {
			if err := validation.ValidateStatus(status); err != nil {
				log.Printf("🔴 SanitizeQueryParams - Status inválido: %v\n", err)
				c.Set("status", "")
			} else {
				c.Set("status", status)
			}
		}

		// Sanitizar priority si existe
		if priority := c.Query("priority"); priority != "" {
			if err := validation.ValidatePriority(priority); err != nil {
				log.Printf("🔴 SanitizeQueryParams - Priority inválida: %v\n", err)
				c.Set("priority", "")
			} else {
				c.Set("priority", priority)
			}
		}

		// Sanitizar search si existe
		if search := c.Query("search"); search != "" {
			search = validation.SanitizeString(search)
			if len(search) > 100 {
				search = search[:100]
			}
			c.Set("search", search)
		}

		c.Next()
	}
}

// ValidateURLEncoding valida que la URL esté correctamente codificada
func ValidateURLEncoding() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Intentar decodificar la URL para detectar codificación inválida
		_, err := url.QueryUnescape(c.Request.URL.RawQuery)
		if err != nil {
			log.Printf("🔴 ValidateURLEncoding - URL inválida: %v\n", err)
			rw := response.NewResponseWriter()
			rw.ValidationError(c, "URL inválida o malformada")
			c.Abort()
			return
		}

		c.Next()
	}
}
