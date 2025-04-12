// Demonstrate how to get an entity its components
package main

import (
	"fmt"

	"github.com/lucdrenth/murphy/src/ecs"
)

type Foo struct {
	value int
	ecs.Component
}
type Bar struct {
	value int
	ecs.Component
}
type Baz struct{ ecs.Component }
type ComponentThatWasNotAdded struct{ ecs.Component }

func main() {
	world := ecs.NewWorld()

	entity, _ := world.Spawn(Foo{value: 25}, Bar{value: 100}, Baz{})
	bar, _ := ecs.Get[Bar](entity, &world)
	fmt.Printf("Value of Bar is %d\n", (*bar).value)

	foo, bar, _ := ecs.Get2[Foo, Bar](entity, &world)
	fmt.Printf("Value of Foo is %d, value of Baz is %d\n", (*foo).value, (*bar).value)

	_, _, _, err := ecs.Get3[Foo, Bar, ComponentThatWasNotAdded](entity, &world)
	fmt.Println(err)
}
