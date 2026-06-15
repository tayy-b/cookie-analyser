package internal

// Business logic for finding the most active cookie command - Returns cookie ids for all matches.
// Assumes csv file is ordered.

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var ErrInvalidDateFormat = errors.New("Invalid date provided. must be in YYYY-MM-DD format")

type Cookie struct {
	ID  string
	Day time.Time
}

func RunCookieFinder(cmd *cobra.Command, args []string) error {
	fileName, _ := cmd.Flags().GetString("file")
	date, _ := cmd.Flags().GetString("date")
	logrus.Debugf("recieved file name %s and date %s", fileName, date)

	targetDate, err := parseDateFromString(date)

	if err != nil {
		return ErrInvalidDateFormat
	}

	cookies, err := LoadCookies(fileName)
	logrus.Debugf("%+v", cookies)

	if err != nil {
		return err
	}

	most_active_cookies := FindMostActive(cookies, targetDate)
	if len(most_active_cookies) == 0 {
		fmt.Println("no cookies found on given day")
		return nil
	}
	logrus.Debugf("%+v", most_active_cookies)

	for _, cookie_ids := range most_active_cookies {
		fmt.Println(cookie_ids)
	}
	return nil

}

// returns a slice of most active cookie ids
func FindMostActive(cookies []Cookie, targetDate time.Time) []string {
	cookieCount := make(map[string]int)
	maxCount := 0
	for _, c := range cookies {

		if targetDate.After(c.Day) {
			break //we can stop if the target date is after the current cookie since it's ordered from most recent
		}
		logrus.Debugf("checking cookie %s", c.ID)
		if c.Day.Equal(targetDate) {
			cookieCount[c.ID]++
			if cookieCount[c.ID] > maxCount {
				maxCount = cookieCount[c.ID]
			}
		}
	}
	mostActiveCookies := []string{}
	for id, count := range cookieCount {
		if count == maxCount {
			mostActiveCookies = append(mostActiveCookies, id)
		}
	}
	return mostActiveCookies
}

// parses user input into the correct format
func parseDateFromString(date string) (time.Time, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, date)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}
