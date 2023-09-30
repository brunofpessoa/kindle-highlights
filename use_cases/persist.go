package useCases

import (
	"database/sql"
	"log"
	"strings"

	"github.com/brunofpessoa/kindle-highlights/domain"
	"github.com/brunofpessoa/kindle-highlights/repository"
	"github.com/brunofpessoa/kindle-highlights/util"
)

func BuildHighlights(fileContent string) (h []domain.Highlight) {
	fileContent = strings.ReplaceAll(fileContent, "\r", "")
	rawHighlights := strings.Split(fileContent, "==========")

	split := func(s string, sep string) (string, string) {
		x := strings.Split(s, sep)
		left := strings.Trim(x[0], "\n")
		right := strings.Trim(x[1], "\n")
		return left, right
	}

	for i := 0; i < len(rawHighlights)-1; i++ {
		info, content := split(rawHighlights[i], "\n\n")
		book, date := split(info, "\n-")
		isoDate, position := util.ExtractDateAndPosition(date)
		book = strings.TrimSpace(book)
		highlight := domain.Highlight{
			Book:     book,
			Date:     isoDate,
			Content:  content,
			Position: position,
		}
		h = append(h, highlight)
	}
	return
}

func PersistHighlights(db *sql.DB, file string, noDuplicates bool) {
	if !util.FileExists(file) {
		log.Fatal("unable to find clips file")
	}

	repository.CreateTables(db)
	fileContent := util.ReadFile(file)
	highlights := BuildHighlights(fileContent)
	repository.InsertData(db, highlights, noDuplicates)
}
