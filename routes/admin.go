package routes

import (
	"Forum/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func AdminRoutes(r *mux.Router) {
	r.Handle("/profile/ADMIN", handlers.SessionMiddleware((http.HandlerFunc(handlers.AdminProfileHandler)))).Methods("GET")
	r.HandleFunc("/promote_user", handlers.PromoteUserToModeratorHandler)
	r.HandleFunc("/demote_moderator", handlers.DemoteModeratorToUserHandler)
	r.HandleFunc("/create_admin", handlers.CreateAdminHandler)
	r.HandleFunc("/respond_report", handlers.RespondToReportHandler)
}
