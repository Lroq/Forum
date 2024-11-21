package database

import (
	"fmt"
	"Forum/post_logic"
	"database/sql"
	"log"

	
)


func DeletePostByID(postID int) error {
	query := "DELETE FROM posts WHERE id = ?"
	_, err := DB.Exec(query, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %v", err)
	}
	return nil
}

func FetchPostsByUsername(username string) ([]post_logic.Post, error) {
	var posts []post_logic.Post
	rows, err := DB.Query("SELECT id, title, content, username, image, gif, like_counter, dislike_counter FROM posts WHERE username = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p post_logic.Post
		var image sql.NullString
		var gif sql.NullString

		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Username, &image, &gif, &p.LikesCount, &p.DislikesCount); err != nil {
			return nil, err
		}

		if image.Valid {
			p.Image = image.String
		} else {
			p.Image = ""
		}

		if gif.Valid {
			p.Gif = gif.String
		} else {
			p.Gif = ""
		}

		posts = append(posts, p)
	}
	return posts, nil
}

func FetchTopPostsByLikes() ([]string, error) {
	var titles []string

	query := `SELECT title FROM posts ORDER BY like_counter DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, err
		}
		titles = append(titles, title)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return titles, nil
}

func FetchPostsWithLikesDislikes() ([]post_logic.Post, error) {
	var posts []post_logic.Post
	rows, err := DB.Query(`
        SELECT p.id, p.title, p.content, p.username, 
               SUM(CASE WHEN l.is_like = 1 THEN 1 ELSE 0 END) AS likes_count,
               SUM(CASE WHEN l.is_dislike = 1 THEN 1 ELSE 0 END) AS dislikes_count
        FROM posts p
        LEFT JOIN likes l ON p.id = l.post_id
        GROUP BY p.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p post_logic.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Username, &p.LikesCount, &p.DislikesCount); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func FetchPostByID(id int) (post_logic.Post, error) {
	var p post_logic.Post
	err := DB.QueryRow("SELECT id, title, content, username, image, gif, like_counter, dislike_counter FROM posts WHERE id = ?", id).Scan(
		&p.ID, &p.Title, &p.Content, &p.Username, &p.Image, &p.Gif, &p.LikesCount, &p.DislikesCount)
	return p, err
}

func GetPosts() ([]post_logic.Post, error) {
	rows, err := DB.Query(`
        SELECT p.id, p.title, p.content, p.username, p.image, p.gif,
               COALESCE(SUM(CASE WHEN l.is_like = 1 THEN 1 ELSE 0 END), 0) AS likes_count,
               COALESCE(SUM(CASE WHEN l.is_dislike = 1 THEN 1 ELSE 0 END), 0) AS dislikes_count
        FROM posts p
        LEFT JOIN likes l ON p.id = l.post_id
        GROUP BY p.id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []post_logic.Post
	for rows.Next() {
		var p post_logic.Post
		var image sql.NullString
		var gif sql.NullString

		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Username, &image, &gif, &p.LikesCount, &p.DislikesCount); err != nil {
			log.Println("Error scanning post: ", err)
			return nil, err
		}

		if image.Valid {
			p.Image = image.String
		} else {
			p.Image = ""
		}

		if gif.Valid {
			p.Gif = gif.String
		} else {
			p.Gif = ""
		}

		posts = append(posts, p)

	}
	return posts, nil
}

func InsertPost(title, content, username, imagePath, gifPath string) error {
	tx, err := DB.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := "INSERT INTO posts (title, content, username, image, gif) VALUES (?, ?, ?, ?, ?)"
	result, err := tx.Exec(query, title, content, username, imagePath, gifPath)
	if err != nil {
		log.Printf("Error executing SQL query: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}

	log.Printf("Rows affected: %d", rowsAffected)
	return nil
}