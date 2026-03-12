package app

import "github.com/lucdrenth/murphecs/src/ecs"

type scheduleType int

const (
	ScheduleTypeStartup   scheduleType = iota // run only once, on startup
	ScheduleTypeRepeating                     // runs repeatedly, in the main loop
	ScheduleTypeCleanup                       // runs only once, before quitting
)

type ScheduleOptions struct {
	// ScheduleType can be one of:
	//   - [ScheduleTypeStartup] - systems in a schedule with this schedule type run once, when starting the app
	//   - [ScheduleTypeRepeating] - systems in a schedule with this schedule type run repeatedly, after startup
	//   - [ScheduleTypeCleanup] - systems in a schedule with this schedule type run once, when closing the app
	ScheduleType scheduleType

	// Order decides when the schedule systems should run, relative to schedules. It can be one of:
	//	- [ScheduleLast] - this is also the default if this field is left nil
	// 	- [ScheduleBefore]
	// 	- [ScheduleAfter]
	Order ecs.ScheduleOrder

	// IsPaused determines the initial pause state of the schedule
	IsPaused bool
}
