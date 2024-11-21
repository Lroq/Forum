package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)


func InsertUser(username, email, password, role string) error {
	// Générer l'UUID
	userUUID := uuid.New().String()

	// Vérifiez d'abord si l'utilisateur existe déjà
	exists, err := UserExists(username, email)
	if err != nil {
		return err // Renvoie l'erreur si la vérification échoue
	}
	if exists {
		return fmt.Errorf("un utilisateur avec ce nom d'utilisateur ou email existe déjà")
	}

	// Continuez avec l'insertion si l'utilisateur n'existe pas
	if role == "" {
		role = "GUEST"
	}
	_, err = DB.Exec("INSERT INTO users (uuid, username, email, password, role) VALUES (?, ?, ?, ?, ?)", userUUID, username, email, password, role)
	return err
}

func UserExists(username, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = ? OR email = ?)`
	err := DB.QueryRow(query, username, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func GetUserID(email string) (string, error) {
	var uuid string
	err := DB.QueryRow("SELECT uuid FROM users WHERE email = ?", email).Scan(&uuid)
	return uuid, err
}

func GetEmailPassword(email string) (string, error) {
	var hashedPassword string
	err := DB.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&hashedPassword)
	return hashedPassword, err
}

func GetUsernameByEmail(email string) (string, error) {
	var username string
	err := DB.QueryRow("SELECT username FROM users WHERE email = ?", email).Scan(&username)
	return username, err
}

func GetIdByUsername(username string) (string, error) {
	var uuid string
	err := DB.QueryRow("SELECT uuid FROM users WHERE username = ?", username).Scan(&uuid)
	return uuid, err
}

func FetchUserByUUID(uuid string) (*User, error) {
	var user User
	var email, password sql.NullString

	err := DB.QueryRow("SELECT uuid, username, email, password, role FROM users WHERE uuid = ?", uuid).Scan(
		&user.UUID, &user.Username, &email, &password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error fetching user by UUID: %w", err)
	}

	if email.Valid {
		user.Email = email.String
	} else {
		user.Email = ""
	}

	if password.Valid {
		user.HashedPassword = password.String
	} else {
		user.HashedPassword = ""
	}

	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	var password sql.NullString

	err := DB.QueryRow("SELECT uuid, username, email, password, role FROM users WHERE email = ?", email).Scan(
		&user.UUID, &user.Username, &user.Email, &password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error fetching user by email: %w", err)
	}

	if password.Valid {
		user.HashedPassword = password.String
	} else {
		user.HashedPassword = ""
	}

	return &user, nil
}

func CreateGuestUser() (string, error) {
	// Générer un nouvel UUID pour l'utilisateur
	userUUID := uuid.New().String()

	// Insérer le nouvel utilisateur avec son UUID dans la base de données
	_, err := DB.Exec("INSERT INTO Users (uuid, role) VALUES (?, 'GUEST')", userUUID)
	if err != nil {
		return "", err
	}

	// Créer un nom d'utilisateur basé sur l'UUID
	username := fmt.Sprintf("invité%s", userUUID[:8]) // Utiliser les 8 premiers caractères de l'UUID pour le nom d'utilisateur

	// Mettre à jour le nom d'utilisateur dans la base de données en utilisant l'UUID
	_, err = DB.Exec("UPDATE Users SET username = ? WHERE uuid = ?", username, userUUID)
	if err != nil {
		return "", err
	}

	return userUUID, nil
}

func GetAllUsers() ([]User, error) {
	rows, err := DB.Query("SELECT uuid, username, email, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.UUID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func UpdateUserProfile(userUUID string, username, email string) error {
	_, err := DB.Exec("UPDATE users SET username = ?, email = ? WHERE uuid = ?", username, email, userUUID)
	return err
}

func DeleteUser(uuid string) error {
	_, err := DB.Exec("DELETE FROM users WHERE uuid = ?", uuid)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	return nil
}
