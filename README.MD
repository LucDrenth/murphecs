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
	world := ecs.World{}
	fmt.Printf("hello %T\n", world)

    // TODO show the basic ECS functionalities
}
```

For more examples, see [/examples](./examples/).

To run any example, run below command and replace `hello_ecs` with the example you want to run.
```bash
go run ./examples/hello_ecs.go
```
