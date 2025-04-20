// Demonstrate how to spawn an entity with components
package main

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/ecs"
)

type Health struct {
	max     int
	current int
	ecs.Component
}

type Friendly struct {
	ecs.Component
}

type NPC struct {
	name string
	ecs.Component
}

func main() {
	world := ecs.NewWorld()

	// You can spawn any amount of unique components!
	entity, err := ecs.Spawn(&world,
		&Friendly{},
		&NPC{name: "Murphy"},
		&Health{max: 100, current: 80},
	)
	fmt.Printf("Spawned entity=%d, err=%v\n", entity, err)
}
