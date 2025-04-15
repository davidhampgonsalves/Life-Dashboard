package fetchers

import (
	"errors"
	"net/http"
	"bytes"
	"encoding/json"
	"io"
	"fmt"
	"strings"
	"time"

	"davidhampgonsalves/lifedashboard/pkg/event"
	"davidhampgonsalves/lifedashboard/pkg/utils"
)

const SystemPrompt = "Summarize todays weather in less that 80 chars. Always start with the high/low temp range using the format \"low temp-high temp🌡️\" and do not include a unit symbol. If there is rain that day note any periods when it stops. End summary with a period. Start with single emoji which characterizes the days weather."
type Content struct {
	Text string `json:"text"`
}

type Response  struct {
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

	weatherJson, err := getHourlyForecastAsString(string(jsonBytes))
	if err != nil {
		return nil, errors.New("failed to filter weather json")
	}
	body := []byte(fmt.Sprintf(`{ "content": "%s", "with_clean_history": true }`, jsonEscape(SystemPrompt + weatherJson))) 
	// fmt.Println(weatherJson)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://gitlab.com/api/v4/chat/completions", bytes.NewBuffer(body))
	apiKey, _ := utils.ReadCredFile("gitlab.txt")
	req.Header.Set("Authorization", "Bearer " + apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)

if err != nil || resp.StatusCode != 201 {
    errorBody, _ := io.ReadAll(resp.Body)
    fmt.Println("Error response:", string(errorBody))
    return nil, fmt.Errorf("gitlab request failed: %v, status: %d", err, resp.StatusCode)
	}
	summaryBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("weather failed to read")
	}
	defer resp.Body.Close()

	summary := string(summaryBytes)
	summary = strings.Trim(summary, "\"")
	summary = strings.ReplaceAll(summary, "\\u0026", "&")

	weather := event.Event{Text: summary}
	return []event.Event{weather}, nil
}

func getHourlyForecastAsString(jsonStr string) (string, error) {
	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			return "", fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	var result strings.Builder
	
	result.WriteString("Hour, Condition, Temp(°C), Precip\n")
	if len(data) > 0 {
			if hourlyFcst, ok := data[0]["hourlyFcst"].(map[string]interface{}); ok {
					if hourly, ok := hourlyFcst["hourly"].([]interface{}); ok {
							for _, h := range hourly {
									hourMap, ok := h.(map[string]interface{})
									if !ok {
											continue
									}

									epochTime := int64(hourMap["epochTime"].(float64))
									condition := hourMap["condition"].(string)
									precip := hourMap["precip"].(string)
									
									t := time.Unix(epochTime, 0)
									timeStr := t.Format("15")
									
									tempMap := hourMap["temperature"].(map[string]interface{})
									temperature := tempMap["metric"].(string)
									
									// Append the formatted output to the result string
									result.WriteString(fmt.Sprintf("%s,%s,%s,%s\n", timeStr, condition, temperature, precip))
							}
					}
			}
	}
	
	return result.String(), nil
}