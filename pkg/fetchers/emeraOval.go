package fetchers

import (
	"errors"
	"net/http"
	"encoding/json"
	"bytes"
	"strings"
	"time"

	"davidhampgonsalves/lifedashboard/pkg/event"
)

type Schedule struct {
	Items []Item `json:"items"`
}

type Item struct {
	Summary string `json:"summary"`
	Start Start `json:"start"`
}

type Start struct {
	DateTime time.Time `json:"datetime"`
}

func OvalSkating() ([]event.Event, error) {
	loc, err := time.LoadLocation("America/Halifax")
	if err != nil { return nil, errors.New("atlantic timezone couldn't be loaded") }
	now := time.Now().In(loc).Truncate(24 * time.Hour)

	resp, err := http.Get("https://clients6.google.com/calendar/v3/calendars/g3bfd4h4ngthv403cn2i0lktdc%40group.calendar.google.com/events?calendarId=g3bfd4h4ngthv403cn2i0lktdc%40group.calendar.google.com&singleEvents=true&eventTypes=default&eventTypes=focusTime&eventTypes=outOfOffice&timeZone=America%2FHalifax&maxAttendees=1&maxResults=250&sanitizeHtml=true&timeMin=2024-12-30T00%3A00%3A00%2B18%3A00&timeMax=2025-01-29T00%3A00%3A00-18%3A00&key=AIzaSyBNlYH01_9Hc5S1J9vuFmu2nUqBZJNAXxs&%24unique=gc456")

	if err != nil || resp.StatusCode != 200 { return nil, errors.New("google calendar schedule for oval failed to load") }
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	var schedule Schedule
	json.Unmarshal(buf.Bytes(), &schedule)

	times :=  []string{}
	for _, item := range schedule.Items {
		if item.Start.DateTime.YearDay() != now.YearDay() { continue }
		// fmt.Printf("----- %+v\n", item.Summary)
		s := strings.ToLower(item.Summary)
		if strings.Contains(s, "speed") || strings.Contains(s, "maintenance") { continue }
		times = append(times, item.Start.DateTime.Format("3:04"))
	}

	if len(times) == 0 {
		return []event.Event{}, nil
	}
	surf := event.Event{Text: "⛸️ " + strings.Join(times, ", ")}
	return []event.Event{surf}, nil
}
