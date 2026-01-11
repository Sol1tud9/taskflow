package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/Sol1tud9/taskflow/internal/domain"
	"github.com/Sol1tud9/taskflow/internal/gateway/cache"
	taskRepo "github.com/Sol1tud9/taskflow/internal/task/repository"
	taskUsecase "github.com/Sol1tud9/taskflow/internal/task/usecase"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, email, name string) (*domain.User, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, id, email, name string) (*domain.User, error)
}

type TeamUseCase interface {
	CreateTeam(ctx context.Context, name, ownerID string) (*domain.Team, error)
	GetTeam(ctx context.Context, id string) (*domain.Team, error)
	AddTeamMember(ctx context.Context, teamID, userID, role string) (*domain.TeamMember, error)
	GetTeamMembers(ctx context.Context, teamID string) ([]*domain.TeamMember, error)
}

type TaskUseCase interface {
	CreateTask(ctx context.Context, input taskUsecase.CreateTaskInput) (*domain.Task, error)
	GetTask(ctx context.Context, id string) (*domain.Task, error)
	ListTasks(ctx context.Context, filter taskRepo.TaskFilter) ([]*domain.Task, int, error)
	UpdateTask(ctx context.Context, id string, input taskUsecase.UpdateTaskInput) (*domain.Task, error)
	DeleteTask(ctx context.Context, id string) error
	GetTaskHistory(ctx context.Context, taskID string) ([]*domain.TaskHistory, error)
}

type ActivityUseCase interface {
	GetUserActivities(ctx context.Context, userID string, from, to int64, limit, offset int) ([]*domain.Activity, int, error)
	GetActivities(ctx context.Context, entityType, entityID string, from, to int64, limit, offset int) ([]*domain.Activity, int, error)
	RecordActivity(ctx context.Context, userID string, entityType domain.EntityType, entityID string, action domain.ActionType, metadata string) error
}

type UserLister interface {
	List(ctx context.Context) ([]*domain.User, error)
}

type TeamLister interface {
	ListTeams(ctx context.Context) ([]*domain.Team, error)
}

type Handler struct {
	cache      *cache.RedisCache
	userUC     UserUseCase
	teamUC     TeamUseCase
	taskUC     TaskUseCase
	activityUC ActivityUseCase
	userLister UserLister
	teamLister TeamLister
}

func NewHandler(
	cache *cache.RedisCache,
	userUC UserUseCase,
	teamUC TeamUseCase,
	taskUC TaskUseCase,
	activityUC ActivityUseCase,
	userLister UserLister,
	teamLister TeamLister,
) *Handler {
	return &Handler{
		cache:      cache,
		userUC:     userUC,
		teamUC:     teamUC,
		taskUC:     taskUC,
		activityUC: activityUC,
		userLister: userLister,
		teamLister: teamLister,
	}
}

func (h *Handler) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", h.CreateUser)
			r.Get("/", h.ListUsers)
			r.Get("/{id}", h.GetUser)
			r.Patch("/{id}", h.UpdateUser)
			r.Get("/{user_id}/activities", h.GetUserActivities)
		})

		r.Route("/teams", func(r chi.Router) {
			r.Post("/", h.CreateTeam)
			r.Get("/", h.ListTeams)
			r.Get("/{id}", h.GetTeam)
			r.Post("/{team_id}/members", h.AddTeamMember)
			r.Get("/{team_id}/members", h.GetTeamMembers)
		})

		r.Route("/tasks", func(r chi.Router) {
			r.Post("/", h.CreateTask)
			r.Get("/", h.ListTasks)
			r.Get("/{id}", h.GetTask)
			r.Patch("/{id}", h.UpdateTask)
			r.Delete("/{id}", h.DeleteTask)
			r.Get("/{task_id}/history", h.GetTaskHistory)
		})

		r.Route("/activities", func(r chi.Router) {
			r.Get("/", h.GetActivities)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			// Log error if logger is available, but client already received headers
			return
		}
	})

	r.Get("/swagger.json", h.SwaggerJSON)
	r.Get("/swagger", h.SwaggerUI)

	return r
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log error if logger is available, but client already received headers
		return
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
