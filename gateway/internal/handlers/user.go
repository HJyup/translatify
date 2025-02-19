package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-common/utils"
	"github.com/HJyup/translatify-gateway/internal/gateway/user"
	"github.com/HJyup/translatify-gateway/internal/models"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type UserHandler struct {
	gateway user.Gateway
}

func NewUserHandler(gateway user.Gateway) *UserHandler {
	return &UserHandler{gateway: gateway}
}

func (h *UserHandler) RegisterRoutes(router *mux.Router) {
	userRouter := router.PathPrefix("/api/v1/users").Subrouter()
	userRouter.HandleFunc("", h.HandleCreateUser).Methods("POST")
	userRouter.Handle("", utils.TokenAuthMiddleware(http.HandlerFunc(h.HandleListUsers))).Methods("GET")
	userRouter.Handle("/{username}", utils.TokenAuthMiddleware(http.HandlerFunc(h.HandleGetUser))).Methods("GET")
	userRouter.Handle("/{userId}", utils.TokenAuthMiddleware(http.HandlerFunc(h.HandleDeleteUser))).Methods("DELETE")
}

// HandleCreateUser godoc
// @Summary Create User
// @Description Create a new user account.
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "User creation data"
// @Success 201 {object} api.CreateUserResponse "User created successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/users [post]
func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var reqBody models.CreateUserRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &reqBody); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), "HandleCreateUser")
	defer span.End()

	resp, err := h.gateway.CreateUser(ctx, &api.CreateUserRequest{
		Username: reqBody.UserName,
		Email:    reqBody.Email,
		Password: reqBody.Password,
		Language: reqBody.Language,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusCreated, resp)
}

// HandleListUsers godoc
// @Summary List Users
// @Description Retrieve a paginated list of users.
// @Tags users
// @Produce json
// @Param limit query int false "Maximum number of users to return"
// @Param pageToken query string false "Pagination token"
// @Success 200 {object} api.ListUsersResponse "List of users"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users [get]
func (h *UserHandler) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := otel.Tracer("http")
	ctx, span := tr.Start(ctx, "HandleListUsers")
	defer span.End()

	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		} else {
			utils.WriteError(w, http.StatusBadRequest, "Invalid limit")
			return
		}
	}
	pageToken := r.URL.Query().Get("pageToken")

	resp, err := h.gateway.ListUsers(ctx, &api.ListUsersRequest{
		Limit:     int32(limit),
		PageToken: pageToken,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Error != "" {
		utils.WriteError(w, http.StatusBadRequest, resp.Error)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

// HandleGetUser godoc
// @Summary Get User
// @Description Retrieve a user's details by username.
// @Tags users
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} api.GetUserResponse "User details"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users/{username} [get]
func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	if username == "" {
		utils.WriteError(w, http.StatusBadRequest, "Username is required")
		return
	}

	ctx := r.Context()
	tr := otel.Tracer("http")
	ctx, span := tr.Start(ctx, "HandleGetUser")
	defer span.End()

	resp, err := h.gateway.GetUser(ctx, &api.GetUserRequest{
		Username: username,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.Error != "" {
		utils.WriteError(w, http.StatusBadRequest, resp.Error)
		return
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}

// HandleDeleteUser godoc
// @Summary Delete User
// @Description Delete a user by userID. Only the authenticated user may delete their account.
// @Tags users
// @Produce json
// @Param userId path string true "User ID to delete"
// @Success 200 {object} map[string]bool "Deletion confirmation"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users/{userId} [delete]
func (h *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	if userId == "" {
		utils.WriteError(w, http.StatusBadRequest, "UserID is required")
		return
	}

	tokenUserID, ok := r.Context().Value("userID").(string)
	if !ok || tokenUserID != userId {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx := r.Context()
	tr := otel.Tracer("http")
	ctx, span := tr.Start(ctx, "HandleDeleteUser")
	defer span.End()

	resp, err := h.gateway.DeleteUser(ctx, &api.DeleteUserRequest{
		UserId: userId,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !resp.Success {
		utils.WriteError(w, http.StatusBadRequest, resp.Error)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]bool{"success": resp.Success})
}
