package main

import (
	"SimpleUserBackend/handlers"
	"database/sql"
	"encoding/json"
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

func loadDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT UNIQUE, password_hash TEXT)")
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
		users, err := handlers.GetUsers(db)
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

// TODO: move to a function
func httpCreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var user handlers.User
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				return
			}
			// TODO: Password hash is not working
			fmt.Println(user.Name, user.PasswordHash)
			_, err = db.Exec("INSERT INTO users (name, password_hash) VALUES (?, ?)", user.Name, user.PasswordHash)
			if err != nil {
				_, err := w.Write([]byte(err.Error()))
				if err != nil {
					return
				}
				return
			}
			w.Write([]byte("User created"))
		}
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
	http.HandleFunc("/user/create", httpCreateUser(db))

	fmt.Println("server started on port ", port)

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
