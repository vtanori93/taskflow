package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Connect establece una conexión con la base de datos
func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión: %w", err)
	}

	// Verificar que la conexión sea válida
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error al conectar a BD: %w", err)
	}

	// Configurar conexión pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

// Close cierra la conexión con la base de datos
func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
