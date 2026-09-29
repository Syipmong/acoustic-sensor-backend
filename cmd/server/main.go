package main

import(
	"crypto/ecdsa"
	"crypto/elliptic"
	"database/sql"
	"log"
	"net/http"
	"os"

	deliveryhttp "acoustic-sensor-backend/internal/delivery/http"
	"acoustic-sensor-backend/internal/repository"
	"acoustic-sensor-backend/internal/usecase"


	- "github.com/lib/pq"
)

func loadNodePublicKey() *ecdsa.PublicKey{
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
	}
}

func main(){
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == ""{
		dsn = "postgres://postgres:password@localhost:5432/aetech_forensics?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("PostgresSQL connection pool established")

	nodePubKey := loadNodePublicKey()

	alertRepo := repository.NewPostgresRepository(db)
	alertUsecase := usecase.NewAlertUsecase(alertRepo, nodePubKey)
	alertHandler := deliveryhttp.newAlertHandler(alertUsecase)
	
}
