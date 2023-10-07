package repository

import (
	"database/sql"
	"log"

	"github.com/brunofpessoa/kindle-highlights/domain"
	"github.com/brunofpessoa/kindle-highlights/util"
)

func CreateTables(db *sql.DB) {
	db.Exec(`
		CREATE TABLE IF NOT EXISTS books(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL
		);
	`)
	db.Exec(`
		CREATE TABLE IF NOT EXISTS highlights (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			content TEXT NOT NULL,
			date TEXT NOT NULL,
			position TEXT NOT NULL,
			book_id INTEGER NOT NULL,
			FOREIGN KEY (book_id) REFERENCES books(id)
		);
	`)
	db.Exec("CREATE INDEX books_index ON books (name)")
	db.Exec("CREATE INDEX highlights_index ON highlights (content)")
}

func InsertHighlight(tx *sql.Tx, h domain.Highlight, bookID int64) {
	_, err := tx.Exec(
		"INSERT INTO highlights (content, date, position, book_id) VALUES (?, ?, ?, ?)",
		h.Content,
		h.Date,
		h.Position,
		bookID,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertBook(tx *sql.Tx, h domain.Highlight) (bookID int64) {
	err := tx.QueryRow("SELECT id FROM books WHERE name = ?", h.Book).Scan(&bookID)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	// If the book is not duplicated, insert it and get the ID
	if err == sql.ErrNoRows {
		result, err := tx.Exec("INSERT INTO books (name) VALUES (?)", h.Book)
		if err != nil {
			log.Fatal(err)
		}
		bookID, _ = result.LastInsertId()
	}
	return
}

func IsContentDuplicated(tx *sql.Tx, h domain.Highlight, bookID int64) bool {
	var dupID, dupDate string
	dupErr := tx.QueryRow("SELECT id, date FROM highlights WHERE content LIKE ? AND book_id = ?", "%"+h.Content+"%", bookID).
		Scan(&dupID, &dupDate)

	// keep the last one if there is two similar highlights
	if dupErr != sql.ErrNoRows {
		shouldUpdate := util.IsFirstDateMoreRecent(h.Date, dupDate)
		if shouldUpdate {
			UpdateHighlight(tx, dupID, h)
		}
		return true
	}
	return false
}

func InsertData(db *sql.DB, highlights *[]domain.Highlight, noDuplicates bool) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	for _, h := range *highlights {
		bookID := InsertBook(tx, h)
		isDuplicated := IsContentDuplicated(tx, h, bookID)
		if isDuplicated && noDuplicates {
			continue
		}

		InsertHighlight(tx, h, bookID)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}

func UpdateHighlight(tx *sql.Tx, id string, h domain.Highlight) {
	_, err := tx.Exec(
		"UPDATE highlights SET content = ?, date = ?, position = ? WHERE id = ?",
		h.Content,
		h.Date,
		h.Position,
		id,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func getBookID(db *sql.DB, bookName string) (id int) {
	query := `SELECT id FROM books WHERE name LIKE ?;`

	err := db.QueryRow(query, "%"+bookName+"%").Scan(&id)
	if err != nil {
		log.Fatal("Unable to find book with name: " + bookName)
	}
	return
}

func GetRand(db *sql.DB, minLen int, maxLen int) (h domain.Highlight) {
	query := `
        SELECT content, date, name AS book, position
        FROM highlights AS h
        INNER JOIN books AS b ON h.book_id = b.id
        WHERE LENGTH(content) > ? AND LENGTH(content) < ?
        ORDER BY RANDOM()
        LIMIT 1;
    `
	err := db.QueryRow(query, minLen, maxLen).Scan(&h.Content, &h.Date, &h.Book, &h.Position)
	if err != nil {
		log.Fatal("unable to find any highlight with theses parameters")
	}
	return
}

func GetRandByBook(db *sql.DB, minLen int, maxLen int, bookName string) (h domain.Highlight) {
	bookID := getBookID(db, bookName)

	query := `
        SELECT content, date, name AS book, position
        FROM highlights AS h
        INNER JOIN books AS b ON h.book_id = b.id
        WHERE LENGTH(content) > ? AND LENGTH(content) < ? AND book_id = ?
        ORDER BY RANDOM()
        LIMIT 1;
    `

	err := db.QueryRow(query, minLen, maxLen, bookID).
		Scan(&h.Content, &h.Date, &h.Book, &h.Position)
	if err != nil {
		log.Fatal("unable to find any highlight with theses parameters")
	}

	return
}

func GetAll(db *sql.DB, minLen int, maxLen int) (highlights []domain.Highlight) {
	query := `
        SELECT content, date, name AS book, position
        FROM highlights AS h
        INNER JOIN books AS b ON h.book_id = b.id
        WHERE LENGTH(content) > ? AND LENGTH(content) < ?;
    `

	rows, err := db.Query(query, minLen, maxLen)
	for rows.Next() {
		var h domain.Highlight
		if err := rows.Scan(&h.Content, &h.Date, &h.Book, &h.Position); err != nil {
			log.Fatal(err)
		}
		highlights = append(highlights, h)
	}
	if err != nil {
		log.Fatal("unable to find any highlight with theses parameters")
	}

	return
}

func GetAllByBook(
	db *sql.DB,
	minLen int,
	maxLen int,
	bookName string,
) (highlights []domain.Highlight) {
	bookID := getBookID(db, bookName)

	query := `
        SELECT content, date, name AS book, position
        FROM highlights AS h
        INNER JOIN books AS b ON h.book_id = b.id
        WHERE LENGTH(content) > ? AND LENGTH(content) < ? AND book_id = ?;
    `

	rows, err := db.Query(query, minLen, maxLen, bookID)
	for rows.Next() {
		var h domain.Highlight
		if err := rows.Scan(&h.Content, &h.Date, &h.Book, &h.Position); err != nil {
			log.Fatal(err)
		}
		highlights = append(highlights, h)
	}
	if err != nil {
		log.Fatal("unable to find any highlight with theses parameters")
	}

	return
}
