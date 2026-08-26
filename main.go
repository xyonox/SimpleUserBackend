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
	http.HandleFunc("/user/auth", routes.HttpAuthTest())

	// TODO: Login logic => token, simple auth test route

	fmt.Println("server started on port ", port)

	return http.ListenAndServe(fmt.Sprintf(":%v", port), cors(http.DefaultServeMux))
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
