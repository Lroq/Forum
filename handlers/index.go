package handlers

import (
	"Forum/database"
	"Forum/post_logic"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	store = sessions.NewCookieStore([]byte("secret-key"))
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 3600,
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	if session.IsNew {
		log.Println("Creating new guest user...")
		userUUID, err := database.CreateGuestUser()
		if err != nil {
			log.Printf("Error creating guest user: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Printf("Guest user created with ID: %v\n", userUUID)
		role, err := database.GetUserRoleByUUID(userUUID)
		if err != nil {
			log.Printf("Error getting user role: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		username := fmt.Sprintf("inviter%v", userUUID)
		session.Values["user_uuid"] = userUUID
		session.Values["username"] = username
		session.Values["role"] = role
		session.Save(r, w)
	} else {
		log.Println("Existing session detected, skipping guest user creation.")
	}
	var data struct {
		User  *database.User
		Posts []post_logic.Post
	}
	user, err := database.FetchUserByUUID(session.Values["user_uuid"].(string))
	if err != nil {
		log.Printf("Error fetching user: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data.User = user
	posts, err := database.GetPosts()
	if err != nil {
		log.Printf("Error getting posts: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	var postList []post_logic.Post
	for _, post := range posts {
		postList = append(postList, post_logic.Post(post))
	}
	data.Posts = post_logic.ReversePosts(postList)
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func TopPostHandler(w http.ResponseWriter, r *http.Request) {
	titles, err := database.FetchTopPostsByLikes()
	if err != nil {
		http.Error(w, "Failed to fetch top posts", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(titles)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/about.html")
}
