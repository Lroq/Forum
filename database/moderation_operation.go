package database

import (
	"fmt"
)

func FetchAllUserUtilisateur() ([]User, error) {
	rows, err := DB.Query("SELECT uuid, username, email, role FROM users WHERE role = 'utilisateur'")
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

func FetchAllUserModerator() ([]User, error) {
	rows, err := DB.Query("SELECT uuid, username, email, role FROM users WHERE role = 'moderator'")
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

func FetchAllUserAdmin() ([]User, error) {
	rows, err := DB.Query("SELECT uuid, username, email, role FROM users WHERE role = 'ADMIN'")
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

func GetUserRoleByUsername(username string) (string, error) {
	var role string
	err := DB.QueryRow("SELECT role FROM users WHERE username = ?", username).Scan(&role)
	return role, err
}

func GetUserRoleByUUID(userUUID string) (string, error) {
	var role string
	err := DB.QueryRow("SELECT role FROM users WHERE uuid = ?", userUUID).Scan(&role)
	return role, err
}

func UpdateUserRole(userUUID string, role string) error {
	query := "UPDATE users SET role = ? WHERE uuid = ?"
	_, err := DB.Exec(query, role, userUUID)
	if err != nil {
		return fmt.Errorf("failed to update user role: %v", err)
	}
	return nil
}

func CreateAdminUser(uuid, username, email, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	_, err = DB.Exec("INSERT INTO users (uuid,username, email, password, role) VALUES (?,?, ?, ?, ?)", uuid, username, email, hashedPassword, "ADMIN")
	if err != nil {
		return fmt.Errorf("failed to insert admin user: %v", err)
	}
	return nil
}

