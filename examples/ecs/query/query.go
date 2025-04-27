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
	query := ecs.Query1[NPC, ecs.NoFilter, ecs.NoOptional, ecs.NoReadOnly]{}
	query.PrepareOptions()
	query.Exec(&world)
	query.Result().Iter(func(entityId ecs.EntityId, npc *NPC) error {
		fmt.Printf("simple query: %d: %s \n", entityId, npc.name)
		return nil
	})

	// Query all NPC components of entities that hav the Friendly component
	query2 := ecs.Query1[NPC, ecs.With[Friendly], ecs.NoOptional, ecs.NoReadOnly]{}
	query2.PrepareOptions()
	query2.Exec(&world)
	query2.Result().Iter(func(entityId ecs.EntityId, npc *NPC) error {
		fmt.Printf("query with Friendly: %d: %s \n", entityId, npc.name)
		return nil
	})

	// Query all NPC and Dialog components of entities that do not have the Friendly component
	query3 := ecs.Query2[NPC, Dialog, ecs.Without[Friendly], ecs.NoOptional, ecs.NoReadOnly]{}
	query3.PrepareOptions()
	query3.Exec(&world)
	query3.Result().Iter(func(entityId ecs.EntityId, npc *NPC, dialog *Dialog) error {
		fmt.Printf("query without Friendly: %d: %s says %s \n", entityId, npc.name, dialog.text)
		return nil
	})
}
