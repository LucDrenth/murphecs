package features

import (
	"github.com/lucdrenth/murph_engine/src/app"
)

type TickCounterFeature struct {
	app.Feature
	Schedule app.Schedule
}

type TickCounter struct {
	Count uint
}

func (f *TickCounterFeature) Init() {
	f.AddResource(&TickCounter{})
	f.AddSystem(f.Schedule, updateCounter)
}

func updateCounter(counter *TickCounter) {
	counter.Count++
}
