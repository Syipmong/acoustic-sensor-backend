package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"database/sql"
	"log"
	"net/http"
	"os"

	deliveryhttp "github.com/Syipmong/acoustic-sensor-backend/internal/delivery/http"
	"github.com/Syipmong/acoustic-sensor-backend/internal/repository"
	"github.com/Syipmong/acoustic-sensor-backend/internal/usecase"

	_ "github.com/lib/pq"
)

func loadNodePublicKey() *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
	}
}

func main() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://9893f8bb05046918e126e05550985e761362a231e630c8b1c41e6fff1bf3666b:sk_4pvadnSaJ7GZESuOUwwJ7@pooled.db.prisma.io:5432/postgres?sslmode=require"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("PostgresSQL connection pool established")

	nodePubKey := loadNodePublicKey()

	alertRepo := repository.NewPostgresRepository(db)
	alertUsecase := usecase.NewAlertUsecase(alertRepo, nodePubKey)
	alertHandler := deliveryhttp.NewAlertHandler(alertUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/ingest", alertHandler.IngestLoraPacket)
	mux.HandleFunc("/ws/alerts", alertHandler.MobileWebsocketEndpoint)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Sensor Server Active on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server crashed: %v", err)
	}
}
