package database

import (
	"Forum/comment_logic"
	"fmt"
	"log"

	
)

func FetchCommentsByUsername(username string) ([]comment_logic.Comment, error) {
	rows, err := DB.Query("SELECT id, post_title, username, content FROM comments WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []comment_logic.Comment
	for rows.Next() {
		var c comment_logic.Comment
		if err := rows.Scan(&c.ID, &c.PostTitle, &c.Username, &c.Content); err != nil {
			log.Println("Error scanning comment: ", err)
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func FetchCommentsByPostID(postID int) ([]comment_logic.Comment, error) {
	rows, err := DB.Query("SELECT id, post_id, username, content FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []comment_logic.Comment
	for rows.Next() {
		var c comment_logic.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.Username, &c.Content); err != nil {
			log.Println("Error scanning comment: ", err)
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

// InsertComment creates a new comment and saves it in the database.
func InsertComment(postID int, username, content string) error {
	_, err := DB.Exec("INSERT INTO comments (post_id, username, content) VALUES (?, ?, ?)", postID, username, content)
	return err
}

func DeleteCommentByID(commentID int) error {
	query := "DELETE FROM comments WHERE id = ?"
	_, err := DB.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %v", err)
	}
	return nil
}



