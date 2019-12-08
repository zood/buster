package resources

import (
	"log"
	"time"
)

// Post contains all the information about a blog/news item, except the actual body
type Post struct {
	ID    int    `json:"id"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Intro string `json:"intro"`
	Date  int64  `json:"date"`
}

// MonthDayYear returns the post date in PST time because that's where Zood is located
func (p Post) MonthDayYear() string {
	t := time.Unix(p.Date, 0)
	pst, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Printf("Failed to find America/Los_Angeles time zone: %v", err)
		return t.Format("January 2, 2006 MST")
	}

	t = t.In(pst)
	return t.Format("January 2, 2006")
}

type postSlice []Post

func (ps postSlice) Len() int {
	return len(ps)
}

func (ps postSlice) Less(i, j int) bool {
	return ps[i].ID > ps[j].ID
}

func (ps postSlice) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}
