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

const SystemPrompt = "Based on the weather forecast described in json generate a weather summary for todays with the most important details that is at most 80 characters long. Do not include wind information unless gusts are over 100km/h. Always start with the high/low temp range using the format \"low temp-high temp🌡️\" and do not include a unit symbol. If there is rain that day try and note any periods when it stops. Ignore fog information.Prefer terse summaries.End summary with a period.Include a emoji at the start of the summary to characterize the days weather."
type Content struct {
	Text string `json:"text"`
}

type AnthropicResponse  struct {
	Content []Content `json:"content"`
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

	body := []byte(
		fmt.Sprintf(`{ "model": "claude-3-5-sonnet-20241022", "system": "%s", "max_tokens": 1024, "messages": [ { "role": "user", "content": "%s" } ] }`, 
		jsonEscape(SystemPrompt), 
		jsonEscape(string(jsonBytes))))

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(body))
	apiKey, _ := utils.ReadCredFile("anthropic.txt")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")
	resp, err = client.Do(req)

	if err != nil || resp.StatusCode != 200 { return nil, errors.New("anthropic request failed") }
	defer resp.Body.Close()

	anthropicResponse:= &AnthropicResponse{}
	err = json.NewDecoder(resp.Body).Decode(anthropicResponse)
	if err != nil { return nil, errors.New("anthropic response could not be decoded") }
	
	weather := event.Event{Text: anthropicResponse.Content[0].Text}
	return []event.Event{weather}, nil
}
