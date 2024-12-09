package fetchers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"davidhampgonsalves/lifedashboard/pkg/event"

	"github.com/PuerkitoBio/goquery"
)

func toUnicode(observationType string) string {
	ot := strings.ToLower(observationType)

	m, _ := regexp.MatchString("snow|freezing|ice|squalls", ot)
	if m {
		return "❄️"
	}
	m, _ = regexp.MatchString("rain|mist|precipitation|drizzle|thunder", ot)
	if m {
		return "🌂"
	}
	m, _ = regexp.MatchString("fog|haze", ot)
	if m {
		return "🌫"
	}
	m, _ = regexp.MatchString("clear", ot)
	if m {
		return "☀️"
	}
	m, _ = regexp.MatchString("partly cloud", ot)
	if m {
		return "⛅"
	}
	m, _ = regexp.MatchString("cloud", ot)
	if m {
		return "☁"
	}

	return "🦞"
}

func Weather() ([]event.Event, error) {
	resp, err := http.Get("https://weather.gc.ca/city/pages/ns-19_metric_e.html")
	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("weather failed to load")
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	observationType, _ := doc.Find("img[data-v-33b01f9c]").First().Attr("alt")
	description := doc.Find(".pdg-tp-0").First().Find("td").Last().Text()
	weather := event.Event{Text: fmt.Sprintf("%s %s", toUnicode(observationType), description)}
	return []event.Event{weather}, nil
}
