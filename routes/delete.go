package routes

import (
	"Forum/handlers"
	"github.com/gorilla/mux"
)

func DeleteRoutes(r *mux.Router) {
	r.HandleFunc("/delete_post", handlers.DeletePostHandler)
	r.HandleFunc("/delete_comment", handlers.DeleteCommentHandler)
	r.HandleFunc("/delete_user", handlers.DeleteUserHandler)
}
