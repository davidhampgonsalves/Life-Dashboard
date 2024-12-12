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
	resp, err := http.Get("https://weather.gc.ca/en/location/index.html?coords=44.649,-63.602")
	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("weather failed to load")
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	observationType, _ := doc.Find("img.mrgn-tp-md").First().Attr("alt")
	high := doc.Find(".mrgn-lft-sm[title=High]").First().Text()
	low := doc.Find(".mrgn-lft-sm[title=Low]").First().Text()

	rawDescription := doc.Find(".pdg-tp-0").First().Find("td").Last().Text()
	
	re := regexp.MustCompile(`([^.]+\.[^.]+\.)`)
	match := re.FindStringSubmatch(rawDescription)

	weather := event.Event{Text: fmt.Sprintf("%s %s/%s🌡️, %s", toUnicode(observationType), high, low, match[0])}
	return []event.Event{weather}, nil
}
