package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

func Authenticate(db *sql.DB, r *http.Request) (bool, int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false, -1, err
	}
	if cookie == nil {
		return false, -1, nil
	}

	valid, i, err := VerifyUserToken(db, cookie.Value)
	if err != nil {
		return false, -1, err
	}

	if !valid {
		return false, -1, nil
	}
	fmt.Println("user id: ", i)

	return true, i, nil
}
