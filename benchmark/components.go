package benchmark

import "github.com/lucdrenth/murphecs/src/ecs"

type emptyComponentA struct{ ecs.Component }
type emptyComponentB struct{ ecs.Component }
type emptyComponentC struct{ ecs.Component }
type emptyComponentD struct{ ecs.Component }
type componentWithValue struct {
	ecs.Component
	value int
}
