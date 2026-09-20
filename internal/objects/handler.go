package objects

import (
	"encoding/json"
	"net/http"
)

var errGetObjects string = "не удалось получить объекты"

type Handler struct {
	repo *Repository
}

func (h *Handler) ListObjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	list, err := h.repo.ListObjects(ctx)
	if err != nil {
		http.Error(w, errGetObjects, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}
