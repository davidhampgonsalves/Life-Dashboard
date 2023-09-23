package fetchers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"davidhampgonsalves/lifedashboard/pkg/event"

	"github.com/PuerkitoBio/goquery"
)

func SchoolClosures() ([]event.Event, error) {
	resp, err := http.Get("https://www.hrce.ca/about-our-schools/parents/school-cancellations")
	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("school closures data failed to load")
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	title := doc.Find("#page-title").First().Text()
	m, _ := regexp.MatchString("close", strings.ToLower(title))

	if m {
		return []event.Event{event.Event{Text: "📚 Schools Closed"}}, nil
	}

	return []event.Event{}, nil
}
