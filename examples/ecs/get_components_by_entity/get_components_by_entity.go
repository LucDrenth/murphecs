// Demonstrate how to get an entity its components
package main

import (
	"fmt"

	"github.com/lucdrenth/murphy/src/ecs"
)

type Health struct {
	max     int
	current int
	ecs.Component
}

type Friendly struct{ ecs.Component }
type Aggressive struct{ ecs.Component }

type NPC struct {
	name string
	ecs.Component
}

func main() {
	world := ecs.NewWorld()

	entity, _ := ecs.Spawn(&world, Friendly{}, Health{max: 100, current: 80}, NPC{name: "Murphy"})

	// only get the NPC component
	npc, _ := ecs.Get[NPC](entity, &world)
	fmt.Printf("npc name is %s\n", (*npc).name)

	// get bot hthe NPC and the Health component
	npc, health, _ := ecs.Get2[NPC, Health](entity, &world)
	fmt.Printf("npc name is %s, current health is %d\n", (*npc).name, (*health).current)

	// returns an error because the entity does not have the Aggressive component
	_, _, _, err := ecs.Get3[NPC, Health, Aggressive](entity, &world)
	fmt.Println(err)
}
