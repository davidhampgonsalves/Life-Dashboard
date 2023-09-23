package fetchers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"davidhampgonsalves/lifedashboard/pkg/event"

	"github.com/PuerkitoBio/goquery"
)

func Surf() ([]event.Event, error) {
	resp, err := http.Get("https://www.ndbc.noaa.gov/station_page.php?station=44258")
	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("bouy data failed to load")
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	period := "n/a"
	height := "n/a"
	doc.Find(".currentobs td").Each(func(i int, s *goquery.Selection) {
		if i == 10 {
			period = strings.TrimSpace(s.Text())
		}
		if i == 8 {
			height = strings.TrimSpace(s.Text())
		}
	})

	surf := event.Event{Text: fmt.Sprintf("🌊🛟 %s @ %s", height, period)}
	return []event.Event{surf}, nil
}
