package main

import (
	"fmt"

	"github.com/lucdrenth/murph/engine/src/ecs"
)

type NPC struct {
	name string
	ecs.Component
}
type Dialog struct {
	text string
	ecs.Component
}
type Friendly struct {
	ecs.Component
}

func main() {
	world := ecs.NewWorld()

	entity, _ := ecs.Spawn(&world, NPC{name: "Murphy"})

	// Insert Dialog component in to the entity, so that is has both NPC and Dialog
	ecs.Insert(&world, entity, Dialog{text: "good morning"})

	// Insert Dialog component, that already exists, and a new Friendly component.
	// This will return an error about Dialog already being present, so it will be skipped.
	// The Friendly component, that is not already present, will still be added.
	err := ecs.Insert(&world, entity, Dialog{text: "good evening"}, Friendly{})
	fmt.Printf("Insert error: %v\n", err)

	dialog, _ := ecs.Get[Dialog](&world, entity)
	fmt.Printf("Dialog text: %s\n", dialog.text)
}
