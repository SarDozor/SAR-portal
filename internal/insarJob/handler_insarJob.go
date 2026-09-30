package insarjob

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	objects "github.com/darkn3sssss/sar-dozor/internal/model"
	"github.com/jackc/pgx/v5"
)

type HandlerInsarJob struct {
	repo *RepositoryInsarJob
}

func (h *HandlerInsarJob) ListInsarJob(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.ListInsarJob(r.Context())
	if err != nil {
		http.Error(w, "не удалось получить объекты", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *HandlerInsarJob) GetInsarJobByID(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	insarJob, err := h.repo.GetInsarJobByID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "объект не найден", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "не удалось получить объект по id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insarJob)
}

func (h *HandlerInsarJob) CreateInsarJob(w http.ResponseWriter, r *http.Request) {
	var insarJob objects.InsarJob
	if err := json.NewDecoder(r.Body).Decode(&insarJob); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	createdInsarJob, err := h.repo.CreateInsarJob(r.Context(), insarJob)
	if err != nil {
		http.Error(w, "не удалось создать объект", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdInsarJob)
}

func (h *HandlerInsarJob) UpdateInsarJob(w http.ResponseWriter, r *http.Request) {
	var insarJob objects.InsarJob
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	if err = json.NewDecoder(r.Body).Decode(&insarJob); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	updatedInsarJob, err := h.repo.UpdateInsarJob(r.Context(), id, insarJob)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "объект не найден", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "не удалось получить объект по id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedInsarJob)
}

func (h *HandlerInsarJob) DeleteInsarJob(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	deletedInsarJob, err := h.repo.DeleteInsarJob(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "объект не найден", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "не удалось получить объект по id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedInsarJob)
}

func NewHandlerInsarJob(repo *RepositoryInsarJob) *HandlerInsarJob {
	return &HandlerInsarJob{repo: repo}
}
