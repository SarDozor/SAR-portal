package measurement

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/darkn3sssss/sar-dozor/internal/model"

	"github.com/jackc/pgx/v5"
)

type Handler struct {
	repo *Repository
}

func (h *Handler) ListMeasurements(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.ListMeasurements(r.Context())
	if err != nil {
		log.Printf("ListMeasurements: %v", err)
		http.Error(w, "не удалось получить измерения", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) ListMeasurementsByJob(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	jobID, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	list, err := h.repo.ListMeasurementsByJob(r.Context(), jobID)
	if err != nil {
		log.Printf("ListMeasurementsByJob %d: %v", jobID, err)
		http.Error(w, "не удалось получить измерения задачи", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) GetMeasurementByID(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	measurement, err := h.repo.GetMeasurementByID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "измерение не найдено", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("GetMeasurementByID %d: %v", id, err)
		http.Error(w, "не удалось получить измерение по id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(measurement)
}

func (h *Handler) CreateMeasurement(w http.ResponseWriter, r *http.Request) {
	var measurement model.Measurement
	if err := json.NewDecoder(r.Body).Decode(&measurement); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	createdMeasurement, err := h.repo.CreateMeasurement(r.Context(), measurement)
	if err != nil {
		log.Printf("CreateMeasurement: %v", err)
		http.Error(w, "не удалось создать измерение", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdMeasurement)
}

func (h *Handler) UpdateMeasurement(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	var measurement model.Measurement
	if err = json.NewDecoder(r.Body).Decode(&measurement); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	updatedMeasurement, err := h.repo.UpdateMeasurement(r.Context(), id, measurement)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "измерение не найдено", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("UpdateMeasurement %d: %v", id, err)
		http.Error(w, "не удалось обновить измерение", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedMeasurement)
}

func (h *Handler) DeleteMeasurement(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	deletedMeasurement, err := h.repo.DeleteMeasurement(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "измерение не найдено", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("DeleteMeasurement %d: %v", id, err)
		http.Error(w, "не удалось удалить измерение", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedMeasurement)
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}
