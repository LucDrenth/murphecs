package main

import (
	"fmt"
	"time"

	"github.com/lucdrenth/murphecs/examples/app/run"
	"github.com/lucdrenth/murphecs/src/app"
	"github.com/lucdrenth/murphecs/src/ecs"
)

type emptyComponentA struct{ ecs.Component }
type emptyComponentB struct{ ecs.Component }
type emptyComponentC struct{ ecs.Component }
type emptyComponentD struct{ ecs.Component }
type componentWithValue struct {
	ecs.Component
	value int
}

const startup app.Schedule = "Startup"
const update app.Schedule = "Update"

type ticks struct {
	total int64
}

type timeStarted struct {
	millis int64
}

type lastPrintTime struct {
	millis int64
}

func main() {
	var logger app.Logger = &app.NoOpLogger{}
	myApp, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}

	myApp.UseUncappedRunner()

	myApp.AddSchedule(startup, app.ScheduleTypeStartup)
	myApp.AddSchedule(update, app.ScheduleTypeRepeating)
	myApp.AddResource(&logger)
	myApp.AddResource(&ticks{total: 0})
	myApp.AddResource(&timeStarted{
		millis: time.Now().UnixMilli(),
	})
	myApp.AddResource(&lastPrintTime{
		millis: time.Now().UnixMilli(),
	})

	myApp.AddSystem(startup, insertComponents)
	myApp.AddSystem(update, runQuery)
	myApp.AddSystem(update, printTPS)

	run.RunApp(&myApp)
}

// Fills the world with 7 different archetypes
func insertComponents(world *ecs.World) error {
	if _, err := ecs.Spawn(world, &emptyComponentA{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &emptyComponentB{}, &emptyComponentA{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &emptyComponentC{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &emptyComponentA{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &componentWithValue{value: 456}, &emptyComponentA{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &emptyComponentC{}, &emptyComponentA{}, &emptyComponentB{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &emptyComponentD{}, &emptyComponentC{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{}); err != nil {
		return err
	}

	return nil
}

func runQuery(log app.Logger, query ecs.Query2[emptyComponentA, componentWithValue, ecs.With[emptyComponentB]]) {
	total := 0
	query.Result().Iter(func(entityId ecs.EntityId, a *emptyComponentA, b *componentWithValue) error {
		total += b.value
		return nil
	})

	log.Debug(fmt.Sprintf("got: %d", total))
}

func printTPS(counter *ticks, startTime timeStarted, lastPrintTime *lastPrintTime) {
	counter.total++

	now := time.Now().UnixMilli()
	if now-lastPrintTime.millis <= 1_000 {
		return
	}
	lastPrintTime.millis = now

	timeRan := now - startTime.millis
	tps := 1.0 / (float64(timeRan) / float64(counter.total) / 1000.0)
	fmt.Printf("TPS: %d\n", int(tps))
}
