package user

import (
	//"encoding/json"
	"net/http"

	"github.com/eugenius-watchman/ecom_go_rest_api/types"
	"github.com/eugenius-watchman/ecom_go_rest_api/utils"
	"github.com/gorilla/mux"
)

type Handler struct {

}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

// metheod for the handler 
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

} 

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON payloads
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	// check if user exists ... db


	// if no ... we create new user

} 