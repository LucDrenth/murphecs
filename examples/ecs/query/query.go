// Demonstrate how to query entities and components
package main

import (
	"fmt"

	"github.com/lucdrenth/murphecs/src/ecs"
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
	world := ecs.NewDefaultWorld()
	ecs.Spawn(world, NPC{name: "Murphy"}, Friendly{}, Dialog{text: "hi my name is Murphy"})
	ecs.Spawn(world, NPC{name: "Yuki"}, Dialog{text: "my name is Yuki"})
	ecs.Spawn(world, NPC{name: "Fiona"}, Friendly{}, Dialog{text: "hi my name is Fiona"})
	ecs.Spawn(world, NPC{name: "Bob"})

	// Query all NPC components
	query := ecs.Query1[NPC, ecs.Default]{}
	query.Prepare(world, nil)
	query.Exec(world)
	query.Iter(func(entityId ecs.EntityId, npc NPC) {
		fmt.Printf("simple query: %d: %s \n", entityId, npc.name)
	})

	// Query all NPC components of entities that have the Friendly component
	query2 := ecs.Query1[NPC, ecs.With[Friendly]]{}
	query2.Prepare(world, nil)
	query2.Exec(world)
	query2.Iter(func(entityId ecs.EntityId, npc NPC) {
		fmt.Printf("query with Friendly: %d: %s \n", entityId, npc.name)
	})

	// Query all NPC and Dialog components of entities that do not have the Friendly component
	query3 := ecs.Query2[NPC, Dialog, ecs.Without[Friendly]]{}
	query3.Prepare(world, nil)
	query3.Exec(world)
	query3.Iter(func(entityId ecs.EntityId, npc NPC, dialog Dialog) {
		fmt.Printf("query without Friendly: %d: %s says %s \n", entityId, npc.name, dialog.text)
	})

	// Specify component pointer to mutate the component
	query4 := ecs.Query2[NPC, *Dialog, ecs.With[Friendly]]{}
	query4.Prepare(world, nil)
	query4.Exec(world)
	query4.Iter(func(entityId ecs.EntityId, npc NPC, dialog *Dialog) {
		dialog.text = "I am super friendly!"
		fmt.Printf("query with pointer: %d: %s says %s \n", entityId, npc.name, dialog.text)
	})

	// You can specify multiple query options with ecs.QueryOptions2, ecs.QueryOptions3 ...
	_ = ecs.Query1[
		NPC,
		ecs.QueryOptions2[
			ecs.With[Dialog],
			ecs.Lazy,
		],
	]{}

	// You can specify all possible query options with ecs.QueryOptions
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
			ecs.NoOptional,
			ecs.Lazy,
			ecs.DefaultWorld,
		],
	]{}

	// You can set query options using functions instead of with generics. You won't have to call Prepare
	dynamicallyBuildQuery := ecs.Query1[NPC, ecs.Default]{}
	ecs.QueryWith[Dialog](world, &dynamicallyBuildQuery)
	ecs.QueryWithout[Friendly](world, &dynamicallyBuildQuery)
	ecs.QueryWithOptional[NPC](world, &dynamicallyBuildQuery)
	dynamicallyBuildQuery.Exec(world)
}
