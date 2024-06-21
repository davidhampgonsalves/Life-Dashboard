package fetchers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

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

	str := "🌊"
	doc.Find(".day-tables").First().Find("tr.even, tr.odd").Each(func(i int, row *goquery.Selection) {
		time := row.Find("td").First().Text()
		height, _ := strconv.ParseFloat(row.Find("td:nth-child(2)").First().Text(), 64)

		if i == 0 {
			if height > 1.0 {
				str += "⏫"
			} else {
				str += "⏬"
			}
			str += strings.TrimLeft(time, "0")
		} 
	})

	return []event.Event{{Text: str}}, nil
}
