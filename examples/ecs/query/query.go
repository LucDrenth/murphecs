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
	query := ecs.Query1[NPC, ecs.Default]{}
	query.Prepare()
	query.Exec(&world)
	query.Result().Iter(func(entityId ecs.EntityId, npc *NPC) error {
		fmt.Printf("simple query: %d: %s \n", entityId, npc.name)
		return nil
	})

	// Query all NPC components of entities that have the Friendly component
	query2 := ecs.Query1[NPC, ecs.With[Friendly]]{}
	query2.Prepare()
	query2.Exec(&world)
	query2.Result().Iter(func(entityId ecs.EntityId, npc *NPC) error {
		fmt.Printf("query with Friendly: %d: %s \n", entityId, npc.name)
		return nil
	})

	// Query all NPC and Dialog components of entities that do not have the Friendly component
	query3 := ecs.Query2[NPC, Dialog, ecs.Without[Friendly]]{}
	query3.Prepare()
	query3.Exec(&world)
	query3.Result().Iter(func(entityId ecs.EntityId, npc *NPC, dialog *Dialog) error {
		fmt.Printf("query without Friendly: %d: %s says %s \n", entityId, npc.name, dialog.text)
		return nil
	})

	// You can give multiple options with QueryOptions2, QueryOptions3 etc.
	_ = ecs.Query1[
		NPC,
		ecs.QueryOptions2[
			ecs.With[Dialog],
			ecs.AllReadOnly,
		],
	]{}

	// You can specify all possible query options with QueryOptions:
	_ = ecs.Query1[
		NPC,
		ecs.QueryOptions[
			ecs.Or[
				ecs.Without[Dialog],
				ecs.And[
					ecs.With[Dialog],
					ecs.With[Friendly],
				],
			],
			ecs.AllReadOnly,
		],
	]{}
}
