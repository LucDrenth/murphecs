package main

import (
	"fmt"

	"github.com/lucdrenth/murphecs/src/ecs"
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
	world := ecs.NewDefaultWorld()

	entity, _ := world.Spawn(NPC{name: "Murphy"})

	// Insert Dialog component in to the entity, so that is has both NPC and Dialog
	world.Insert(entity, Dialog{text: "good morning"})

	// Insert Dialog component, that already exists, and a new Friendly component.
	// This will return an error about Dialog already being present, so it will be skipped.
	// The Friendly component, that is not already present, will still be added.
	err := world.Insert(entity, Dialog{text: "good evening"}, Friendly{})
	fmt.Printf("Insert error: %v\n", err)

	dialog, _ := world.Get1[Dialog](entity)
	fmt.Printf("Dialog text: %s\n", dialog.text)
}
