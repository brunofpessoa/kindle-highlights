package useCases

import (
	"database/sql"

	"github.com/brunofpessoa/kindle-highlights/domain"
	"github.com/brunofpessoa/kindle-highlights/repository"
	"github.com/fatih/color"
)

func print(h domain.Highlight) {
	book := color.New(color.FgBlue).Add(color.Bold)
	book.Printf("\n%s | %s | %s\n\n", h.Book, h.Position, h.Date)

	content := color.New(color.FgGreen).Add(color.Bold)
	content.Printf("%s\n\n", h.Content)
}

func PrintRand(db *sql.DB, minLen int, maxLen int) {
	h := repository.GetRand(db, minLen, maxLen)
	print(h)
}

func PrintAll(db *sql.DB, minLen int, maxLen int) {
	highlights := repository.GetAll(db, minLen, maxLen)
	for _, h := range highlights {
		print(h)
	}
}

func PrintRandByBook(db *sql.DB, minLen int, maxLen int, bookName string) {
	h := repository.GetRandByBook(db, minLen, maxLen, bookName)
	print(h)
}

func PrintAllByBook(db *sql.DB, minLen int, maxLen int, bookName string) {
	highlights := repository.GetAllByBook(db, minLen, maxLen, bookName)
	for _, h := range highlights {
		print(h)
	}
}
