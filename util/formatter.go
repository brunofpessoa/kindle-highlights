package util

import (
	"fmt"
	"strings"
	"time"
)

var months = map[string]string{
	// pt-BR
	"janeiro":   "01",
	"fevereiro": "02",
	"março":     "03",
	"abril":     "04",
	"maio":      "05",
	"junho":     "06",
	"julho":     "07",
	"agosto":    "08",
	"setembro":  "09",
	"outubro":   "10",
	"novembro":  "11",
	"dezembro":  "12",
	// en
	"january":   "01",
	"february":  "02",
	"march":     "03",
	"april":     "04",
	"may":       "05",
	"june":      "06",
	"july":      "07",
	"august":    "08",
	"september": "09",
	"october":   "10",
	"november":  "11",
	"december":  "12",
}

func ExtractDateAndPosition(date string) (formattedDate string, position string) {
	rawParts := strings.Split(date, ", ")

	rawPosition := strings.Split(rawParts[0], " ")
	datePart := rawParts[1]

	replacedDate := strings.ReplaceAll(datePart, "de ", "")
	dateParts := strings.Split(replacedDate, " ")

	day := dateParts[0]
	month := months[strings.ToLower(dateParts[1])]
	year := dateParts[2]
	time := dateParts[3]
	formattedDate = fmt.Sprintf("%s-%s-%s %s", year, month, day, time)

	indicatorPart := rawPosition[4]
	positionPart := rawPosition[5]

	position = fmt.Sprintf("%s: %s", indicatorPart, positionPart)

	return
}

func IsFirstDateMoreRecent(dateStr1, dateStr2 string) bool {
	layout := "2006-01-02 15:04:05"
	date1, _ := time.Parse(layout, dateStr1)
	date2, _ := time.Parse(layout, dateStr2)
	if date1.After(date2) {
		return true
	}
	return false
}
