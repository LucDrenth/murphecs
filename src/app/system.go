package app

import "github.com/lucdrenth/murph_engine/src/ecs"

type System func(app SubApp, params ...SystemParam) error

type SystemParam interface {
	// Fetch gets all data necessary to execute the system
	Prepare(*ecs.World)
	Run(*ecs.World)
}

type SystemSet struct {
	systems []System
}

func (s *SystemSet) Run(app SubApp) {
	for i := range s.systems {
		s.systems[i](app)
	}
}
