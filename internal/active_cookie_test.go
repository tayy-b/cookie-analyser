package internal

import (
	"slices"
	"testing"
	"time"
)

type findMostActiveTest struct {
	description string
	cookies     []Cookie
	targetDate  time.Time
	expected    []string
}

func day(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return t.UTC()
}

var findMostActiveTestCases = []findMostActiveTest{
	{
		description: "single most active cookie",
		cookies: []Cookie{
			{ID: "AtY0laUfhglK3lC7", Day: day("2018-12-09")},
			{ID: "SAZuXPGUrfbcn5UA", Day: day("2018-12-09")},
			{ID: "5UAVanZf6UtGyKVS", Day: day("2018-12-09")},
			{ID: "AtY0laUfhglK3lC7", Day: day("2018-12-09")},
			{ID: "SAZuXPGUrfbcn5UA", Day: day("2018-12-08")},
			{ID: "4sMM2LxV07bPJzwf", Day: day("2018-12-08")},
			{ID: "fbcn5UAVanZf6UtG", Day: day("2018-12-08")},
			{ID: "4sMM2LxV07bPJzwf", Day: day("2018-12-08")},
		},
		targetDate: day("2018-12-09"),
		expected:   []string{"AtY0laUfhglK3lC7"},
	},
	{
		description: "multiple cookies most active",
		cookies: []Cookie{
			{ID: "AtY0laUfhglK3lC7", Day: day("2018-12-09")},
			{ID: "SAZuXPGUrfbcn5UA", Day: day("2018-12-09")},
			{ID: "5UAVanZf6UtGyKVS", Day: day("2018-12-08")},
			{ID: "5UAVanZf6UtGyKVS", Day: day("2018-12-08")},
			{ID: "4sMM2LxV07bPJzwf", Day: day("2018-12-08")},
			{ID: "4sMM2LxV07bPJzwf", Day: day("2018-12-08")},
			{ID: "AtY0laUfhglK3lC7", Day: day("2018-12-08")},
		},
		targetDate: day("2018-12-08"),
		expected:   []string{"4sMM2LxV07bPJzwf", "5UAVanZf6UtGyKVS"},
	},
	{
		description: "no cookies for requested day",
		cookies: []Cookie{
			{ID: "AtY0laUfhglK3lC7", Day: day("2018-12-08")},
			{ID: "5UAVanZf6UtGyKVS", Day: day("2018-12-08")},
		},
		targetDate: day("2018-12-09"),
		expected:   []string{},
	},
}

func TestFindMostActive(t *testing.T) {
	for _, tc := range findMostActiveTestCases {
		t.Run(tc.description, func(t *testing.T) {
			actual := FindMostActive(tc.cookies, tc.targetDate)

			slices.Sort(actual)
			slices.Sort(tc.expected)

			if !slices.Equal(actual, tc.expected) {
				t.Errorf(
					"FindMostActive() = %#v, want %#v",
					actual,
					tc.expected,
				)
			}
		})
	}
}

func BenchmarkFindMostActive(b *testing.B) {
	targetDate := day("2018-12-09")

	var cookies []Cookie
	d10 := day("2018-12-10")
	d09 := day("2018-12-09")
	d08 := day("2018-12-08")

	for i := 0; i < 5000; i++ {
		cookies = append(cookies,
			Cookie{ID: "a", Day: d10},
		)
	}

	for i := 0; i < 5000; i++ {
		cookies = append(cookies,
			Cookie{ID: "b", Day: d09},
		)
	}

	for i := 0; i < 5000; i++ {
		cookies = append(cookies,
			Cookie{ID: "c", Day: d08},
		)
	}

	b.ResetTimer()

	for b.Loop() {
		FindMostActive(cookies, targetDate)
	}
}
