package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Laefye/go-search/internal/rabbitmq/events"
	"github.com/Laefye/go-search/internal/service/search"
	"github.com/Laefye/go-search/internal/service/top"
)

type Handler struct {
	topService    *top.TopService
	searchService *search.SearchService
}

func NewHandler(topService *top.TopService, searchService *search.SearchService) *Handler {
	return &Handler{
		topService:    topService,
		searchService: searchService,
	}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("/top", h.GetTopQueries)
	router.HandleFunc("/search", h.Search)
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

	w.Header().Set("Content-Type", "application/json, charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	userID := r.URL.Query().Get("user_id")

	queryEvent := events.SearchQueryEvent{
		Query:     query,
		UserID:    userID,
		Timestamp: time.Now(),
	}

	err := h.searchService.Publish(r.Context(), queryEvent)
	if err != nil {
		if errors.Is(err, search.ErrInvalidQuery) {
			http.Error(w, "Invalid query parameter", http.StatusBadRequest)
			return
		}

		http.Error(w, "Failed to publish search event", http.StatusInternalServerError)
		log.Printf("Error publishing search event: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
