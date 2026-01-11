package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Sol1tud9/taskflow/internal/domain"
	taskRepo "github.com/Sol1tud9/taskflow/internal/task/repository"
	taskUsecase "github.com/Sol1tud9/taskflow/internal/task/usecase"
)

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	AssigneeID  string `json:"assignee_id"`
	CreatorID   string `json:"creator_id"`
	TeamID      string `json:"team_id"`
	DueDate     int64  `json:"due_date"`
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := taskUsecase.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		AssigneeID:  req.AssigneeID,
		CreatorID:   req.CreatorID,
		TeamID:      req.TeamID,
		DueDate:     req.DueDate,
	}

	task, err := h.taskUC.CreateTask(r.Context(), input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = h.activityUC.RecordActivity(r.Context(), req.CreatorID, domain.EntityTypeTask, task.ID, domain.ActionTypeCreated, `{"title":"`+task.Title+`"}`)

	_ = h.cache.Set(r.Context(), "task:"+task.ID, task)

	respondJSON(w, http.StatusCreated, task)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cacheKey := "task:" + id
	var cachedTask domain.Task

	var task *domain.Task
	if err := h.cache.Get(r.Context(), cacheKey, &cachedTask); err == nil {
		task = &cachedTask
	} else {
		fetchedTask, err := h.taskUC.GetTask(r.Context(), id)
		if err != nil {
			respondError(w, http.StatusNotFound, "task not found")
			return
		}
		task = fetchedTask
		_ = h.cache.Set(r.Context(), cacheKey, task)
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	if limit <= 0 {
		limit = 20
	}

	filter := taskRepo.TaskFilter{
		TeamID:     query.Get("team_id"),
		AssigneeID: query.Get("assignee_id"),
		Status:     query.Get("status"),
		Limit:      limit,
		Offset:     offset,
	}

	tasks, total, err := h.taskUC.ListTasks(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if tasks == nil {
		tasks = []*domain.Task{}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"tasks": tasks,
		"total": total,
	})
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	AssigneeID  string `json:"assignee_id"`
	DueDate     int64  `json:"due_date"`
	UserID      string `json:"user_id"`
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateTaskRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := taskUsecase.UpdateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		AssigneeID:  req.AssigneeID,
		DueDate:     req.DueDate,
		UserID:      req.UserID,
	}

	task, err := h.taskUC.UpdateTask(r.Context(), id, input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = h.cache.Set(r.Context(), "task:"+id, task)

	respondJSON(w, http.StatusOK, task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.taskUC.DeleteTask(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = h.cache.Delete(r.Context(), "task:"+id)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"id":      id,
	})
}

func (h *Handler) GetTaskHistory(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "task_id")

	history, err := h.taskUC.GetTaskHistory(r.Context(), taskID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if history == nil {
		history = []*domain.TaskHistory{}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"history": history,
		"total":   len(history),
	})
}
