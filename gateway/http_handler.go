package main

import (
	common "github.com/HJyup/translatify-common"
	"net/http"

	pb "github.com/HJyup/translatify-common/api"
)

type Handler struct {
	client pb.UserServiceClient
}

func NewHandler(client pb.UserServiceClient) *Handler {
	return &Handler{client: client}
}

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/users/register", h.handleRegister)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	_, err := h.client.CreateUser(r.Context(), &pb.CreateUserRequest{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
