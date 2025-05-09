package ecs

import (
	"testing"
)

func BenchmarkToComponentId(b *testing.B) {
	type aComponent struct{ Component }
	component := aComponent{}

	for b.Loop() {
		ComponentIdOf(component)
	}
}

func BenchmarkGetComponentId(b *testing.B) {
	type aComponent struct{ Component }

	for b.Loop() {
		ComponentIdFor[aComponent]()
		ComponentIdFor[*aComponent]()
	}
}
