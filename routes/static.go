package routes

import (
	"net/http"
	"github.com/gorilla/mux"
)

func StaticRoutes(r *mux.Router) {
	r.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))
	r.PathPrefix("/script/").Handler(http.StripPrefix("/script/", http.FileServer(http.Dir("./script"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.PathPrefix("/Ressources/").Handler(http.StripPrefix("/Ressources/", http.FileServer(http.Dir("./Ressources"))))
}
