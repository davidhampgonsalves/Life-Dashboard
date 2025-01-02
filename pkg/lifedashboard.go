package pkg

import (
	"bytes"
	"davidhampgonsalves/lifedashboard/pkg/draw"
	"davidhampgonsalves/lifedashboard/pkg/event"
	"davidhampgonsalves/lifedashboard/pkg/fetchers"
	"fmt"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

func Generate() []byte {
	fetchers := []func() ([]event.Event, error){
		fetchers.Weather,
		fetchers.GoogleCalendar("davidhampgonsalves@gmail.com"),
		fetchers.GoogleCalendar("ms7011nsnge4elr2cgvrmhap6g@group.calendar.google.com"),
		fetchers.SurfCaptain,
		fetchers.Tide,
		fetchers.SunAndMoon,
		fetchers.OvalSkating,

		// fetchers.SchoolClosures, // they changed the location, check back when canceled
		// fetchers.NsPower,
		// fetchers.Surfline,
	}

	events := []event.Event{}
	for _, fetcher := range fetchers {
		fetchedEvents, err := fetcher()
		if err == nil {
			events = append(events, fetchedEvents...)
		} else {
			fmt.Printf("> Fetcher error: %s\n", err)
		}
	}

	c, ctx, font := draw.Init()

	draw.Background(ctx)

	yPos := draw.Events(ctx, font, events)

	fmt.Printf("done events@%f\n", yPos)
	if yPos > 200.0 {
		draw.Date(ctx, font, yPos)
	}

	var imgData bytes.Buffer
	renderers.PNG(canvas.DPMM(1))(&imgData, c)

	return imgData.Bytes()
}
