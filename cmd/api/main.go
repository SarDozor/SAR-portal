package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/darkn3sssss/sar-dozor/internal/db"
	insarjob "github.com/darkn3sssss/sar-dozor/internal/insarJob"
	"github.com/darkn3sssss/sar-dozor/internal/measurement"
	objects "github.com/darkn3sssss/sar-dozor/internal/object"

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

	repositoryObjects := objects.NewRepositoryObject(pool)
	handlerObjects := objects.NewHandlerObject(repositoryObjects)

	repositoryInsarJob := insarjob.NewRepositoryInsarJob(pool)
	handlerInsarJob := insarjob.NewHandlerInsarJob(repositoryInsarJob)

	repositoryMeasurement := measurement.NewRepositoryMeasurement(pool)
	handlerMeasurement := measurement.NewHandlerMeasurement(repositoryMeasurement)

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

	mux.HandleFunc("GET /insar", handlerInsarJob.ListInsarJob)
	mux.HandleFunc("GET /insar/{id}", handlerInsarJob.GetInsarJobByID)
	mux.HandleFunc("PUT /insar/{id}", handlerInsarJob.UpdateInsarJob)
	mux.HandleFunc("POST /insar", handlerInsarJob.CreateInsarJob)
	mux.HandleFunc("DELETE /insar/{id}", handlerInsarJob.DeleteInsarJob)

	mux.HandleFunc("GET /measurements", handlerMeasurement.ListMeasurement)
	mux.HandleFunc("GET /measurements/{id}", handlerMeasurement.GetMeasurementByID)
	mux.HandleFunc("PUT /measurements/{id}", handlerMeasurement.UpdateMeasurement)
	mux.HandleFunc("POST /measurements", handlerMeasurement.CreateMeasurement)
	mux.HandleFunc("DELETE /measurements/{id}", handlerMeasurement.DeleteMeasurement)

	log.Println("сервер слушает порт :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
