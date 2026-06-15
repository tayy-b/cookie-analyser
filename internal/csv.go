package internal

//parses csv data and loads into cookies

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"time"
)

var ErrInvalidCSVDateFormat = errors.New("csv date is in invalid formated. expected ")

// opens csv file and loads data into a slice of cookies
func LoadCookies(fileName string) ([]Cookie, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("no data rows in csv")
	}

	var cookies []Cookie

	for _, row := range records[1:] {
		if len(row) < 2 {
			return nil, fmt.Errorf("invalid row: %#v", row)
		}

		date, err := normaliseTimestamp(row[1])
		if err != nil {
			return nil, err
		}
		cookies = append(cookies, Cookie{
			ID:  row[0],
			Day: date,
		})
	}

	return cookies, nil
}

// formats time stamps into year/month/day format
func normaliseTimestamp(val string) (time.Time, error) {
	timestamp, err := time.Parse(time.RFC3339, val)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"invalid timestamp %q: %w",
			val,
			err,
		)
	}
	year, month, day := timestamp.Date()
	date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return date, nil
}
