package handlers


import (
	"Forum/comment_logic"
	"Forum/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func createComment(postID string, username string, content string) error {
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		return fmt.Errorf("invalid post ID: %v", err)
	}
	err = database.InsertComment(postIDInt, username, content)
	if err != nil {
		return fmt.Errorf("failed to insert comment: %v", err)
	}
	return nil
}

func UserCommentsHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to get username from session", http.StatusInternalServerError)
		return
	}

	comments, err := database.FetchCommentsByUsername(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user comments: %v", err), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}


func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Extract comment content and post ID from form data
	content := r.FormValue("content")
	postID := r.FormValue("postID")
	// Get the username from the session
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Error(w, "Failed to get username from session", http.StatusInternalServerError)
		return
	}
	// Create and save the new comment
	err = createComment(postID, username, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	database.UpdateComments()
	// Redirect back to the post detail page
	http.Redirect(w, r, fmt.Sprintf("/posts/%s", postID), http.StatusFound)
}

func fetchCommentsForPost(postID string) ([]comment_logic.Comment, error) {
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		return nil, fmt.Errorf("invalid post ID: %v", err)
	}
	comments, err := database.FetchCommentsByPostID(postIDInt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %v", err)
	}
	return comments, nil
}
