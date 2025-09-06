package fetchers

import (
	"errors"
	"net/http"
	"bytes"
	"encoding/json"
	"io"
	"fmt"

	"davidhampgonsalves/lifedashboard/pkg/event"
	"davidhampgonsalves/lifedashboard/pkg/utils"
)

const SystemPrompt = "Based on the weather forecast described in json generate a weather summary for todays with the most important details that is at most 80 characters long. Do not include wind information unless gusts are over 100km/h. Do not include humidex. Always start with the high/low temp range using the format \"low temp-high temp🌡️\" and do not include a unit symbol. If there is rain that day try and note any periods when it stops. Ignore fog information.Prefer terse summaries.Include a emoji at the start of the summary to characterize the days weather."

type ResponsePayload struct {
	Candidates []Candidate `json:"candidates"`
}
type Candidate struct {
	Content Content `json:"content"`
}
type Content struct {
	Parts []Part `json:"parts"`
}
type Part struct {
	Text string `json:"text"`
}


func jsonEscape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
			panic(err)
	}
	return string(b[1:len(b)-1])
}

func Weather() ([]event.Event, error) {
	resp, err := http.Get("https://weather.gc.ca/api/app/en/Location/44.649,-63.602?type=city")
	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("weather failed to load")
	}
	jsonBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("weather failed to read")
	}
	defer resp.Body.Close()

	bodyContent := fmt.Sprintf(`{ "contents": [ { "parts": [ { "text": "%s\n%s" } ] } ] }`, jsonEscape(SystemPrompt), jsonEscape(string(jsonBytes)))
	body := []byte(bodyContent)

	apiKey, _ := utils.ReadCredFile("gemini.txt")
	client := &http.Client{}

	req, _ := http.NewRequest("POST", "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent", bytes.NewBuffer(body))
	req.Header.Set("X-goog-api-key", apiKey)
	req.Header.Set("content-type", "application/json")
	resp, err = client.Do(req)

	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, errors.New("weather failed to read")
	}

	if err != nil || resp.StatusCode != 200 { return nil, errors.New("Gemini error") }
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	var payload ResponsePayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		fmt.Println("Error decoding gemini JSON:", err)
		return nil, errors.New("weather failed to read")
	}

	weather := event.Event{Text: payload.Candidates[0].Content.Parts[0].Text}
	return []event.Event{weather}, nil
}
