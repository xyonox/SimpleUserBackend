package main

import (
	"SimpleUserBackend/handlers"
	"SimpleUserBackend/routes"
	"database/sql"
	"errors"
	"fmt"
	"math/bits"
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

	// TODO: Login logic => token, simple auth test route

	fmt.Println("server started on port ", port)

	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

/*func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
*/

func main() {

	passwords := []string{
		"WhatsUp",
		"whatsup",
		"WHATSUP",
		"wHAtsuP",
		"whbtsup",
		"whatsvp",
		"xhatsup",
		"zhatsup",
		"xhbtups",
	}

	originalHash, salt, err := handlers.HashPassword(passwords[0])
	if err != nil {
		panic(err)
	}

	fmt.Printf("Original: %s\n", passwords[0])
	fmt.Printf("Hash:     %x\n\n", originalHash)
	fmt.Println(passwords[0])

	for _, password := range passwords[1:] {

		hash, err := handlers.HashPasswordWithSalt(password, salt)
		if err != nil {
			panic(err)
		}

		differentBits := 0

		minLength := len(originalHash)
		if len(hash) < minLength {
			minLength = len(hash)
		}

		for i := 0; i < minLength; i++ {
			differentBits += bits.OnesCount8(originalHash[i] ^ hash[i])
		}

		totalBits := len(originalHash) * 8
		percentage := float64(differentBits) / float64(totalBits) * 100

		fmt.Printf("Passwort: %s\n", password)
		fmt.Printf("Hash:     %x\n", hash)
		fmt.Println(hash)
		fmt.Printf("Unterschiedliche Bits: %d / %d (%.2f%%)\n\n",
			differentBits,
			totalBits,
			percentage,
		)
	}
}
