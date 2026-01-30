package service

import (
	"context"
	"fmt"
	"log"

	"github.com/taskflow/backend/internal/domain"
	"github.com/taskflow/backend/internal/errors"
	"github.com/taskflow/backend/internal/models"
)

// UserService maneja la lógica de negocio de usuarios
type UserService struct {
	userRepo domain.UserRepository
}

// NewUserService crea una nueva instancia de UserService
func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetAllUsers obtiene todos los usuarios registrados en el sistema
func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("🔴 PANIC en GetAllUsers: %v\n", rec)
		}
	}()

	log.Printf("📝 GetAllUsers Service - Fetching all users\n")

	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		log.Printf("🔴 ERROR en GetAllUsers Service - Repository Error: %v (type: %T)\n", err, err)
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al obtener usuarios: %v", err))
	}

	log.Printf("✅ GetAllUsers Service - Success: returned %d users\n", len(users))
	return users, nil
}
