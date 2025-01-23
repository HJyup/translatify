package main

import "net/http"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/register", h.handleRegister)
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// handle register
}
