// Demonstrate how to initialise the ECS
package main

import (
	"fmt"

	"github.com/lucdrenth/murphy/src/ecs"
)

func main() {
	world := ecs.World{}
	fmt.Printf("hello %T\n", world)
}
