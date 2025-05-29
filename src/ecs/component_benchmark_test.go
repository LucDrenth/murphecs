package ecs

import (
	"testing"
)

func BenchmarkToComponentId(b *testing.B) {
	type componentA struct{ Component }

	world := NewDefaultWorld()
	component := componentA{}

	for b.Loop() {
		ComponentIdOf(component, world)
	}
}

func BenchmarkGetComponentId(b *testing.B) {
	type componentA struct{ Component }

	world := NewDefaultWorld()

	for b.Loop() {
		ComponentIdFor[componentA](world)
		ComponentIdFor[*componentA](world)
	}
}
