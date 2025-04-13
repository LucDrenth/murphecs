// Demonstrate required components.
//
// When spawning an entity, each components will be checked for its RequiredComponents(),
// and they will be exhaustively added if they are not already given.

package main

import (
	"fmt"

	"github.com/lucdrenth/murph/engine/src/ecs"
)

type NPC struct{ ecs.Component }

type Health struct {
	max     int
	current int
	ecs.Component
}

type Dialog struct {
	text string
	ecs.Component
}

// TODO this does not overwrite the implementation from ecs.Component
func (NPC) RequiredComponents() []ecs.IComponent {
	return []ecs.IComponent{
		Health{max: 100, current: 50},
		Dialog{text: "I am an NPC!"},
	}
}

func main() {
	world := ecs.NewWorld()

	// because NPC requires Health and Dialog, they will also be added to the entity.
	entity, _ := ecs.Spawn(&world, NPC{})
	dialog, health, _ := ecs.Get2[Dialog, Health](&world, entity)
	fmt.Printf("npc has %d/%d health and the following dialog: %s\n", (*health).current, (*health).max, (*dialog).text)

	// NPC requires Dialog and provides a default for that component.
	// But because we specify Dialog here, it will not use the default implementation from the required component.
	entity, _ = ecs.Spawn(&world, NPC{}, Dialog{text: "Good morning."})
	dialog, health, _ = ecs.Get2[Dialog, Health](&world, entity)
	fmt.Printf("npc has %d/%d health and the following dialog: %s\n", (*health).current, (*health).max, (*dialog).text)
}
