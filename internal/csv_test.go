package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTempCSV(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	file := filepath.Join(dir, "cookies.csv")

	err := os.WriteFile(file, []byte(content), 0644)
	require.NoError(t, err)

	return file
}

func TestLoadCookiesFileNotFound(t *testing.T) {
	cookies, err := LoadCookies("invalid.csv")

	require.Error(t, err)
	assert.Nil(t, cookies)
}

func TestLoadCookiesNoDataRows(t *testing.T) {
	file := createTempCSV(t, "cookie,timestamp\n")

	cookies, err := LoadCookies(file)

	require.Error(t, err)
	assert.EqualError(t, err, "no data rows in csv")
	assert.Nil(t, cookies)
}

func TestLoadCookiesInvalidTimestamp(t *testing.T) {
	file := createTempCSV(t,
		`cookie,timestamp
AtY0laUfhglK3lC7,invaliddate
`)

	cookies, err := LoadCookies(file)

	require.Error(t, err)
	assert.Contains(t, err.Error(), `invalid timestamp "invaliddate"`)
	assert.Nil(t, cookies)
}

func TestLoadCookiesValidSingleRow(t *testing.T) {
	file := createTempCSV(t,
		`cookie,timestamp
AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00
`)

	cookies, err := LoadCookies(file)

	require.NoError(t, err)
	require.Len(t, cookies, 1)

	expectedDate := time.Date(
		2018,
		time.December,
		9,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	assert.Equal(t, Cookie{
		ID:  "AtY0laUfhglK3lC7",
		Day: expectedDate,
	}, cookies[0])
}

func TestLoadCookiesValidMultipleRows(t *testing.T) {
	file := createTempCSV(t,
		`cookie,timestamp
AtY0laUfhglK3lC7,2018-12-09T14:19:00+00:00
FtY0laUfhglK3lC7,2018-12-08T23:59:59+00:00
`)

	cookies, err := LoadCookies(file)

	require.NoError(t, err)
	require.Len(t, cookies, 2)

	assert.Equal(t, "AtY0laUfhglK3lC7", cookies[0].ID)
	assert.Equal(t, "FtY0laUfhglK3lC7", cookies[1].ID)
}

func TestLoadCookiesNormalizesTimeToMidnightUTC(t *testing.T) {
	file := createTempCSV(t,
		`cookie,timestamp
AtY0laUfhglK3lC7,2018-12-09T23:59:59+00:00
`)

	cookies, err := LoadCookies(file)

	require.NoError(t, err)
	require.Len(t, cookies, 1)

	expected := time.Date(
		2018,
		time.December,
		9,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	assert.True(t, cookies[0].Day.Equal(expected))
	assert.Equal(t, time.UTC, cookies[0].Day.Location())
	assert.Equal(t, 0, cookies[0].Day.Hour())
	assert.Equal(t, 0, cookies[0].Day.Minute())
	assert.Equal(t, 0, cookies[0].Day.Second())
}
