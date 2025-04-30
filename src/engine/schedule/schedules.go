package schedule

import "github.com/lucdrenth/murph_engine/src/app"

// Schedule that gets executed only once
const (
	Startup app.Schedule = "Startup"
	Cleanup app.Schedule = "Cleanup"
)

// Core schedule that gets executed in a loop
const (
	PreUpdate  app.Schedule = "PreUpdate"
	Update     app.Schedule = "Update"
	PostUpdate app.Schedule = "PostUpdate"
)

// Render schedule that gets executed in a loop
const (
	PreRender  app.Schedule = "PreRender"
	Render     app.Schedule = "Render"
	PostRender app.Schedule = "PostRender"
)
