// Benchmark user facing ECS functions
package ecs_test

import (
	"testing"

	"github.com/lucdrenth/murph_engine/src/ecs"
)

type emptyComponentA struct{ ecs.Component }
type emptyComponentB struct{ ecs.Component }
type emptyComponentC struct{ ecs.Component }
type emptyComponentD struct{ ecs.Component }

func BenchmarkSpawn(b *testing.B) {
	b.Run("VariadicOneComponent", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			ecs.Spawn(&world, emptyComponentA{})
		}
	})

	b.Run("VariadicTwoComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			ecs.Spawn(&world, emptyComponentA{}, emptyComponentB{})
		}
	})

	b.Run("VariadicThreeComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			ecs.Spawn(&world, emptyComponentA{}, emptyComponentB{}, emptyComponentC{})
		}
	})

	b.Run("VariadicFourComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			ecs.Spawn(&world, emptyComponentA{}, emptyComponentB{}, emptyComponentC{}, emptyComponentD{})
		}
	})
}

func BenchmarkInsert(b *testing.B) {
	b.Run("VariadicOneComponent", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world)

			ecs.Insert(&world, entity, emptyComponentA{})
		}
	})

	b.Run("VariadicTwoComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world)

			ecs.Insert(&world, entity, emptyComponentA{}, emptyComponentB{})
		}
	})

	b.Run("VariadicThreeComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world)

			ecs.Insert(&world, entity, emptyComponentA{}, emptyComponentB{}, emptyComponentC{})
		}
	})

	b.Run("VariadicFourComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world)

			ecs.Insert(&world, entity, emptyComponentA{}, emptyComponentB{}, emptyComponentC{}, emptyComponentD{})
		}
	})
}

func BenchmarkRemove(b *testing.B) {
	b.Run("OneComponent", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world, emptyComponentA{}, emptyComponentB{}, emptyComponentC{}, emptyComponentD{})

			ecs.Remove[emptyComponentA](&world, entity)
		}
	})

	b.Run("TwoComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world, emptyComponentA{}, emptyComponentB{}, emptyComponentC{}, emptyComponentD{})

			ecs.Remove2[emptyComponentA, emptyComponentB](&world, entity)
		}
	})

	b.Run("ThreeComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world, emptyComponentA{}, emptyComponentB{}, emptyComponentC{}, emptyComponentD{})

			ecs.Remove3[emptyComponentA, emptyComponentB, emptyComponentC](&world, entity)
		}
	})

	b.Run("FourComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world, emptyComponentA{}, emptyComponentB{}, emptyComponentC{}, emptyComponentD{})

			ecs.Remove4[emptyComponentA, emptyComponentB, emptyComponentC, emptyComponentD](&world, entity)
		}
	})
}

func BenchmarkDelete(b *testing.B) {
	world := ecs.NewWorld()

	for b.Loop() {
		entity, _ := ecs.Spawn(&world, emptyComponentA{}, emptyComponentB{}, emptyComponentC{}, emptyComponentD{})
		ecs.Delete(&world, entity)
	}
}

func BenchmarkGet(b *testing.B) {
	world := ecs.NewWorld()
	ecs.Spawn(&world, emptyComponentA{}, emptyComponentB{}, emptyComponentC{}, emptyComponentD{})
	target, _ := ecs.Spawn(&world, emptyComponentA{}, emptyComponentB{}, emptyComponentC{}, emptyComponentD{})
	ecs.Spawn(&world, emptyComponentA{}, emptyComponentB{}, emptyComponentC{}, emptyComponentD{})

	b.Run("Get", func(b *testing.B) {
		for b.Loop() {
			ecs.Get[emptyComponentA](&world, target)
		}
	})

	b.Run("Get2", func(b *testing.B) {
		for b.Loop() {
			ecs.Get2[emptyComponentA, emptyComponentB](&world, target)
		}
	})

	b.Run("Get3", func(b *testing.B) {
		for b.Loop() {
			ecs.Get3[emptyComponentA, emptyComponentB, emptyComponentC](&world, target)
		}
	})

	b.Run("Get4", func(b *testing.B) {
		for b.Loop() {
			ecs.Get4[emptyComponentA, emptyComponentB, emptyComponentC, emptyComponentD](&world, target)
		}
	})
}
