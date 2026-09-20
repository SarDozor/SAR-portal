package objects

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo *Repository
}

func (h *Handler) ListObjects(w http.ResponseWriter, r *http.Request) {
	list, err := h.repo.ListObjects(r.Context())
	if err != nil {
		http.Error(w, "не удалось получить объекты", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) CreateObject(w http.ResponseWriter, r *http.Request) {
	var obj Object
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

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}
