package insarjob

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

func (h *Handler) ListInsarJobs(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.ListInsarJobs(r.Context())
	if err != nil {
		log.Printf("ListInsarJobs: %v", err)
		http.Error(w, "не удалось получить задачи", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) ListInsarJobsByObject(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	objectID, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	list, err := h.repo.ListInsarJobsByObject(r.Context(), objectID)
	if err != nil {
		log.Printf("ListInsarJobsByObject %d: %v", objectID, err)
		http.Error(w, "не удалось получить задачи объекта", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) GetInsarJobByID(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	insarJob, err := h.repo.GetInsarJobByID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "задача не найдена", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("GetInsarJobByID %d: %v", id, err)
		http.Error(w, "не удалось получить задачу по id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insarJob)
}

func (h *Handler) CreateInsarJob(w http.ResponseWriter, r *http.Request) {
	var insarJob model.InsarJob
	if err := json.NewDecoder(r.Body).Decode(&insarJob); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	createdInsarJob, err := h.repo.CreateInsarJob(r.Context(), insarJob)
	if err != nil {
		log.Printf("CreateInsarJob: %v", err)
		http.Error(w, "не удалось создать задачу", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdInsarJob)
}

func (h *Handler) UpdateInsarJob(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	var insarJob model.InsarJob
	if err = json.NewDecoder(r.Body).Decode(&insarJob); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	updatedInsarJob, err := h.repo.UpdateInsarJob(r.Context(), id, insarJob)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "задача не найдена", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("UpdateInsarJob %d: %v", id, err)
		http.Error(w, "не удалось обновить задачу", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedInsarJob)
}

func (h *Handler) DeleteInsarJob(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	deletedInsarJob, err := h.repo.DeleteInsarJob(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "задача не найдена", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("DeleteInsarJob %d: %v", id, err)
		http.Error(w, "не удалось удалить задачу", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedInsarJob)
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}
