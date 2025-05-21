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
	ecs.Spawn(&world, &NPC{name: "Murphy"}, &Friendly{}, &Dialog{text: "hi my name is Murphy"})
	ecs.Spawn(&world, &NPC{name: "Yuki"}, &Dialog{text: "my name is Yuki"})
	ecs.Spawn(&world, &NPC{name: "Fiona"}, &Friendly{}, &Dialog{text: "hi my name is Fiona"})
	ecs.Spawn(&world, &NPC{name: "Bob"})

	// Query all NPC components
	query := ecs.Query1[NPC, ecs.Default]{}
	query.Prepare(&world)
	query.Exec(&world)
	query.Result().Iter(func(entityId ecs.EntityId, npc *NPC) error {
		fmt.Printf("simple query: %d: %s \n", entityId, npc.name)
		return nil
	})

	// Query all NPC components of entities that have the Friendly component
	query2 := ecs.Query1[NPC, ecs.With[Friendly]]{}
	query2.Prepare(&world)
	query2.Exec(&world)
	query2.Result().Iter(func(entityId ecs.EntityId, npc *NPC) error {
		fmt.Printf("query with Friendly: %d: %s \n", entityId, npc.name)
		return nil
	})

	// Query all NPC and Dialog components of entities that do not have the Friendly component
	query3 := ecs.Query2[NPC, Dialog, ecs.Without[Friendly]]{}
	query3.Prepare(&world)
	query3.Exec(&world)
	query3.Result().Iter(func(entityId ecs.EntityId, npc *NPC, dialog *Dialog) error {
		fmt.Printf("query without Friendly: %d: %s says %s \n", entityId, npc.name, dialog.text)
		return nil
	})

	// You can specify multiple query options with ecs.QueryOptions2, ecs.QueryOptions3 ...
	_ = ecs.Query1[
		NPC,
		ecs.QueryOptions2[
			ecs.With[Dialog],
			ecs.AllReadOnly,
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
			ecs.AllReadOnly,
			ecs.Lazy,
			ecs.DefaultWorld,
		],
	]{}

	// You can set query options using functions instead of with generics. You won't have to call Prepare
	dynamicallyBuildQuery := ecs.Query1[NPC, ecs.Default]{}
	ecs.QueryWith[Dialog](&world, &dynamicallyBuildQuery)
	ecs.QueryWithout[Friendly](&world, &dynamicallyBuildQuery)
	ecs.QueryWithAllReadOnly(&dynamicallyBuildQuery)
	ecs.QueryWithOptional[NPC](&world, &dynamicallyBuildQuery)
	dynamicallyBuildQuery.Exec(&world)
}
