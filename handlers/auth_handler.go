package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func Authenticate(db *sql.DB, r *http.Request) (bool, int, error) {
	cookie, err := r.Cookie("session_token")
	if errors.Is(err, http.ErrNoCookie) {
		return false, -1, nil
	}
	if err != nil {
		return false, -1, err
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
