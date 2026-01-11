package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/Sol1tud9/taskflow/internal/domain"
)

type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userUC.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = h.activityUC.RecordActivity(r.Context(), user.ID, domain.EntityTypeUser, user.ID, domain.ActionTypeCreated, `{"email":"`+user.Email+`","name":"`+user.Name+`"}`)

	_ = h.cache.Set(r.Context(), "user:"+user.ID, user)
	_ = h.cache.Delete(r.Context(), "users:list")

	respondJSON(w, http.StatusCreated, user)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	var users []*domain.User
	var cached struct {
		Users []*domain.User `json:"users"`
		Total int            `json:"total"`
	}

	if err := h.cache.Get(r.Context(), "users:list", &cached); err == nil {
		users = cached.Users
	} else {
		listUsers, err := h.userLister.List(r.Context())
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if listUsers == nil {
			listUsers = []*domain.User{}
		}

		users = listUsers
		_ = h.cache.Set(r.Context(), "users:list", map[string]interface{}{
			"users": users,
			"total": len(users),
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"users": users,
		"total": len(users),
	})
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cacheKey := "user:" + id
	var cachedUser domain.User

	var user *domain.User
	if err := h.cache.Get(r.Context(), cacheKey, &cachedUser); err == nil {
		user = &cachedUser
	} else {
		fetchedUser, err := h.userUC.GetUser(r.Context(), id)
		if err != nil {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		user = fetchedUser
		_ = h.cache.Set(r.Context(), cacheKey, user)
	}

	respondJSON(w, http.StatusOK, user)
}

type UpdateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userUC.UpdateUser(r.Context(), id, req.Email, req.Name)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = h.cache.Set(r.Context(), "user:"+id, user)
	_ = h.cache.Delete(r.Context(), "users:list")

	respondJSON(w, http.StatusOK, user)
}

type CreateTeamRequest struct {
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req CreateTeamRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	team, err := h.teamUC.CreateTeam(r.Context(), req.Name, req.OwnerID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = h.activityUC.RecordActivity(r.Context(), req.OwnerID, domain.EntityTypeTeam, team.ID, domain.ActionTypeCreated, `{"name":"`+team.Name+`"}`)

	_ = h.cache.Set(r.Context(), "team:"+team.ID, team)
	_ = h.cache.Delete(r.Context(), "teams:list")

	respondJSON(w, http.StatusCreated, team)
}

func (h *Handler) ListTeams(w http.ResponseWriter, r *http.Request) {
	var teams []*domain.Team
	var cached struct {
		Teams []*domain.Team `json:"teams"`
		Total int            `json:"total"`
	}

	if err := h.cache.Get(r.Context(), "teams:list", &cached); err == nil {
		teams = cached.Teams
	} else {
		listTeams, err := h.teamLister.ListTeams(r.Context())
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if listTeams == nil {
			listTeams = []*domain.Team{}
		}

		teams = listTeams
		_ = h.cache.Set(r.Context(), "teams:list", map[string]interface{}{
			"teams": teams,
			"total": len(teams),
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"teams": teams,
		"total": len(teams),
	})
}

func (h *Handler) GetTeam(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	cacheKey := "team:" + id
	var cachedTeam domain.Team

	var team *domain.Team
	if err := h.cache.Get(r.Context(), cacheKey, &cachedTeam); err == nil {
		team = &cachedTeam
	} else {
		fetchedTeam, err := h.teamUC.GetTeam(r.Context(), id)
		if err != nil {
			respondError(w, http.StatusNotFound, "team not found")
			return
		}
		team = fetchedTeam
		_ = h.cache.Set(r.Context(), cacheKey, team)
	}

	respondJSON(w, http.StatusOK, team)
}

type AddTeamMemberRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func (h *Handler) AddTeamMember(w http.ResponseWriter, r *http.Request) {
	teamID := chi.URLParam(r, "team_id")

	var req AddTeamMemberRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	member, err := h.teamUC.AddTeamMember(r.Context(), teamID, req.UserID, req.Role)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_ = h.cache.Delete(r.Context(), "team:"+teamID+":members")
	_ = h.cache.Delete(r.Context(), "team:"+teamID)

	respondJSON(w, http.StatusCreated, member)
}

func (h *Handler) GetTeamMembers(w http.ResponseWriter, r *http.Request) {
	teamID := chi.URLParam(r, "team_id")

	var members []*domain.TeamMember
	cacheKey := "team:" + teamID + ":members"

	var cached struct {
		Members []*domain.TeamMember `json:"members"`
		Total   int                  `json:"total"`
	}

	if err := h.cache.Get(r.Context(), cacheKey, &cached); err == nil {
		members = cached.Members
	} else {
		fetchedMembers, err := h.teamUC.GetTeamMembers(r.Context(), teamID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if fetchedMembers == nil {
			fetchedMembers = []*domain.TeamMember{}
		}

		members = fetchedMembers
		_ = h.cache.Set(r.Context(), cacheKey, map[string]interface{}{
			"members": members,
			"total":   len(members),
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"members": members,
		"total":   len(members),
	})
}

func (h *Handler) GetUserActivities(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	from, _ := strconv.ParseInt(r.URL.Query().Get("from"), 10, 64)
	to, _ := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)

	if limit <= 0 {
		limit = 20
	}

	activities, total, err := h.activityUC.GetUserActivities(r.Context(), userID, from, to, limit, offset)
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
