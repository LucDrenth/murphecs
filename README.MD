# About
Murphy is an Entity Component System written in Golang.

# Installation
```bash
go get github.com/lucdrenth/murphy
```

# Usage
Below is an example of this ECS. 
```go
package main

import (
	"fmt"

	"github.com/lucdrenth/murphy/src/ecs"
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
