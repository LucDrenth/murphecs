# About
Murphecs is an Entity Component System written in Go.

# Why Murphecs
There are a bunch of other ECS implementation in Go already, so why did I choose to make Murphecs? 

1. Other ECS I have seen all focus on performance over usability. While performance is important, and is one of the main reasons for choosing to use an ECS, I believe usability comes first. 
2. Other ECS I have seen all focus on the Entity Component part, and leave out the System part. Murphecs has easy to use systems functionality that allows for automatic parallelization.

# Installation
```bash
go get github.com/lucdrenth/murphecs
```

# Usage
Below is an example of the ECS. 
```go
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
		query := ecs.Query2[*Position, Velocity, ecs.Default]{}
		query.Prepare(world, nil)
		query.Exec(world)

		for position, velocity := range query.Range() {
			position.X += velocity.X
			position.Y += velocity.Y
		}
	}
}
```

For more examples, see [/examples](./examples/).

To run any example, run below command and replace `readme` with the directory of the example you want to run.
```bash
go run ./examples/readme
```

# Contributing
If you want to contribute to Murphecs, feel free to open an issue or make a PR. If the change could be considered controversial, consider making an issue first to discuss the changes before working on the implementation. 

## Benchmarking
To compare benchmarks with that of another commit, run `./scripts/compare-benchmarks.sh`. By default it will compare the benchmarks in [/benchmarks](./benchmark) to that of the main branch. You can also pass a custom commit hash or run only 1 specific benchmark function. Run `./scripts/compare-benchmarks.sh -help` for more info.

Be aware that this script will temporarily stash your current changes. If anything goes wrong and you are missing changes, just run `git stash pop`.
