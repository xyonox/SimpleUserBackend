package main

import (
	"database/sql"
	"net/http"

	_ "modernc.org/sqlite"
)

const (
	dbFile = "SimpleUserBackend.db"
)

// For database
// go get modernc.org/sqlite

func loadDB() (*sql.DB, error) {
	return nil, nil
}

func httpTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, world!"))
		if err != nil {
			return
		}
	}
}

func run() error {

	/*db, err := loadDB()
	if err != nil {
		return err
	}
	if db == nil {
		return errors.New("db is nil")
	}
	*/

	http.HandleFunc("/helloworld", httpTest())

	return http.ListenAndServe(":8080", nil)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
