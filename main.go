package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	_ "modernc.org/sqlite"
)

const (
	dbFile = "SimpleUserBackend.db"
	port   = 8080
)

// For database
// go get modernc.org/sqlite

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
}

func getUsers(db *sql.DB) ([]*User, error) {
	rows, err := db.Query("SELECT * FROM users")
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
		err := rows.Scan(&user.ID, &user.Name, &user.PasswordHash)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func getUser(db *sql.DB, id int) (*User, error) {
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	user := User{}
	err := row.Scan(&user.ID, &user.Name, &user.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func loadDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT, password_hash TEXT)")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func httpTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, world!"))
		if err != nil {
			return
		}
	}
}

func httpGetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := getUsers(db)
		if err != nil {
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				return
			}
			return
		}

		_, err = w.Write([]byte(fmt.Sprintf("%v", users)))
	}
}

func run() error {

	db, err := loadDB()
	if err != nil {
		return err
	}
	if db == nil {
		return errors.New("db is nil")
	}

	err = db.Ping()
	if err == nil {
		fmt.Println("PONG!")
	}

	fmt.Println("db is ready")

	fmt.Println("starting server")

	http.HandleFunc("/helloworld", httpTest())
	http.HandleFunc("/users", httpGetUsers(db))

	fmt.Println("server started on port ", port)

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
