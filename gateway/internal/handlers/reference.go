package handlers

import (
	"fmt"
	"net/http"

	"github.com/HJyup/translatify-common/utils"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gorilla/mux"
)

type ReferenceHandler struct{}

func NewReferenceHandler() *ReferenceHandler {
	return &ReferenceHandler{}
}

func (h *ReferenceHandler) HandleAPIReference(w http.ResponseWriter, r *http.Request) {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: "../docs/swagger.json",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "translatify",
		},
		DarkMode: true,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("Error generating API reference: %v", err))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprintln(w, htmlContent)
}

func (h *ReferenceHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/reference", h.HandleAPIReference).Methods("GET")
}
