// Demonstrate how to despawn an entity
package main

import (
	"fmt"

	"github.com/lucdrenth/murphecs/src/ecs"
)

type NPC struct{ ecs.Component }

func main() {
	world := ecs.NewDefaultWorld()
	entity, _ := world.Spawn(&NPC{})

	fmt.Printf("Before deleting: %d entity in the world\n", world.CountEntities())
	world.Despawn(entity)
	fmt.Printf("After deleting: %d entities in the world\n", world.CountEntities())
}
