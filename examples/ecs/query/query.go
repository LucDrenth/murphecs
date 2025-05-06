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

	queryWithGenerics(&world)
	fmt.Println()
	queryWithDirectQuery(&world)
}

// We can define queries as generics. This has both advantages and disadvantages over defining queries
// with functions.
//
// Advantages:
//   - When used in system params, they are executed automatically.
//   - When used in system params, they are able to be executed in parallel with other systems
//
// Disadvantages:
//   - When using multiple query options, you need to specify all of them, which is a bit more work.
func queryWithGenerics(world *ecs.World) {
	// Query all NPC components
	query := ecs.Query1[NPC, ecs.Default]{}
	query.Prepare()
	query.Exec(world)
	query.Result().Iter(func(entityId ecs.EntityId, npc *NPC) error {
		fmt.Printf("simple query: %d: %s \n", entityId, npc.name)
		return nil
	})

	// Query all NPC components of entities that have the Friendly component
	query2 := ecs.Query1[NPC, ecs.With[Friendly]]{}
	query2.Prepare()
	query2.Exec(world)
	query2.Result().Iter(func(entityId ecs.EntityId, npc *NPC) error {
		fmt.Printf("query with Friendly: %d: %s \n", entityId, npc.name)
		return nil
	})

	// Query all NPC and Dialog components of entities that do not have the Friendly component
	query3 := ecs.Query2[NPC, Dialog, ecs.Without[Friendly]]{}
	query3.Prepare()
	query3.Exec(world)
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
		],
	]{}
}

// We can compose queries using functions. This has both advantages and disadvantages over using a
// generics query.
//
// Advantages:
//   - Easier to write than generic queries because you don't have to specify each option
//   - In contrast to queries in system params, these queries are not automatically executed every
//     system call.
//
// Disadvantages:
//   - It is not possible to run a system that creates such a query in parallel because it needs to be
//     executed manually, which Run a system that creates such a query in parallel because the system
//     would need the *ecs.World.
func queryWithDirectQuery(world *ecs.World) {
	// Query all NPC components
	query := ecs.NewQuery1[NPC]()
	result := query.Exec(world)
	result.Iter(func(entityId ecs.EntityId, npc *NPC) error {
		fmt.Printf("simple query: %d: %s \n", entityId, npc.name)
		return nil
	})

	// Query all NPC components of entities that have the Friendly component
	query2 := ecs.NewQuery1[NPC]()
	ecs.QueryWith[Friendly](&query2)
	result2 := query2.Exec(world)
	result2.Iter(func(entityId ecs.EntityId, npc *NPC) error {
		fmt.Printf("query with Friendly: %d: %s \n", entityId, npc.name)
		return nil
	})

	// Query all NPC and Dialog components of entities that do not have the Friendly component
	query3 := ecs.NewQuery2[NPC, Dialog]()
	ecs.QueryWithout[Friendly](&query3)
	result3 := query3.Exec(world)
	result3.Iter(func(entityId ecs.EntityId, npc *NPC, dialog *Dialog) error {
		fmt.Printf("query without Friendly: %d: %s says %s \n", entityId, npc.name, dialog.text)
		return nil
	})
}
