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
			handlers.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		responseBody := make(map[string]any)

		var loginRequest handlers.User
		err := json.NewDecoder(r.Body).Decode(&loginRequest)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		user, err := handlers.GetUserByName(db, loginRequest.Name)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		passwordVerify, err := handlers.VerifyPassword(loginRequest.PasswordHash, user.PasswordHash)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if !passwordVerify {
			handlers.WriteError(w, http.StatusUnauthorized, "Password verification failed")
			return
		}

		token, err := handlers.SetUsersToken(db, user.ID, 1)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
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
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func HttpAuthTest(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		responseBody := make(map[string]any)

		authBool, i, err := handlers.Authenticate(db, r)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !authBool {
			handlers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		w.WriteHeader(http.StatusOK)
		responseBody["message"] = "token is valid"
		jsonBody, _ := json.Marshal(responseBody)
		_, err = w.Write(jsonBody)
		fmt.Println("user id: ", i)
	}
}

func HttpGetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		users, err := handlers.GetUsers(db)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		type UserNonPassword struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		var usersNonPassword []UserNonPassword
		for _, user := range users {
			usersNonPassword = append(usersNonPassword, UserNonPassword{
				ID:   user.ID,
				Name: user.Name,
			})
		}

		jsonBody, err := json.Marshal(usersNonPassword)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonBody)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

// TODO: move to a function
func HttpCreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		responseBody := make(map[string]any)

		if r.Method == http.MethodPost {
			var user handlers.User
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				handlers.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}

			hashedPassword, err := handlers.HashPassword(user.PasswordHash)
			if err != nil {
				handlers.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			user.PasswordHash = hashedPassword

			_, err = db.Exec("INSERT INTO users (name, password_hash) VALUES (?, ?)", user.Name, user.PasswordHash)
			if err != nil {
				handlers.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
			w.WriteHeader(http.StatusCreated)
			responseBody["message"] = "User created"
			jsonBody, _ := json.Marshal(responseBody)
			_, err = w.Write(jsonBody)
			if err != nil {
				handlers.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			handlers.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}
	}
}
