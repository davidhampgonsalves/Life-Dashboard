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
	"path/filepath"
	"strings"
	"time"

	"davidhampgonsalves/lifedashboard/pkg/event"

	"github.com/PuerkitoBio/goquery"
)

func Lunch() ([]event.Event, error) {
	lunchDir := os.Getenv("LUNCH_DIR")
	if lunchDir == "" {
		lunchDir = "."
	}

	files, err := filepath.Glob(filepath.Join(lunchDir, "*.eml"))
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no eml files found in %s", lunchDir)
	}

	now := time.Now()
	var events []event.Event
	var errs []error
	seen := map[string]bool{}

	for _, file := range files {
		meals, err := lunchMeals(file, now)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for _, meal := range meals {
			if seen[meal] {
				continue
			}
			seen[meal] = true
			events = append(events, event.Event{Text: fmt.Sprintf("🍲 %s", meal)})
		}
	}

	if len(events) == 0 && len(errs) == len(files) {
		return nil, errors.Join(errs...)
	}

	return events, nil
}

func lunchMeals(file string, now time.Time) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("could not open %s", file)
	}
	defer f.Close()

	msg, err := mail.ReadMessage(f)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid email", file)
	}

	html, err := lunchHtmlPart(msg.Header.Get("Content-Type"), msg.Header.Get("Content-Transfer-Encoding"), msg.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	if html == "" {
		return nil, fmt.Errorf("%s has no html part", file)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}

	var meals []string
	doc.Find("table tr").Each(func(_ int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 2 {
			return
		}
		date := normalizeSpace(cells.Eq(cells.Length() - 2).Text())
		meal := normalizeSpace(cells.Eq(cells.Length() - 1).Text())
		if meal == "" || !isLunchToday(date, now) {
			return
		}
		meals = append(meals, meal)
	})

	return meals, nil
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
