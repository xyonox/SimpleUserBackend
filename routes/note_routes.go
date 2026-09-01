package routes

import (
	"SimpleUserBackend/handlers"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func HttpCreateNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		responseBody := make(map[string]any)

		// TODO implement auth
		authBool, userID, err := handlers.Authenticate(db, r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			responseBody["error"] = fmt.Sprintf("Authentication failed: %s", err.Error())
			jsonBody, err := json.Marshal(responseBody)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(err.Error()))
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				fmt.Println("Error: ", err)
				return
			}
			_, err = w.Write(jsonBody)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			return
		}

		if !authBool {
			w.WriteHeader(http.StatusUnauthorized)
			responseBody["error"] = "Authentication failed"
			jsonBody, err := json.Marshal(responseBody)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(err.Error()))
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				fmt.Println("Error: ", err)
				return
			}
			_, err = w.Write(jsonBody)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			responseBody["error"] = "Method not allowed"
			jsonBody, err := json.Marshal(responseBody)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, err := w.Write([]byte(err.Error()))
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				fmt.Println("Error: ", err)
				return
			}
			_, err = w.Write(jsonBody)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}

		var simpleNote handlers.SimpleNote
		err = json.NewDecoder(r.Body).Decode(&simpleNote)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			responseBody["error"] = "Method not allowed"
			jsonBody, err := json.Marshal(responseBody)
			if err != nil {
				_, err := w.Write([]byte(err.Error()))
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				fmt.Println("Error: ", err)
				return
			}
			_, err = w.Write(jsonBody)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}

		err = handlers.CreateNote(db, simpleNote.Title, simpleNote.Content, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			responseBody["error"] = "Method not allowed"
			jsonBody, err := json.Marshal(responseBody)
			if err != nil {
				_, err := w.Write([]byte(err.Error()))
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				fmt.Println("Error: ", err)
				return
			}
			_, err = w.Write(jsonBody)
			if err != nil {
				fmt.Println("Error: ", err)
			}
		}
		w.WriteHeader(http.StatusCreated)
		responseBody["message"] = "Note created"
		jsonBody, _ := json.Marshal(responseBody)
		_, err = w.Write(jsonBody)
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
