package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/eugenius-watchman/ecom_go_rest_api/cmd/service/user"
	"github.com/gorilla/mux"
)

// blueprint
type APIServer struct {
	addr string
	db *sql.DB
}

// constructor
func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db: db,
	}
}
 
// run method
func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// handler for users
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)


	log.Println("Listening on", s.addr)
	
	return http.ListenAndServe(s.addr, router)
}
