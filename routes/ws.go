package routes

import (
	"Forum/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func WSRoutes(r *mux.Router) {
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWs(&handlers.HubInstance, w, r)
	})
}
