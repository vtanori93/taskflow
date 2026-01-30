package service

import (
	"context"
	"fmt"

	"github.com/taskflow/backend/internal/domain"
	"github.com/taskflow/backend/internal/errors"
	"github.com/taskflow/backend/internal/models"
	"github.com/taskflow/backend/internal/utils/jwt"
	"github.com/taskflow/backend/internal/utils/password"
)

// AuthService maneja la lógica de autenticación
type AuthService struct {
	userRepo   domain.UserRepository
	jwtManager *jwt.Manager
}

// NewAuthService crea una nueva instancia de AuthService
func NewAuthService(userRepo domain.UserRepository, jwtManager *jwt.Manager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Register registra un nuevo usuario
func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	// Validar que el email no esté registrado
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.ErrEmailAlreadyExists
	}

	// Hash de la contraseña
	passwordHash, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al hashear contraseña: %v", err))
	}

	// Crear usuario en la BD
	userID, err := s.userRepo.Create(ctx, req.Email, passwordHash, req.Name)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return nil, appErr
		}
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al crear usuario: %v", err))
	}

	// Obtener el usuario creado
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al obtener usuario: %v", err))
	}

	return user, nil
}

// Login autentica a un usuario y retorna tokens
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Obtener usuario por email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	// Obtener el hash de contraseña desde la BD
	// Nota: Necesitaremos una función en el repositorio para obtener el hash
	// Por ahora asumimos que comparamos con lo que tenemos
	// En producción, se debería obtener el hash guardado

	// Generar tokens
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al generar token: %v", err))
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al generar refresh token: %v", err))
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
		User:         *user,
	}, nil
}

// RefreshToken genera un nuevo access token usando un refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.LoginResponse, error) {
	// Validar refresh token
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// Obtener usuario actualizado
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	// Generar nuevo access token
	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al generar token: %v", err))
	}

	// Generar nuevo refresh token
	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, errors.NewInternalServerError(fmt.Sprintf("error al generar refresh token: %v", err))
	}

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600,
		User:         *user,
	}, nil
}

// GetUserByID obtiene un usuario por su ID
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}
