package schedule

import "github.com/lucdrenth/murph_engine/src/app"

const (
	Startup app.Schedule = "Startup"
	Cleanup app.Schedule = "Cleanup"

	// Core
	Update app.Schedule = "Update"

	// Engine
	Render app.Schedule = "Render"
)
