package main

import (
	"SimpleUserBackend/routes"
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

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5500" || origin == "http://127.0.0.1:5500" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Add("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// For database
// go get modernc.org/sqlite

func loadDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbFile+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT UNIQUE, password_hash TEXT)")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS sessions (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER NOT NULL, token_hash TEXT NOT NULL UNIQUE, expires_at INTEGER NOT NULL, FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE)")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id)")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TRIGGER IF NOT EXISTS notes_set_updated_at
		AFTER UPDATE ON notes
		FOR EACH ROW
		WHEN NEW.updated_at = OLD.updated_at
		BEGIN
			UPDATE notes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
		END
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
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

	http.HandleFunc("/helloworld", routes.HttpTest())
	http.HandleFunc("/users", routes.HttpGetUsers(db))
	http.HandleFunc("/user/create", routes.HttpCreateUser(db))
	http.HandleFunc("/user/login", routes.HttpLogin(db))
	http.HandleFunc("/user/auth", routes.HttpAuthTest(db))

	fmt.Println("server started on port ", port)

	return http.ListenAndServe(fmt.Sprintf(":%v", port), cors(http.DefaultServeMux))
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
