package ecs

import (
	"testing"
)

func BenchmarkToComponentType(b *testing.B) {
	type aComponent struct{ Component }
	component := aComponent{}

	for b.Loop() {
		toComponentType(component)
	}
}

func BenchmarkGetComponentType(b *testing.B) {
	type aComponent struct{ Component }

	for b.Loop() {
		GetComponentType[aComponent]()
	}
}
