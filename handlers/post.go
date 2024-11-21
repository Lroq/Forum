package handlers

import (
	"Forum/database"
	"encoding/json"
	"Forum/post_logic"
	"Forum/comment_logic"
	"html/template"
	"github.com/gorilla/mux"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.IsNew {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Error(w, "Failed to get username from session", http.StatusInternalServerError)
		return
	}

	userUUID, ok := session.Values["user_uuid"].(string)
	if !ok {
		http.Error(w, "Failed to get user UUID from session", http.StatusInternalServerError)
		return
	}

	role, err := database.GetUserRoleByUUID(userUUID)
	if err != nil {
		log.Printf("Failed to get user role: %v", err)
		http.Error(w, "Failed to validate user role", http.StatusInternalServerError)
		return
	}

	if role == "GUEST" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("postContent")

	var imagePath string
	var gifPath string

	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		os.MkdirAll("uploads", os.ModePerm)
		imagePath = header.Filename
		fullImagePath := filepath.Join("uploads", imagePath)
		dest, err := os.Create(fullImagePath)
		if err != nil {
			log.Printf("Failed to create image file: %v", err)
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
		defer dest.Close()
		_, err = io.Copy(dest, file)
		if err != nil {
			log.Printf("Failed to copy image to destination: %v", err)
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("No image uploaded: %v", err)
	}

	gifFile, gifHeader, err := r.FormFile("gif")
	if err == nil {
		defer gifFile.Close()
		os.MkdirAll("uploads", os.ModePerm)
		gifPath = gifHeader.Filename
		fullGifPath := filepath.Join("uploads", gifPath)
		gifDest, err := os.Create(fullGifPath)
		if err != nil {
			log.Printf("Failed to create gif file: %v", err)
			http.Error(w, "Failed to save gif", http.StatusInternalServerError)
			return
		}
		defer gifDest.Close()
		_, err = io.Copy(gifDest, gifFile)
		if err != nil {
			log.Printf("Failed to copy gif to destination: %v", err)
			http.Error(w, "Failed to save gif", http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("No gif uploaded: %v", err)
	}

	err = database.InsertPost(title, content, username, imagePath, gifPath)
	if err != nil {
		log.Printf("Failed to create post: %v", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}
	post, err := database.FetchPostByID(postIDInt)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Post: %+v", post)
}

func UserPostsHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving session: %v", err), http.StatusInternalServerError)
		return
	}
	if session.Values["user_uuid"] == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username, ok := session.Values["username"].(string)
	if !ok {
		http.Error(w, "Failed to get user ID from session", http.StatusInternalServerError)
		return
	}

	posts, err := database.FetchPostsByUsername(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user posts: %v", err), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func PostDetailHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving session: %v", err), http.StatusInternalServerError)
		return
	}
	if session.Values["user_uuid"] == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	userUUID, ok := session.Values["user_uuid"].(string)
	if !ok {
		http.Error(w, "Failed to get user ID from session", http.StatusInternalServerError)
		return
	}
	user, err := database.FetchUserByUUID(userUUID)
	if err != nil {
		log.Printf("Failed to get user role: %v", err)
		http.Error(w, "Failed to validate user role", http.StatusInternalServerError)
		return
	}
	postID := getPostIDFromRequest(r)
	post, err := fetchPostDetails(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Fetch comments for the post
	comments, err := fetchCommentsForPost(postID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Create a struct to hold both post and comments
	data := struct {
		Post     post_logic.Post
		Comments []comment_logic.Comment
		User     *database.User
	}{
		Post:     post,
		Comments: comments,
		User:     user,
	}
	tmpl := template.Must(template.ParseFiles("templates/post_detail.html"))
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// getPostIDFromRequest pourrait extraire l'ID du post de la requête.
func getPostIDFromRequest(r *http.Request) string {
	vars := mux.Vars(r)
	return vars["id"]
}

// fetchPostDetails pourrait récupérer les détails d'un post à partir de la base de données.
func fetchPostDetails(postID string) (post_logic.Post, error) {
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		return post_logic.Post{}, fmt.Errorf("invalid post ID: %v", err)
	}
	post, err := database.FetchPostByID(postIDInt)
	if err != nil {
		return post_logic.Post{}, fmt.Errorf("failed to fetch post: %v", err)
	}
	return post, nil
}




