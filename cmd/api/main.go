package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/darkn3sssss/sar-dozor/internal/db"
	"github.com/darkn3sssss/sar-dozor/internal/insarjob"
	"github.com/darkn3sssss/sar-dozor/internal/measurement"
	"github.com/darkn3sssss/sar-dozor/internal/object"

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
	if dsn == "" {
		log.Fatal("не задана переменная DATABASE_URL")
	}
	ctx := context.Background()
	pool, err := db.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("подключились к базе")
	defer pool.Close()

	handlerObjects := object.NewHandler(object.NewRepository(pool))
	handlerInsarJobs := insarjob.NewHandler(insarjob.NewRepository(pool))
	handlerMeasurements := measurement.NewHandler(measurement.NewRepository(pool))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HealthResponse{Status: "ok"})
	})

	mux.HandleFunc("GET /objects", handlerObjects.ListObjects)
	mux.HandleFunc("GET /objects/{id}", handlerObjects.GetObjectByID)
	mux.HandleFunc("PUT /objects/{id}", handlerObjects.UpdateObject)
	mux.HandleFunc("POST /objects", handlerObjects.CreateObject)
	mux.HandleFunc("DELETE /objects/{id}", handlerObjects.DeleteObject)
	mux.HandleFunc("GET /objects/{id}/jobs", handlerInsarJobs.ListInsarJobsByObject)

	mux.HandleFunc("GET /jobs", handlerInsarJobs.ListInsarJobs)
	mux.HandleFunc("GET /jobs/{id}", handlerInsarJobs.GetInsarJobByID)
	mux.HandleFunc("PUT /jobs/{id}", handlerInsarJobs.UpdateInsarJob)
	mux.HandleFunc("POST /jobs", handlerInsarJobs.CreateInsarJob)
	mux.HandleFunc("DELETE /jobs/{id}", handlerInsarJobs.DeleteInsarJob)
	mux.HandleFunc("GET /jobs/{id}/measurements", handlerMeasurements.ListMeasurementsByJob)

	mux.HandleFunc("GET /measurements", handlerMeasurements.ListMeasurements)
	mux.HandleFunc("GET /measurements/{id}", handlerMeasurements.GetMeasurementByID)
	mux.HandleFunc("PUT /measurements/{id}", handlerMeasurements.UpdateMeasurement)
	mux.HandleFunc("POST /measurements", handlerMeasurements.CreateMeasurement)
	mux.HandleFunc("DELETE /measurements/{id}", handlerMeasurements.DeleteMeasurement)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Println("сервер слушает порт :8080")
	log.Fatal(server.ListenAndServe())
}
