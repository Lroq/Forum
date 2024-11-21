package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type User struct {
	UUID           string
	Username       string
	Email          string
	HashedPassword string
	Role           string
}

type Report struct {
	ID            int
	ModeratorUUID string
	ModeratorName string
	Content       string
	Resolved      bool
	Response      string
	PostID        int
	PostTitle     string
	PostContent   string
}

type Role struct {
	ID        int
	GUEST     string
	USER      string
	Moderator string
	ADMIN     string
}

func InitDB(DBPath string) {
	var err error
	DB, err = sql.Open("sqlite3", DBPath)
	if err != nil {
		log.Fatal(err)
	}
	createTable := `CREATE TABLE IF NOT EXISTS users (
        uuid TEXT PRIMARY KEY,
        username TEXT NULL,
        email TEXT NULL,
        password TEXT UNIQUE,
        role TEXT DEFAULT 'GUEST'
    );`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
	createPostsTable := `CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        content TEXT,
        username TEXT,
		image TEXT,
		gif TEXT,
		like_counter INTEGER DEFAULT 0,
		dislike_counter INTEGER DEFAULT 0,
        FOREIGN KEY(username) REFERENCES users(username)
    );`
	_, err = DB.Exec(createPostsTable)
	if err != nil {
		log.Fatal(err)
	}
	createTableLikes := `CREATE TABLE IF NOT EXISTS likes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_uuid TEXT,
        post_id INTEGER,
        is_like BOOLEAN,
        is_dislike BOOLEAN,
		post_title TEXT,
        FOREIGN KEY(user_uuid) REFERENCES users(uuid),
        FOREIGN KEY(post_id) REFERENCES posts(id),
		FOREIGN KEY(post_title) REFERENCES posts(title),
        UNIQUE(user_uuid, post_id)
    );`
	_, err = DB.Exec(createTableLikes)
	if err != nil {
		log.Fatal(err)
	}

	updateLikes := `UPDATE likes
	SET post_title = (
    SELECT title 
    FROM posts 
    WHERE posts.id = likes.post_id
	);`
	_, err = DB.Exec(updateLikes)
	if err != nil {
		log.Fatal(err)
	}

	createTableComments := `CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER,
		username TEXT,
		content TEXT,
		post_title TEXT,
		FOREIGN KEY(post_id) REFERENCES posts(id)
	);`
	_, err = DB.Exec(createTableComments)
	if err != nil {
		log.Fatal(err)
	}

	createTableReports := `CREATE TABLE IF NOT EXISTS reports (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		moderator_uuid INTEGER,
		moderator_name TEXT,
		content TEXT,
		response TEXT DEFAULT 'Pas encore de réponse..',
		resolved BOOLEAN DEFAULT FALSE,
		post_id INTEGER,
		post_title TEXT,
		post_content TEXT,
		FOREIGN KEY (moderator_name) REFERENCES users(username)
    	FOREIGN KEY (post_id) REFERENCES posts(id)
		FOREIGN KEY (post_title) REFERENCES posts(title)
		FOREIGN KEY (post_content) REFERENCES posts(content)
	);`
	_, err = DB.Exec(createTableReports)
	if err != nil {
		log.Fatal(err)
	}

}

func CloseDB() {
	DB.Close()
}

func UpdateComments() error {
	_, err := DB.Exec(`UPDATE comments
	SET post_title = (
	SELECT title 
	FROM posts 
	WHERE posts.id = comments.post_id
	);`)
	if err != nil {
		return err
	}
	return nil
}