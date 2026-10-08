package fetchers

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"os"
	"strings"
	"time"

	"davidhampgonsalves/lifedashboard/pkg/event"

	"github.com/PuerkitoBio/goquery"
)

const lunchFile = "lunch.eml"

func Lunch() ([]event.Event, error) {
	f, err := os.Open(lunchFile)
	if err != nil {
		return nil, fmt.Errorf("could not open %s", lunchFile)
	}
	defer f.Close()

	msg, err := mail.ReadMessage(f)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid email", lunchFile)
	}

	html, err := lunchHtmlPart(msg.Header.Get("Content-Type"), msg.Header.Get("Content-Transfer-Encoding"), msg.Body)
	if err != nil {
		return nil, err
	}
	if html == "" {
		return nil, errors.New("lunch email has no html part")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var events []event.Event
	seen := map[string]bool{}

	doc.Find("table tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 2 {
			return
		}
		date := normalizeSpace(cells.Eq(cells.Length() - 2).Text())
		meal := normalizeSpace(cells.Eq(cells.Length() - 1).Text())
		if meal == "" || !isLunchToday(date, now) || seen[meal] {
			return
		}
		seen[meal] = true
		events = append(events, event.Event{Text: fmt.Sprintf("🍲 %s", meal)})
	})

	return events, nil
}

func lunchHtmlPart(contentType, encoding string, body io.Reader) (string, error) {
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		parts := multipart.NewReader(body, params["boundary"])
		for {
			part, err := parts.NextPart()
			if err == io.EOF {
				return "", nil
			}
			if err != nil {
				return "", err
			}
			html, err := lunchHtmlPart(part.Header.Get("Content-Type"), part.Header.Get("Content-Transfer-Encoding"), part)
			if err == nil && html != "" {
				return html, nil
			}
		}
	}

	if mediaType != "text/html" {
		return "", nil
	}

	if strings.EqualFold(strings.TrimSpace(encoding), "quoted-printable") {
		body = quotedprintable.NewReader(body)
	}
	html, err := io.ReadAll(body)

	return string(html), err
}

func isLunchToday(date string, now time.Time) bool {
	weekday, monthDay, found := strings.Cut(date, ",")
	if !found || strings.TrimSpace(weekday) != now.Weekday().String() {
		return false
	}

	d, err := time.Parse("January 2", strings.TrimSpace(monthDay))
	if err != nil {
		return false
	}

	return d.Month() == now.Month() && d.Day() == now.Day()
}

func normalizeSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
