package routes

import (
	"Forum/handlers"
	"github.com/gorilla/mux"
)

func Appliroute(r *mux.Router) {
	r.HandleFunc("/", handlers.IndexHandler)
	r.HandleFunc("/signup", handlers.SignupHandler)
	r.HandleFunc("/login", handlers.LoginHandler)
	r.HandleFunc("/image", handlers.ImageHandler)
	r.HandleFunc("/about", handlers.AboutHandler)
}
