package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/darkn3sssss/sar-dozor/internal/db"
	"github.com/darkn3sssss/sar-dozor/internal/objects"
	"github.com/joho/godotenv"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("не удалось загрузить .env, использую переменные окружения как есть")
	}
	dsn := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	pool, err := db.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("подключились к базе")
	defer pool.Close()
	repository := objects.NewRepository(pool)
	handler := objects.NewHandler(repository)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
	})
	http.HandleFunc("/objects", handler.ListObjects)
	log.Println("сервер слушает порт :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
