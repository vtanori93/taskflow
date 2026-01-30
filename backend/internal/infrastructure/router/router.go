package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/taskflow/backend/internal/handler"
	"github.com/taskflow/backend/internal/infrastructure/response"
	"github.com/taskflow/backend/internal/middleware"
	"github.com/taskflow/backend/internal/utils/jwt"
)

// Setup configura todas las rutas de la API
func Setup(
	engine *gin.Engine,
	authHandler *handler.AuthHandler,
	taskHandler *handler.TaskHandler,
	userHandler *handler.UserHandler,
	jwtManager *jwt.Manager,
) {
	// Middleware global
	rw := response.NewResponseWriter()
	engine.Use(middleware.CORSMiddleware())
	engine.Use(middleware.ErrorHandlingMiddleware())
	engine.Use(middleware.ValidationMiddleware(rw))
	engine.Use(middleware.SanitizeQueryParams())
	engine.Use(middleware.ValidateURLEncoding())

	// Swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rutas públicas
	api := engine.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		health := api.Group("/health")
		{
			health.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "ok"})
			})
		}
	}

	// Rutas protegidas
	protected := engine.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		// Auth routes
		auth := protected.Group("/auth")
		{
			auth.GET("/profile", authHandler.GetProfile)
		}

		// Task routes
		tasks := protected.Group("/tasks")
		{
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("", taskHandler.GetTasks)
			tasks.GET("/my", taskHandler.GetMyTasks)
			tasks.GET("/stats", taskHandler.GetTaskStats)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
			tasks.PATCH("/:id/status", taskHandler.UpdateTaskStatus)
			tasks.POST("/:id/assign", taskHandler.AssignTask)
		}

		// User routes
		users := protected.Group("/users")
		{
			users.GET("", userHandler.GetAllUsers)
		}
	}
}
