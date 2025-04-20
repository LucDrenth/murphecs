// Demonstrate how to query entities and components
package main

import (
	"fmt"

	"github.com/lucdrenth/murph_engine/src/ecs"
)

type NPC struct {
	name string
	ecs.Component
}
type Friendly struct{ ecs.Component }
type Dialog struct {
	text string
	ecs.Component
}

func main() {
	world := ecs.NewWorld()
	ecs.Spawn(&world, &NPC{name: "Murphy"}, &Friendly{}, &Dialog{text: "hi my name is Murphy"})
	ecs.Spawn(&world, &NPC{name: "Yuki"}, &Dialog{text: "my name is Yuki"})
	ecs.Spawn(&world, &NPC{name: "Fiona"}, &Friendly{}, &Dialog{text: "hi my name is Fiona"})
	ecs.Spawn(&world, &NPC{name: "Bob"})

	// Query all NPC components
	queryResult := ecs.Query1[NPC](&world)
	queryResult.Iter(func(entityId ecs.EntityId, npc *NPC) error {
		fmt.Printf("simple query: %d: %s \n", entityId, npc.name)
		return nil
	})

	// Query all NPC components of entities that have both the Friendly and the Dialog component
	queryResult = ecs.Query1[NPC](&world, ecs.With[Friendly](), ecs.With[Dialog]())
	queryResult.Iter(func(entityId ecs.EntityId, npc *NPC) error {
		fmt.Printf("query with Friendly and Dialog: %d: %s \n", entityId, npc.name)
		return nil
	})

	// Query all NPC and Dialog components of entities that do not have the Friendly component
	queryResult2 := ecs.Query2[NPC, Dialog](&world, ecs.Without[Friendly]())
	queryResult2.Iter(func(entityId ecs.EntityId, npc *NPC, dialog *Dialog) error {
		fmt.Printf("query without Friendly: %d: %s says %s \n", entityId, npc.name, dialog.text)
		return nil
	})

	// Query all NPC and (optionally) Dialog component of all entities that do not have the friendly component
	queryResult2 = ecs.Query2[NPC, Dialog](&world, ecs.Optional[Dialog](), ecs.Without[Friendly]())
	queryResult2.Iter(func(entityId ecs.EntityId, npc *NPC, dialog *Dialog) error {
		fmt.Printf("query with optional Dialog: %d: %s has dialog %v \n", entityId, npc.name, dialog)
		return nil
	})
}
