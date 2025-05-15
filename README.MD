# About
Murph Engine is a Game Engine based on a custom Entity Component System written in Go.

!! This project is currently in a very early stage of development.

# Installation
```bash
go get github.com/lucdrenth/murph_engine
```

# Usage
Below is an example of the ECS. 
```go
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
	world := ecs.DefaultWorld()
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
		query := ecs.Query2[Position, Velocity, ecs.Default]{}
		query.Prepare(&world)
		query.Exec(&world)

		for position, velocity := range query.Result().Range() {
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

# Development
## Benchmarking
To compare benchmarks with changed code, run below command with your new code
```bash
go test -count=10 -bench=. ./benchmark/ecs/ > new.txt
```
and then below command with your old code
```bash
go test -count=10 -bench=. ./benchmark/ecs/ > new.txt
```
Now run [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) on the files created above:
```bash
go tool benchstat old.txt new.txt
```
