package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
	"github.com/taskflow/backend/internal/domain"
	"github.com/taskflow/backend/internal/errors"
	"github.com/taskflow/backend/internal/models"
)

// UserRepository implementa domain.UserRepository usando PostgreSQL
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository crea una nueva instancia de UserRepository
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

// GetByEmail obtiene un usuario por email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, email, name, created_at, updated_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}

	return &user, nil
}

// GetByID obtiene un usuario por ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1::UUID",
		id,
	).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, fmt.Errorf("error al obtener usuario: %w", err)
	}

	return &user, nil
}

// Create crea un nuevo usuario
func (r *UserRepository) Create(ctx context.Context, email, passwordHash, name string) (string, error) {
	var userID string

	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3) RETURNING id",
		email, passwordHash, name,
	).Scan(&userID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return "", errors.ErrEmailAlreadyExists
		}
		return "", fmt.Errorf("error al registrar usuario: %w", err)
	}

	return userID, nil
}

// Update actualiza un usuario
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	// Implementación simple - puede extenderse según necesidades
	return nil
}

// GetAllUsers obtiene todos los usuarios registrados en el sistema
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("🔴 PANIC en GetAllUsers: %v\n", rec)
		}
	}()

	log.Printf("📝 GetAllUsers - Fetching all users\n")

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, email, name, created_at, updated_at 
		 FROM users 
		 ORDER BY created_at DESC`,
	)

	if err != nil {
		log.Printf("🔴 ERROR en GetAllUsers - QueryContext Error: %v (type: %T)\n", err, err)
		return nil, fmt.Errorf("error al obtener usuarios: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Printf("🔴 ERROR en GetAllUsers - Scan Error: %v\n", err)
			return nil, fmt.Errorf("error al escanear usuario: %w", err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("🔴 ERROR en GetAllUsers - Rows Error: %v\n", err)
		return nil, fmt.Errorf("error iterando usuarios: %w", err)
	}

	log.Printf("✅ GetAllUsers - Success: returned %d users\n", len(users))
	return users, nil
}
