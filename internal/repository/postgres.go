package repository

import (
	"database/sql"
	"fmt"

	"acoustic-sensor-backend/internal/domain"
	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(connectionString string)