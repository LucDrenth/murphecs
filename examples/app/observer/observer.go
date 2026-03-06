// Demonstrate how to use observers.
//
// In this demo, we spawn npc that randomly trigger the npc-specific talk observer.
// At some point the global extinction observer will be triggered and all npc entities
// will be despawned.
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
	startup app.Schedule = "Startup"
	update  app.Schedule = "Update"
)

type npc struct {
	ecs.Component

	name string
}

type talk struct {
	ecs.Observer

	text string
}

type extinction struct {
	ecs.Observer
}

func main() {
	var logger app.Logger = &app.SimpleConsoleLogger{}
	myApp, err := app.New(logger, ecs.DefaultWorldConfigs())
	if err != nil {
		panic(err)
	}
	myApp.SetTickRate(time.Second / 10)
	myApp.
		AddSchedule(startup, app.ScheduleOptions{ScheduleType: app.ScheduleTypeStartup}).
		AddSchedule(update, app.ScheduleOptions{ScheduleType: app.ScheduleTypeRepeating}).
		AddResource(&logger)

	myApp.
		AddSystem(startup, spawn).
		AddSystem(startup, observerExtinction).
		AddSystem(update, letNpcTalk).
		AddSystem(update, randomlyTriggerExtinction)

	run.RunApps(myApp)
}

func spawn(world *ecs.World) error {
	for i := range 10 {
		npcEntity, err := ecs.Spawn(world, npc{name: fmt.Sprintf("%d", i+1)})
		if err != nil {
			return fmt.Errorf("failed to spawn npc %d: %w", i+1, err)
		}

		err = ecs.Observe(world, npcEntity, npcTalkObserver)
		if err != nil {
			return fmt.Errorf("failed to add observer for npc %d: %w", i+1, err)
		}
	}

	return nil
}

func npcTalkObserver(world *ecs.World, observer talk) {
	fmt.Println(observer.text)
}

func letNpcTalk(world *ecs.World, query *ecs.Query1[npc, ecs.Default]) error {
	return query.IterUntilErr(func(entityId ecs.EntityId, npc npc) error {
		if rand.IntN(20) == 0 {
			return ecs.TriggerEntity(world, entityId, talk{text: fmt.Sprintf("I am NPC %s", npc.name)})
		}
		return nil
	})
}

func observerExtinction(
	world *ecs.World,
	npcQuery *ecs.Query0[ecs.QueryOptions2[
		ecs.With[npc],
		ecs.Lazy,
	]],
) {
	ecs.On(world, func(world *ecs.World, _ extinction) {
		fmt.Println("! extinction !")

		err := npcQuery.Exec(world)
		if err != nil {
			fmt.Printf("failed to exec npc query: %v\n", err)
			return
		}

		npcQuery.Iter(func(entityId ecs.EntityId) {
			err := ecs.Despawn(world, entityId)
			if err != nil {
				fmt.Printf("failed to despawn npc %d\n", entityId)
			}
		})
	})
}

func randomlyTriggerExtinction(world *ecs.World) {
	if rand.IntN(100) == 0 {
		ecs.Trigger(world, extinction{})
	}
}
