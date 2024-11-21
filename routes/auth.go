package routes

import (
	"Forum/handlers"
	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router) {
	r.HandleFunc("/auth/google/login", handlers.GoogleLoginHandler)
	r.HandleFunc("/auth/google/callback", handlers.GoogleCallbackHandler)
	r.HandleFunc("/auth/github/login", handlers.GithubLoginHandler)
	r.HandleFunc("/auth/github/callback", handlers.GithubCallbackHandler)
	r.HandleFunc("/auth/facebook/login", handlers.FacebookLoginHandler)
	r.HandleFunc("/auth/facebook/callback", handlers.FacebookCallbackHandler)
}
