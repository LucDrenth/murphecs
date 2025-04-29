package tick

import (
	"github.com/lucdrenth/murph_engine/src/app"
	"github.com/lucdrenth/murph_engine/src/engine/schedule"
)

type Counter struct {
	Count uint
}

func Init(app *app.BasicSubApp) {
	app.AddResource(&Counter{})
	app.AddSystem(schedule.PreUpdate, updateCounter)
}

func updateCounter(counter *Counter) {
	counter.Count++
}
