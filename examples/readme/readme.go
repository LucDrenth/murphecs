// The example in the readme. It is added here to catch compiler errors
package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/lucdrenth/murphecs/src/ecs"
)

type Position struct {
	ecs.Component

	X, Y float64
}

type Velocity struct {
	ecs.Component

	X, Y float64
}

func main() {
	world := ecs.NewDefaultWorld()
	fmt.Printf("Hello %T! \n", world)

	for range 3 {
		// Create a new Entity with a Position and a Velocity component
		ecs.Spawn(world,
			Position{X: rand.Float64() * 100, Y: rand.Float64() * 100},
			Velocity{X: rand.NormFloat64(), Y: rand.NormFloat64()},
		)
	}

	for range 5 {
		// Loop over the entities with the Position and the Velocity component
		query := ecs.Query2[Position, Velocity, ecs.Default]{}
		query.Prepare(world, nil)
		query.Exec(world)

		for position, velocity := range query.Range() {
			position.X += velocity.X
			position.Y += velocity.Y
		}
	}
}
