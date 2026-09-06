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

		authBool, userID, err := handlers.Authenticate(db, r)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !authBool {
			handlers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if r.Method != http.MethodPost {
			handlers.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		var simpleNote handlers.SimpleNote
		err = json.NewDecoder(r.Body).Decode(&simpleNote)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
		}

		err = handlers.CreateNote(db, simpleNote.Title, simpleNote.Content, userID)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
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

func HttpUpdateNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authBool, _, err := handlers.Authenticate(db, r)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !authBool {
			handlers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if r.Method != http.MethodPut {
			handlers.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		var noteRq handlers.Note
		err = json.NewDecoder(r.Body).Decode(&noteRq)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
		}

		if noteRq.ID == 0 {
			handlers.WriteError(w, http.StatusBadRequest, "ID is required")
			return
		}
		if noteRq.Title == "" && noteRq.Content == "" {
			handlers.WriteError(w, http.StatusBadRequest, "Title or Content is required")
			return
		}

		err = handlers.UpdateNote(db, &noteRq)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		handlers.WriteJSON(w, http.StatusOK, "Note updated")

	}
}

func HttpDeleteNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authBool, _, err := handlers.Authenticate(db, r)
		if err != nil {
			handlers.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !authBool {
			handlers.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if r.Method != http.MethodDelete {
			handlers.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

	}
}
