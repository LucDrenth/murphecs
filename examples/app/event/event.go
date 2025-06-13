// Demonstrate how to use the event system.
//
// Events are structs that embed [app.Event]. They are written using an EventWriter and read using an
// EventReader, which we declare as system parameters.
//
// A Written event can be read from the next schedule until the end of the
// next run of the schedule it is written in.
// That means that they are not available for reading in the same schedule that they are written in.
//
// Written events can be read by multiple EventReader's.
//
// Because reacting to a written event can only be done in the next schedule, they are not useful for
// immediately reacting to something.
// Use them if you only want to react to an event at a specific time. For example, lets take a WindowResize
// event. We wouldn't immediately want to respond to such an events if we are halfway through a render
// cycle because then our render would get all inconsistent. We do want to respond to such an event at the
// beginning of the render cycle, so we'd put our system with the EventWriter there.
package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/lucdrenth/murphecs/examples/app/run"
	"github.com/lucdrenth/murphecs/src/app"
	"github.com/lucdrenth/murphecs/src/ecs"
)

const (
	update app.Schedule = "Update"
)

type highRollEvent struct {
	app.Event
	number int
}

func main() {
	var logger app.Logger = &app.SimpleConsoleLogger{}
	myApp, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}
	myApp.SetTickRate(time.Millisecond * 100)

	myApp.
		AddSchedule(update, app.ScheduleTypeRepeating).
		AddResource(&logger)

	myApp.
		AddSystem(update, rollDice).
		AddSystem(update, logHighRolls)

	run.RunApps(myApp)
}

// rollDice rolls a 10-sided dice. If the dice lands on 10, 11 or 12, we send a highRollEvent.
func rollDice(eventWriter *app.EventWriter[*highRollEvent]) {
	number := rand.IntN(12) + 1
	if number >= 10 {
		eventWriter.Write(&highRollEvent{number: number})
	}
}

func logHighRolls(logger app.Logger, eventReader *app.EventReader[*highRollEvent]) {
	for highRollEvent := range eventReader.Read {
		logger.Info(fmt.Sprintf("high dice roll: %d", highRollEvent.number))
	}
}
