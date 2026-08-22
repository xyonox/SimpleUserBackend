package routes

import (
	"SimpleUserBackend/handlers"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func HttpTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, world!"))
		if err != nil {
			return
		}
	}
}

func HttpGetUsers(db *sql.DB) http.HandlerFunc {
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
func HttpCreateUser(db *sql.DB) http.HandlerFunc {
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
			_, err = w.Write([]byte("User created"))
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}
	}
}
