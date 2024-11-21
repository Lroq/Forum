package routes

import (
	"Forum/handlers"
	"github.com/gorilla/mux"
)

func PostRoutes(r *mux.Router) {
	r.HandleFunc("/new_post", handlers.NewPostHandler)
	r.HandleFunc("/post", handlers.GetPostHandler)
	r.HandleFunc("/like", handlers.LikeHandler)
	r.HandleFunc("/dislike", handlers.DislikeHandler)
	r.HandleFunc("/top_posts", handlers.TopPostHandler)
	r.HandleFunc("/posts/{id:[0-9]+}", handlers.PostDetailHandler)
	r.HandleFunc("/add-comment", handlers.AddCommentHandler)
	r.HandleFunc("/upload_image", handlers.ImageUploadHandler)
	r.HandleFunc("/check-authorization", handlers.CheckAuthorizationHandler)
}
