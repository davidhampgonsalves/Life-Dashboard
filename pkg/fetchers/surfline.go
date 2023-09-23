package fetchers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"davidhampgonsalves/lifedashboard/pkg/event"
)

type Rating struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type Json struct {
	Data struct {
		Rating []struct {
			Timestamp int64  `json:"timestamp"`
			UtcOffset int    `json:"utcOffset"`
			Rating    Rating `json:"rating"`
		} `json:"rating"`
	} `json:"data"`
}

func Surfline() ([]event.Event, error) {
	resp, err := http.Get("https://services.surfline.com/kbyg/spots/forecasts/rating?spotId=584204204e65fad6a77094cb&days=3&intervalHours=1&cacheEnabled=true")

	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("surfline data failed to load")
	}
	defer resp.Body.Close()

	var report Json
	err = json.NewDecoder(resp.Body).Decode(&report)
	if err != nil {
		return nil, errors.New("surfline data failed parse JSON")
	}

	str := "🏄"
	summary := make(map[string]Rating)
	for _, rating := range report.Data.Rating {
		t := time.Unix(rating.Timestamp, 0)
		day := t.Format("Jan _2")
		prev := summary[day]
		if prev.Value < rating.Rating.Value {
			summary[day] = rating.Rating
		} else {
			summary[day] = prev
		}
	}

	for day, rating := range summary {
		str = fmt.Sprintf("%s %s %s,", str, day, strings.ToLower(strings.ReplaceAll(rating.Key, "_", " ")))
	}

	surf := event.Event{Text: strings.Trim(str, ",")}
	return []event.Event{surf}, nil
}
