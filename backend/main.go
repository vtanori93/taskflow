package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/taskflow/backend/docs"
	"github.com/taskflow/backend/internal/config"
	"github.com/taskflow/backend/internal/handler"
	"github.com/taskflow/backend/internal/infrastructure/database"
	"github.com/taskflow/backend/internal/infrastructure/response"
	"github.com/taskflow/backend/internal/infrastructure/router"
	"github.com/taskflow/backend/internal/repository/postgres"
	"github.com/taskflow/backend/internal/service"
	"github.com/taskflow/backend/internal/utils/jwt"
)

// @title TaskFlow API
// @version 1.0
// @description API REST para gestión de tareas colaborativas
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@taskflow.local

// @license.name MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT Token en formato "Bearer {token}"

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// Establecer modo de Gin
	if cfg.ServerEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Conectar a base de datos
	db, err := database.Connect(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Error conectando a BD: %v", err)
	}
	defer database.Close(db)

	fmt.Println("✅ Conectado a PostgreSQL")

	// Crear JWT Manager
	jwtManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTExpirationTime, cfg.JWTRefreshExpiration)

	// Crear ResponseWriter para respuestas estandarizadas
	rw := response.NewResponseWriter()

	// Crear repositorios
	userRepo := postgres.NewUserRepository(db)
	taskRepo := postgres.NewTaskRepository(db)

	// Crear servicios
	authService := service.NewAuthService(userRepo, jwtManager)
	taskService := service.NewTaskService(taskRepo)
	userService := service.NewUserService(userRepo)

	// Crear handlers con inyección de ResponseWriter
	authHandler := handler.NewAuthHandler(authService, rw)
	taskHandler := handler.NewTaskHandler(taskService, rw)
	userHandler := handler.NewUserHandler(userService, rw)

	// Crear engine de Gin
	engine := gin.Default()

	// Setup de rutas
	router.Setup(engine, authHandler, taskHandler, userHandler, jwtManager)

	// Iniciar servidor
	addr := fmt.Sprintf(":%d", cfg.ServerPort)

	fmt.Printf("\n🚀 API iniciada en http://localhost%s\n", addr)
	fmt.Printf("📚 Documentación Swagger: http://localhost%s/swagger/index.html\n\n", addr)

	if err := engine.Run(addr); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
