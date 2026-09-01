package handlers

import "database/sql"

type Note struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	UserID    int    `json:"user_id"`
}

func GetNotesByUserID(db *sql.DB, userID int) ([]*Note, error) {
	return nil, nil
}

func CreateNote(db *sql.DB, note *Note) error {
	return nil
}

func UpdateNote(db *sql.DB, note *Note) error {
	return nil
}

func DeleteNote(db *sql.DB, id int) error {
	return nil
}
