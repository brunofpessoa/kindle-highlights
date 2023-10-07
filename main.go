package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	useCases "github.com/brunofpessoa/kindle-highlights/use_cases"
	"github.com/brunofpessoa/kindle-highlights/util"
)

func main() {
	var dbFileName, fileName, bookName string
	var minLen, maxLen int
	var showAll, noDuplicates bool

	flag.IntVar(&maxLen, "max", 500, "Max length of the highlight")
	flag.IntVar(&minLen, "min", 10, "Min length of the highlight")
	flag.StringVar(&dbFileName, "db", "database.sqlite3", "Database name")
	flag.StringVar(&fileName, "file", "My Clippings.txt", "File path")
	flag.StringVar(
		&bookName,
		"book",
		"",
		"Specifies which book to take the highlight from. Can be the full name or part of it",
	)
	flag.BoolVar(&showAll, "all", false, "Show all highlights")
	flag.BoolVar(
		&noDuplicates,
		"no-duplicates",
		true,
		"Determine if should remove duplicated highlights. If true will keep the last one created",
	)

	flag.Parse()

	if minLen > maxLen {
		log.Fatal("min length must be less than max length")
	}

	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if !util.FileExists(dbFileName) {
		fmt.Println("Wait, creating database...")
		useCases.PersistHighlights(db, fileName, noDuplicates)
		fmt.Println("Done! Run again to print your highlights")
		return
	}

	if bookName != "" && showAll {
		useCases.PrintAllByBook(db, minLen, maxLen, bookName)
	} else if bookName != "" {
		useCases.PrintRandByBook(db, minLen, maxLen, bookName)
	} else if showAll {
		useCases.PrintAll(db, minLen, maxLen)
	} else {
		useCases.PrintRand(db, minLen, maxLen)
	}
}
