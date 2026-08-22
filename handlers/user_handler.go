package handlers

import (
	"database/sql"
	"fmt"
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

func GetUser(db *sql.DB, id int) (*User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	user := User{}
	err := row.Scan(&user.ID, &user.Name, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
