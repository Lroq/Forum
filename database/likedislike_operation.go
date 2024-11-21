package database

import (
	"Forum/post_logic"
	"database/sql"
	"log"
)

func UpdateLikeDislike(userUUID string, postID int, isLike, isDislike bool) error {
	_, err := DB.Exec("INSERT OR REPLACE INTO likes (user_uuid, post_id, is_like, is_dislike) VALUES (?, ?, ?, ?)", userUUID, postID, isLike, isDislike)
	return err
}

func IncrementLikeCounter(postID int) error {
	_, err := DB.Exec("UPDATE posts SET like_counter = like_counter + 1 WHERE id = ?", postID)
	return err
}

func IncrementDislikeCounter(postID int) error {
	_, err := DB.Exec("UPDATE posts SET dislike_counter = dislike_counter + 1 WHERE id = ?", postID)
	return err
}

func FetchLikesByUserUUID(userUUID string) ([]post_logic.Like, error) {
	rows, err := DB.Query("SELECT id, post_title FROM likes WHERE user_uuid = ? AND is_like = 1", userUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []post_logic.Like
	for rows.Next() {
		var like post_logic.Like
		var postTitle sql.NullString
		if err := rows.Scan(&like.ID, &postTitle); err != nil {
			return nil, err
		}
		if postTitle.Valid {
			like.PostTitle = postTitle.String
		} else {
			like.PostTitle = ""
		}
		likes = append(likes, like)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return likes, nil
}

func FetchDislikesByUserUUID(userUUID string) ([]post_logic.Like, error) {
	rows, err := DB.Query("SELECT id, COALESCE(post_title, ''), post_id FROM likes WHERE user_uuid = ? AND is_dislike = 1", userUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dislikes []post_logic.Like
	for rows.Next() {
		var d post_logic.Like
		var postTitle sql.NullString

		if err := rows.Scan(&d.ID, &postTitle, &d.PostID); err != nil {
			log.Println("Error scanning dislike: ", err)
			return nil, err
		}

		if postTitle.Valid {
			d.PostTitle = postTitle.String
		} else {
			d.PostTitle = ""
		}

		dislikes = append(dislikes, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dislikes, nil
}

func GetLikeCountByPostID(postID int) (int, error) {
	var likeCount int
	err := DB.QueryRow("SELECT like_counter FROM posts WHERE id = ?", postID).Scan(&likeCount)
	return likeCount, err
}