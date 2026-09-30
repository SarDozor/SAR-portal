package measurement

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	objects "github.com/darkn3sssss/sar-dozor/internal/model"

	"github.com/jackc/pgx/v5"
)

type HandleMeasurement struct {
	repo *RepositoryMeasurement
}

func (h *HandleMeasurement) ListMeasurement(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.ListMeasurement(r.Context())
	if err != nil {
		http.Error(w, "не удалось получить объекты", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *HandleMeasurement) GetMeasurementByID(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	measurement, err := h.repo.GetMeasurementByID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "объект не найден", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "не удалось получить объект по id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurement)
}

func (h *HandleMeasurement) CreateMeasurement(w http.ResponseWriter, r *http.Request) {
	var createdMeasurement objects.Measurement
	if err := json.NewDecoder(r.Body).Decode(&createdMeasurement); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	createdMeasurement, err := h.repo.CreateMeasurement(r.Context(), createdMeasurement)
	if err != nil {
		http.Error(w, "не удалось создать объект", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdMeasurement)
}

func (h *HandleMeasurement) UpdateMeasurement(w http.ResponseWriter, r *http.Request) {
	var measurement objects.Measurement
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&measurement); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	updatedMeasurement, err := h.repo.UpdateMeasurement(r.Context(), id, measurement)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "объект не найден", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "не удалось получить объект по id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedMeasurement)

}

func (h *HandleMeasurement) DeleteMeasurement(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	deletedMeasurement, err := h.repo.DeleteMeasurement(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "объект не найден", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "не удалось получить объект по id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedMeasurement)
}

func NewHandlerMeasurement(repo *RepositoryMeasurement) *HandleMeasurement {
	return &HandleMeasurement{repo: repo}
}
