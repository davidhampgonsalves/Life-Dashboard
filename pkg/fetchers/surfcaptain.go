package fetchers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"davidhampgonsalves/lifedashboard/pkg/event"

	"github.com/PuerkitoBio/goquery"
)

func SurfCaptain() ([]event.Event, error) {
	resp, err := http.Get("https://surfcaptain.com/forecast/cow-bay-nova-scotia")

	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("surfcaptain page failed to load")
	}
	// bodyBytes, err := io.ReadAll(resp.Body)
	// fmt.Println("BODY")
	// fmt.Println(string(bodyBytes))
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	fmt.Println()
	str := "🏄"

	// , .day-summary-cond
	doc.Find(".day-summary-surf").Each(func(i int, node *goquery.Selection) {
		if i > 0 {
			str += " "
		}
		str += strings.TrimSpace(node.Contents().Get(1).Data)
	})

	surf := event.Event{Text: str}
	return []event.Event{surf}, nil
}
