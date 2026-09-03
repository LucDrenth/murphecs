// Demonstrate how to remove a component from an entity
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

func main() {
	world := ecs.NewDefaultWorld()

	// Spawn an entity with two components: NPC and Dialog
	entity, _ := world.Spawn(NPC{name: "Murphy"}, Dialog{text: "Hello!"})

	// Get the dialog to see its current value
	dialog, err := world.Get1[Dialog](entity)
	fmt.Printf("Before remove: dialog=%v, err=%v\n", dialog.text, err)

	// Remove Dialog component so that entity only has 1 component left: NPC
	world.Remove1[Dialog](entity)

	// Getting the removed component will now fail
	dialog, err = world.Get1[Dialog](entity)
	fmt.Printf("After remove: dialog=%v, err=%v\n", dialog, err)

	// Getting the component that was not removed still works
	npc, err := world.Get1[NPC](entity)
	fmt.Printf("After remove: npc=%v, err=%v\n", npc.name, err)

	// Removing the Dialog component after it was already removed will result in an error
	err = world.Remove1[Dialog](entity)
	fmt.Printf("Error when removing a component that is not present: %v\n", err)
}
