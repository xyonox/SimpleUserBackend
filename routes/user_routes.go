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
			fmt.Println("Error: ", err)
			return
		}
	}
}

func HttpLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{Name: "test", Value: "123"})

		fmt.Println("test")
		w.WriteHeader(http.StatusOK)
	}
}

func HttpAuthTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		coockie, err := r.Cookie("test")
		if err != nil {
			_, err := w.Write([]byte("no cookie"))
			if err != nil {
				return
			}
			return
		}
		_, err = w.Write([]byte(coockie.Value))
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Println("cookie: ", coockie.Value)
	}
}

func HttpGetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		users, err := handlers.GetUsers(db)
		if err != nil {
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				return
			}
			return
		}

		type UserNonPassword struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		usersNonPassword := []UserNonPassword{}
		for _, user := range users {
			usersNonPassword = append(usersNonPassword, UserNonPassword{
				ID:   user.ID,
				Name: user.Name,
			})
		}

		jsonBody, err := json.Marshal(usersNonPassword)

		_, err = w.Write(jsonBody)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
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

			hashedPassword, err := handlers.HashPassword(user.PasswordHash)
			if err != nil {
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
			}
			user.PasswordHash = hashedPassword

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
