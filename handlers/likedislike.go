package handlers

import (
	"Forum/database"
	"net/http"
	"strconv"
	"encoding/json"
	"fmt"
	"log"
)


func LikeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_uuid"] == nil || session.IsNew {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	postID := r.FormValue("postID")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}
	userUUID, ok := session.Values["user_uuid"].(string)
	if !ok {
		http.Error(w, "Failed to get user ID from session", http.StatusInternalServerError)
		return
	}
	err = database.UpdateLikeDislike(userUUID, postIDInt, true, false)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = database.IncrementLikeCounter(postIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil || session.Values["user_uuid"] == nil || session.IsNew {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	postID := r.FormValue("postID")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}
	userUUID, ok := session.Values["user_uuid"].(string)
	if !ok {
		http.Error(w, "Failed to get user ID from session", http.StatusInternalServerError)
		return
	}
	err = database.UpdateLikeDislike(userUUID, postIDInt, false, true)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = database.IncrementDislikeCounter(postIDInt)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UserLikesHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		http.Error(w, "Error retrieving session", http.StatusInternalServerError)
		return
	}

	if session.Values["user_uuid"] == nil {
		log.Println("User UUID not found in session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userUUID, ok := session.Values["user_uuid"].(string)
	if !ok {
		log.Println("Failed to get user UUID from session")
		http.Error(w, "Failed to get user ID from session", http.StatusInternalServerError)
		return
	}

	likes, err := database.FetchLikesByUserUUID(userUUID)
	if err != nil {
		log.Printf("Error fetching user likes e: %v", err)
		http.Error(w, "Error fetching user likes", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(likes)
	if err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func UserDislikesHandler(w http.ResponseWriter, r *http.Request) {
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
	dislikes, err := database.FetchDislikesByUserUUID(userUUID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching user dislikes: %v", err), http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(dislikes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}