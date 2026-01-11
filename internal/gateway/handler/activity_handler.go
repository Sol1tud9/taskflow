package handler

import (
	"net/http"
	"strconv"

	"github.com/Sol1tud9/taskflow/internal/domain"
)

func (h *Handler) GetActivities(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	entityType := query.Get("entity_type")
	entityID := query.Get("entity_id")
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	from, _ := strconv.ParseInt(query.Get("from"), 10, 64)
	to, _ := strconv.ParseInt(query.Get("to"), 10, 64)

	if limit <= 0 {
		limit = 20
	}

	activities, total, err := h.activityUC.GetActivities(r.Context(), entityType, entityID, from, to, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if activities == nil {
		activities = []*domain.Activity{}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"activities": activities,
		"total":      total,
	})
}
