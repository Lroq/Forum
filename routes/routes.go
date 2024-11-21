package routes

import (
	"github.com/gorilla/mux"
)

func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	// appli ropute
	Appliroute(r)
	
	// Static routes
	StaticRoutes(r)
	
	// Authentication routes
	AuthRoutes(r)
	
	// Post routes
	PostRoutes(r)
	
	// Profile routes
	ProfileRoutes(r)
	
	// Moderation routes
	ModerationRoutes(r)

	//admin routes
	AdminRoutes(r)

	//report routes
	ReportRoutes(r)

	// delete routes
	DeleteRoutes(r)

	// Websocket routes
	WSRoutes(r)

	return r
}
