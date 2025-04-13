# About
Murph Engine is a Game Engine based on a custom Entity Component System written in Go.

!! This project is currently in a very early stage of development.

# Installation
```bash
go get github.com/lucdrenth/murph-engine
```

# Usage
Below is an example of the ECS. 
```go
package main

import (
	"fmt"

	"github.com/lucdrenth/murph/engine/src/ecs"
)

func main() {
	world := ecs.NewWorld()
	fmt.Printf("hello %T\n", world)

    // TODO show the basic ECS functionalities
}
```

For more examples, see [/examples](./examples/).

To run any example, run below command and replace `readme` with the directory of the example you want to run.
```bash
go run ./examples/readme
```
