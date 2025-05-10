package ecs

import (
	"testing"
)

func BenchmarkToComponentId(b *testing.B) {
	world := DefaultWorld()
	type aComponent struct{ Component }
	component := aComponent{}

	for b.Loop() {
		ComponentIdOf(component, &world)
	}
}

func BenchmarkGetComponentId(b *testing.B) {
	world := DefaultWorld()
	type aComponent struct{ Component }

	for b.Loop() {
		ComponentIdFor[aComponent](&world)
		ComponentIdFor[*aComponent](&world)
	}
}
