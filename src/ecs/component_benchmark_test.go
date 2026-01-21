package ecs

import (
	"testing"
)

func BenchmarkToComponentId(b *testing.B) {
	type componentA struct{ Component }

	world := NewDefaultWorld()
	component := componentA{}

	b.Run("component by reference", func(b *testing.B) {
		for b.Loop() {
			ComponentIdOf(&component, world)
		}
	})

	b.Run("component by value", func(b *testing.B) {
		for b.Loop() {
			ComponentIdOf(component, world)
		}
	})
}

func BenchmarkGetComponentId(b *testing.B) {
	type componentA struct{ Component }

	world := NewDefaultWorld()

	b.Run("pointer", func(b *testing.B) {
		for b.Loop() {
			ComponentIdFor[*componentA](world)
		}
	})

	b.Run("raw", func(b *testing.B) {
		for b.Loop() {
			ComponentIdFor[componentA](world)
		}
	})
}
