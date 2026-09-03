package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Syipmong/acoustic-sensor-backend/internal/domain"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository{
	return &PostgresRepository{db:db}
}



