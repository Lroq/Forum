package routes

import (
	"Forum/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func ProfileRoutes(r *mux.Router) {
	r.Handle("/profile/utilisateur", handlers.SessionMiddleware((http.HandlerFunc(handlers.UtilisateurProfileHandler)))).Methods("GET")
	r.HandleFunc("/update_profile", handlers.UpdateProfileHandler)
	r.HandleFunc("/user_posts", handlers.UserPostsHandler)
	r.HandleFunc("/user_comments", handlers.UserCommentsHandler)
	r.HandleFunc("/user_likes", handlers.UserLikesHandler)
	r.HandleFunc("/user_dislikes", handlers.UserDislikesHandler)
}
