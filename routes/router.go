package routes

import (
	"fmt"
	"github.com/MuhaFAH/AuthCheck/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	Mux *mux.Router
}

func NewRouter(app *handlers.App) *Router {
	router := mux.NewRouter()

	router.HandleFunc("/auth/{guid}", app.IssuingTokensHandler).Methods("GET")
	router.HandleFunc("/refresh/{guid}", app.RefreshTokensHandler).Methods("POST")

	return &Router{router}
}

func (r *Router) Run(port string) error {
	return http.ListenAndServe(fmt.Sprintf(":%s", port), r.Mux)
}
