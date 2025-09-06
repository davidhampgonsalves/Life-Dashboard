package fetchers

import (
	"errors"
	"net/http"
	"io/ioutil"
	"strings"
	"fmt"

	"davidhampgonsalves/lifedashboard/pkg/event"
	"davidhampgonsalves/lifedashboard/pkg/utils"
)

func Todos() ([]event.Event, error) {
	token, _ := utils.ReadCredFile("silver-bullet.txt")

	req, _ := http.NewRequest("GET", "http://home.davidhampgonsalves.com:3000/.fs/todo.md", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("content-type", "text/markdown")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != 200 { return nil, errors.New("silver bullet request failed") }
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	todos := strings.ReplaceAll(string(body), "*", "")

	var events []event.Event
	for _, todo := range strings.Split(todos, "\n") {
		events = append(events, event.Event{Text: todo})
	}

	return events, nil

	// return []event.Event{}, nil
}
