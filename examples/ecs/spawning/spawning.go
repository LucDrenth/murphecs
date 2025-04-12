package main

import (
	"fmt"

	"github.com/lucdrenth/murphy/src/ecs"
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
	fmt.Printf("hello %T\n", world)

	// You can spawn any amount of unique components!
	entity, err := world.Spawn(
		Friendly{},
		NPC{name: "Murph"},
		Health{max: 100, current: 80},
	)
	fmt.Printf("entity=%d, err=%v\n", entity, err)
}
