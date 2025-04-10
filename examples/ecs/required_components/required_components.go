// Demonstrate required components.
//
// When spawning an entity, each components will be checked for its RequiredComponents(),
// and they will be exhaustively added if they are not already given.

package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/lucdrenth/murphy/src/ecs"
)

type A struct{ ecs.Component }
type B struct{ ecs.Component }
type C struct{ ecs.Component }

func (a A) requiredComponents() []ecs.IComponent {
	return []ecs.IComponent{
		B{}, C{},
	}
}

func main() {
	a := A{}

	requiredComponents := []string{}
	for _, c := range a.requiredComponents() {
		requiredComponents = append(requiredComponents, reflect.TypeOf(c).String())
	}

	fmt.Printf("Component %T requires: %s\n", a, strings.Join(requiredComponents, ", "))
}
