package fetchers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"davidhampgonsalves/lifedashboard/pkg/event"

	"github.com/PuerkitoBio/goquery"
)

func Tide() ([]event.Event, error) {
	resp, err := http.Get("https://tides.gc.ca/en/stations/490")

	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("tide page failed to load")
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.New("tide page failed to parse")
	}

	hfx, _ := time.LoadLocation("America/Halifax")
	now := time.Now().In(hfx)

	str := ""
	doc.Find(".day-tables").EachWithBreak(func(_ int, table *goquery.Selection) bool {
		id, _ := table.Attr("id")
		date := strings.TrimPrefix(id, "day-table-")
		table.Find("tr.even, tr.odd").EachWithBreak(func(_ int, row *goquery.Selection) bool {
			t, err := time.ParseInLocation("2006-01-02 15:04", date+" "+strings.TrimSpace(row.Find("td").First().Text()), hfx)
			if err != nil || t.Before(now) {
				return true
			}
			height, _ := strconv.ParseFloat(strings.TrimSpace(row.Find("td:nth-child(2)").First().Text()), 64)
			if height > 1.0 {
				str = "🌊⏫" + t.Format(time.Kitchen)
			} else {
				str = "🌊⏬" + t.Format(time.Kitchen)
			}
			return false
		})
		return str == ""
	})

	if str == "" {
		return nil, errors.New("no upcoming tide found")
	}
	return []event.Event{{Text: str}}, nil
}
