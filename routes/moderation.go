package routes

import (
	"Forum/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func ModerationRoutes(r *mux.Router) {
	//Moderator routes
	r.Handle("/profile/moderator", handlers.SessionMiddleware((http.HandlerFunc(handlers.ModeratorProfileHandler)))).Methods("GET")
	r.HandleFunc("/send_report", handlers.SendReportHandler)
}
