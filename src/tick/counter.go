package tick

import (
	"github.com/lucdrenth/murph_engine/src/app"
)

type Counter struct {
	Count uint
}

func Init(app *app.BasicSubApp, schedule app.Schedule) {
	app.AddResource(&Counter{})
	app.AddSystem(schedule, updateCounter)
}

func updateCounter(counter *Counter) {
	counter.Count++
}
