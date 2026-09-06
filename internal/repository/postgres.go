package repository

import (
	"context"
	"database/sql"

	"github.com/Syipmong/acoustic-sensor-backend/internal/domain"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository{
	return &PostgresRepository{db:db}
}

//Log Alert to insert a veried forensic into the table
func (r *PostgresRepository) LogAlert(ctx context.Context, alert *domain.Alert) error{

	query := `
	INSERT INTO acoustic_alerts
	(node_id, epoch_time, sequence_num, latitude, longitude, threat_class, confidence, battery_volts, signature_r_s)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, created_at`

	err := r.db.QueryRowContext(
		ctx, query, alert.NodeID,
		alert.EpochTime,
		alert.SequenceNum,
		alert.Latitude,
		alert.Longitude,
		alert.ThreatClass,
		alert.Confidence,
		alert.BatteryVolts,
		alert.SignatureRS,
	).Scan(&alert.ID, &alert.CreatedAt)
	return err
}