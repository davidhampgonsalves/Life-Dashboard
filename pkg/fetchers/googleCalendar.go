package fetchers

import (
	"context"
	"davidhampgonsalves/lifedashboard/pkg/event"
	"fmt"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GoogleCalendar(calendarName string) func() ([]event.Event, error) {
	return func() ([]event.Event, error) {
		ctx := context.Background()
		googleCalendar, err := calendar.NewService(ctx, option.WithCredentialsFile("creds/jwt.keys.json"))
		if err != nil {
			return nil, err
		}

		year, month, day := time.Now().Date()
		hfx, _ := time.LoadLocation("America/Halifax")
		startTime := time.Date(year, month, day, 0, 0, 0, 0, hfx)
		endTime := startTime.AddDate(0, 0, 1)

		calEvents, err := googleCalendar.Events.List(calendarName).TimeMin(startTime.Format(time.RFC3339)).TimeMax(endTime.Format(time.RFC3339)).ShowDeleted(false).SingleEvents(true).MaxResults(10).OrderBy("startTime").Do()

		if err != nil {
			return nil, err
		}

		var events []event.Event
		for _, item := range calEvents.Items {
			span := item.Start.DateTime
			if span != "" {
				startTime, err := time.Parse(time.RFC3339, span)
				if err != nil {
					return nil, err
				}
				span = startTime.Format(time.Kitchen) + " "
			}
			events = append(events, event.Event{Text: fmt.Sprintf("🗓️%s%s", span, item.Summary)})
		}

		return events, nil
	}
}
