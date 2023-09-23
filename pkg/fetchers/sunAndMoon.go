package fetchers

import (
	"fmt"
	"math"
	"time"

	"davidhampgonsalves/lifedashboard/pkg/event"

	MoonPhase "github.com/janczer/goMoonPhase"
	"github.com/nathan-osman/go-sunrise"
)

var Emojis = map[int]string{
	0: "🌑",
	1: "🌒",
	2: "🌓",
	3: "🌔",
	4: "🌕",
	5: "🌖",
	6: "🌗",
	7: "🌘",
	8: "🌑",
}

func SunAndMoon() ([]event.Event, error) {
	hfx, _ := time.LoadLocation("America/Halifax")
	n := time.Now().In(hfx)
	rise, set := sunrise.SunriseSunset(44.64, -63.59, n.Year(), n.Month(), n.Day())
	m := MoonPhase.New(n)

	phase := int(math.Floor((m.Phase() + 0.0625) * 8))
	return []event.Event{{Text: fmt.Sprintf("🌞 🔼%s 🔽%s %s", rise.Format(time.Kitchen), set.Format(time.Kitchen), Emojis[phase])}}, nil
}
