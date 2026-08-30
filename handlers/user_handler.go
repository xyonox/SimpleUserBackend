package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
}

func GetUsers(db *sql.DB) ([]*User, error) {
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println("error closing rows")
		}
	}(rows)

	users := []*User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func GetUser(db *sql.DB, id int) (User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	user := User{}
	err := row.Scan(&user.ID, &user.Name, &user.PasswordHash)
	if err != nil {
		return user, err
	}
	return user, nil
}

func SetUsersToken(db *sql.DB, id int, daysOfExpiring int) (string, error) {
	token, err := GenerateToken()
	if err != nil {
		return "", err
	}

	_, err = db.Exec("DELETE FROM sessions WHERE user_id = ?", id)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(time.Duration(daysOfExpiring) * time.Hour * 24)

	hashedToken := HashToken(token)

	_, err = db.Exec("INSERT INTO sessions (user_id, token_hash, expires_at) VALUES (?, ?, ?)", id, hashedToken, expiresAt.Unix())
	return token, err
}

func VerifyUserToken(db *sql.DB, token string) (bool, int, error) {

	var userID int
	err := db.QueryRow(
		`SELECT user_id
     FROM sessions
     WHERE token_hash = ? AND expires_at > ?`,
		HashToken(token),
		time.Now().Unix(),
	).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, -1, nil
	}
	if err != nil {
		return false, -1, err
	}
	return true, userID, nil
}

func GetUserByName(db *sql.DB, name string) (User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE name = ?", name)
	user := User{}
	err := row.Scan(&user.ID, &user.Name, &user.PasswordHash)
	if err != nil {
		return user, err
	}
	return user, nil
}
