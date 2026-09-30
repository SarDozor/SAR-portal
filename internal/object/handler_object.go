package objects

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	objects "github.com/darkn3sssss/sar-dozor/internal/model"

	"github.com/jackc/pgx/v5"
)

type HandlerObject struct {
	repo *RepositoryObject
}

func (h *HandlerObject) ListObjects(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.ListObjects(r.Context())
	if err != nil {
		http.Error(w, "не удалось получить объекты", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *HandlerObject) GetObjectByID(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	obj, err := h.repo.GetObjectByID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "объект не найден", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "не удалось получить объект по id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}

func (h *HandlerObject) CreateObject(w http.ResponseWriter, r *http.Request) {
	var obj objects.Object
	if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	createdObj, err := h.repo.CreateObject(r.Context(), obj)
	if err != nil {
		http.Error(w, "не удалось создать объект", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdObj)
}

func (h *HandlerObject) UpdateObject(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	var obj objects.Object
	if err = json.NewDecoder(r.Body).Decode(&obj); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	updatedObj, err := h.repo.UpdateObject(r.Context(), id, obj)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "объект не найден", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "не удалось обновить объект", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedObj)
}

func (h *HandlerObject) DeleteObject(w http.ResponseWriter, r *http.Request) {
	strID := r.PathValue("id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "неверный id", http.StatusBadRequest)
		return
	}
	deletedObject, err := h.repo.DeleteObject(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		http.Error(w, "объект не найден", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "не удалось удалить объект", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedObject)
}

func NewHandlerObject(repo *RepositoryObject) *HandlerObject {
	return &HandlerObject{repo: repo}
}
