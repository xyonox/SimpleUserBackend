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

	fmt.Println("server started on port ", port)

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
