package fetchers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"regexp"

	// "strings"

	"davidhampgonsalves/lifedashboard/pkg/event"

	"github.com/PuerkitoBio/goquery"
)

func parseRange(input string) (int, error) {
	re := regexp.MustCompile(`^\d+`)
	match := re.FindString(input)
	if match == "" {
			return 0, fmt.Errorf("invalid input: %s", input)
	}
	num, err := strconv.Atoi(match)
	if err != nil {
			return 0, err
	}
	return num, nil
}

func SurfCaptain() ([]event.Event, error) {
	resp, err := http.Get("https://surfcaptain.com/forecast/cow-bay-nova-scotia")

	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("surfcaptain page failed to load")
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	str := ""
	doc.Find(".hourly-surf.clean").Each(func(_i int, node *goquery.Selection) {
		if str != "" {
			return
		}
		size, err := parseRange(node.Contents().Get(0).Data)
		if err != nil {
			return 
		}

		if size > 2 {
			day := node.ParentsFiltered(".fcst-summary-swell").Find(".summary-date-day").Text()
			str = fmt.Sprintf("🏄 %d ft on %s", size, day)
		}
	})

	if str == "" {
		return []event.Event{}, nil
	}
	surf := event.Event{Text: str}
	return []event.Event{surf}, nil
}
