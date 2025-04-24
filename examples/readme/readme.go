// The example in the readme. It is added here to catch compiler errors
package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/lucdrenth/murph_engine/src/ecs"
)

type Position struct {
	X float64
	Y float64

	ecs.Component
}

type Velocity struct {
	X float64
	Y float64

	ecs.Component
}

func main() {
	world := ecs.NewWorld()
	fmt.Printf("Hello %T! \n", world)

	for range 3 {
		// Create a new Entity with a Position and a Velocity component
		ecs.Spawn(&world,
			&Position{X: rand.Float64() * 100, Y: rand.Float64() * 100},
			&Velocity{X: rand.NormFloat64(), Y: rand.NormFloat64()},
		)
	}

	for range 5 {
		// Loop over the entities with the Position and the Velocity component
		query := ecs.Query2[Position, Velocity, ecs.NoFilter, ecs.AllRequired]{}
		query.PrepareOptions()
		query.Exec(&world)

		for position, velocity := range query.Result().Range() {
			position.X += velocity.X
			position.Y += velocity.Y
		}
	}
}
