package routes

import (
	"SimpleUserBackend/handlers"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func HttpTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Hello, world!"))
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}
}

func HttpLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodPost {
			_, err := w.Write([]byte("Method not allowed"))
			if err != nil {
				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Println("Error: ", err)
				return
			}
			return
		}

		responseBody := make(map[string]any)

		var sentUser handlers.User
		err := json.NewDecoder(r.Body).Decode(&sentUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			responseBody["error"] = err.Error()
			jsonBody, _ := json.Marshal(responseBody)
			_, err := w.Write(jsonBody)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}

		user, err := handlers.GetUserByName(db, sentUser.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			responseBody["error"] = err.Error()
			jsonBody, _ := json.Marshal(responseBody)
			_, err := w.Write(jsonBody)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			return
		}

		passwordVerify, err := handlers.VerifyPassword(sentUser.PasswordHash, user.PasswordHash)
		if err != nil {
			return
		}

		if !passwordVerify {
			responseBody["error"] = "password is not correct"
			jsonBody, _ := json.Marshal(responseBody)
			_, err := w.Write(jsonBody)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			return
		}

		token, err := handlers.SetUsersToken(db, user.ID, 1)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			responseBody["error"] = err.Error()
			jsonBody, _ := json.Marshal(responseBody)
			_, err := w.Write(jsonBody)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // lokal über HTTP; in Produktion mit HTTPS auf true setzen
			SameSite: http.SameSiteLaxMode,
			MaxAge:   60 * 60 * 24, // 24 Stunden
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    strconv.Itoa(user.ID),
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   60 * 60 * 24, // 24 Stunden
		})

		responseBody["message"] = "Login successful"

		jsonBody, _ := json.Marshal(responseBody)
		_, err = w.Write(jsonBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Error: ", err)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func HttpAuthTest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		responseBody := make(map[string]any)

		cookie, err := r.Cookie("session_token")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			responseBody["error"] = err.Error()
			jsonBody, _ := json.Marshal(responseBody)
			_, err := w.Write(jsonBody)
			if err != nil {
				return
			}
			return
		}

		valid, i, err := handlers.VerifyUserToken(db, cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			responseBody["error"] = err.Error()
			jsonBody, _ := json.Marshal(responseBody)
			_, err := w.Write(jsonBody)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			return
		}

		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			responseBody["error"] = "token is not valid"
			jsonBody, _ := json.Marshal(responseBody)
			_, err := w.Write(jsonBody)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("token is valid"))
		fmt.Println("user id: ", i)
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
		w.Header().Set("Content-Type", "application/json")
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
