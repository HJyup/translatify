package handlers

import (
	"encoding/json"
	pb "github.com/HJyup/translatify-common/api"
	"github.com/HJyup/translatify-gateway/internal/models"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/HJyup/translatify-common/token"
	"github.com/HJyup/translatify-common/utils"
	"github.com/HJyup/translatify-gateway/internal/gateway/user"
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

	userRouter.Handle("", token.TokenAuthMiddleware(http.HandlerFunc(h.HandleListUsers))).Methods("GET")
	userRouter.Handle("/{userId}", token.TokenAuthMiddleware(http.HandlerFunc(h.HandleGetUser))).Methods("GET")
	userRouter.Handle("/{userId}", token.TokenAuthMiddleware(http.HandlerFunc(h.HandleDeleteUser))).Methods("DELETE")
}

func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var reqBody models.CreateUserRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &reqBody); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), "HandleCreateUser")
	defer span.End()

	createdUser, err := h.gateway.CreateUser(ctx, &pb.CreateUserRequest{
		Username: reqBody.UserName,
		Email:    reqBody.Email,
		Password: reqBody.Password,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tokenString, err := token.CreateToken(createdUser.User.UserId, createdUser.User.Username, createdUser.User.Email)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, "failed to generate token: "+err.Error())
		return
	}

	resp := map[string]interface{}{
		"user":  createdUser,
		"token": tokenString,
	}
	utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *UserHandler) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tr := otel.Tracer("http")
	ctx, span := tr.Start(ctx, "HandleListUsers")
	defer span.End()

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid limit")
	}

	pageToken := r.URL.Query().Get("pageToken")

	users, err := h.gateway.ListUsers(ctx, &pb.ListUsersRequest{
		Limit:     int32(limit),
		PageToken: pageToken,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	if userId == "" {
		utils.WriteError(w, http.StatusBadRequest, "userId is required")
		return
	}

	ctx := r.Context()
	tr := otel.Tracer("http")
	ctx, span := tr.Start(ctx, "HandleGetUser")
	defer span.End()

	userDetails, err := h.gateway.GetUser(ctx, &pb.GetUserRequest{
		UserId: userId,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, userDetails)
}

func (h *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	if userId == "" {
		utils.WriteError(w, http.StatusBadRequest, "userId is required")
		return
	}

	tokenUserID, ok := r.Context().Value("userID").(string)
	if !ok || tokenUserID != userId {
		utils.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx := r.Context()
	tr := otel.Tracer("http")
	ctx, span := tr.Start(ctx, "HandleDeleteUser")
	defer span.End()

	res, err := h.gateway.DeleteUser(ctx, &pb.DeleteUserRequest{
		UserId: userId,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": strconv.FormatBool(res.Success)})
}
