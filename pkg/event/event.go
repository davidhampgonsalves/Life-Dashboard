package event

import (
	"time"
)

type Event struct {
	Start time.Time
	End   time.Time
	Text  string
}
