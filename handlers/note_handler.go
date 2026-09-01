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

type SimpleNote struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func GetNotesByUserID(db *sql.DB, userID int) ([]*Note, error) {
	rows, err := db.Query("SELECT * FROM notes WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}

	var notes []*Note
	for rows.Next() {
		note := Note{}
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.UserID)
		if err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	return notes, nil
}

func getNoteByID(db *sql.DB, id int) (*Note, error) {
	row := db.QueryRow("SELECT * FROM notes WHERE id = ?", id)
	note := Note{}
	err := row.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.UserID)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func CreateNote(db *sql.DB, title string, content string, userID int) error {
	note := Note{Title: title, Content: content, UserID: userID}
	_, err := db.Exec("INSERT INTO notes (title, content, user_id) VALUES (?, ?, ?)", note.Title, note.Content, note.UserID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateNote(db *sql.DB, note *Note) error {
	_, err := db.Exec("UPDATE notes SET title = ?, content = ? WHERE id = ?", note.Title, note.Content, note.ID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteNote(db *sql.DB, noteID int) error {
	_, err := db.Exec("DELETE FROM notes WHERE id = ?", noteID)
	if err != nil {
		return err
	}
	return nil
}
