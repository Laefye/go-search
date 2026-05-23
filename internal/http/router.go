package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Laefye/go-search/internal/service/top"
)

type Handler struct {
	topService *top.TopService
}

func NewHandler(topService *top.TopService) *Handler {
	return &Handler{topService: topService}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/top", h.GetTopQueries)
}

func (h *Handler) GetTopQueries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil || limit <= 0 {
		http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
		return
	}

	response, err := h.topService.GetTopQueries(ctx, limit)
	if err != nil {
		http.Error(w, "Failed to get top queries", http.StatusInternalServerError)
		log.Printf("Error getting top queries: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
