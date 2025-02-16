package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	pb "github.com/HJyup/translatify-common/api"
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

	resp, err := h.gateway.CreateUser(ctx, &pb.CreateUserRequest{
		Username: reqBody.UserName,
		Email:    reqBody.Email,
		Password: reqBody.Password,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusCreated, resp)
}

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

	resp, err := h.gateway.ListUsers(ctx, &pb.ListUsersRequest{
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

	resp, err := h.gateway.GetUser(ctx, &pb.GetUserRequest{
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

	resp, err := h.gateway.DeleteUser(ctx, &pb.DeleteUserRequest{
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
