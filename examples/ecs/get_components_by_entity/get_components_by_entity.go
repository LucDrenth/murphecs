// Demonstrate how to get an entity its components
package main

import (
	"fmt"

	"github.com/lucdrenth/murphecs/src/ecs"
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
	world := ecs.NewDefaultWorld()

	entity, _ := world.Spawn(Friendly{}, Health{max: 100, current: 80}, NPC{name: "Murphy"})

	// only get the NPC component
	npc, _ := world.Get1[NPC](entity)
	fmt.Printf("npc name is %s\n", npc.name)

	// use component pointer type to get a mutable reference
	npcPointer, _ := world.Get1[*NPC](entity)
	npcPointer.name = "Yuki"
	npc, _ = world.Get1[NPC](entity)
	fmt.Printf("npc name is %s\n", npc.name)

	// get both the NPC and the Health component
	npc, health, _ := world.Get2[NPC, Health](entity)
	fmt.Printf("npc name is %s, current health is %d\n", npc.name, health.current)

	// returns an error because the entity does not have the Aggressive component
	_, _, _, err := world.Get3[NPC, Health, Aggressive](entity)
	fmt.Println(err)
}
