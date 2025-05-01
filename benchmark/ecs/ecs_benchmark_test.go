// Benchmark ECS using its public (user facing) functions
package ecs_test

import (
	"fmt"
	"testing"

	"github.com/lucdrenth/murph_engine/src/ecs"
)

type emptyComponentA struct{ ecs.Component }
type emptyComponentB struct{ ecs.Component }
type emptyComponentC struct{ ecs.Component }
type emptyComponentD struct{ ecs.Component }
type componentWithValue struct {
	ecs.Component
	value int
}

func BenchmarkSpawn(b *testing.B) {
	b.Run("VariadicOneComponent", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			ecs.Spawn(&world, &emptyComponentA{})
		}
	})

	b.Run("VariadicTwoComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{})
		}
	})

	b.Run("VariadicThreeComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{})
		}
	})

	b.Run("VariadicFourComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
		}
	})
}

func BenchmarkInsert(b *testing.B) {
	b.Run("VariadicOneComponent", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world)

			ecs.Insert(&world, entity, &emptyComponentA{})
		}
	})

	b.Run("VariadicTwoComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world)

			ecs.Insert(&world, entity, &emptyComponentA{}, &emptyComponentB{})
		}
	})

	b.Run("VariadicThreeComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world)

			ecs.Insert(&world, entity, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{})
		}
	})

	b.Run("VariadicFourComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world)

			ecs.Insert(&world, entity, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
		}
	})
}

func BenchmarkRemove(b *testing.B) {
	b.Run("OneComponent", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})

			ecs.Remove[emptyComponentA](&world, entity)
		}
	})

	b.Run("TwoComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})

			ecs.Remove2[emptyComponentA, emptyComponentB](&world, entity)
		}
	})

	b.Run("ThreeComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})

			ecs.Remove3[emptyComponentA, emptyComponentB, emptyComponentC](&world, entity)
		}
	})

	b.Run("FourComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewWorld()
			entity, _ := ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})

			ecs.Remove4[emptyComponentA, emptyComponentB, emptyComponentC, emptyComponentD](&world, entity)
		}
	})
}

func BenchmarkDelete(b *testing.B) {
	world := ecs.NewWorld()

	for b.Loop() {
		entity, _ := ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
		ecs.Delete(&world, entity)
	}
}

func BenchmarkGet(b *testing.B) {
	world := ecs.NewWorld()
	ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
	target, _ := ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
	ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})

	b.Run("Get1", func(b *testing.B) {
		for b.Loop() {
			ecs.Get1[emptyComponentA](&world, target)
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

func BenchmarkQuery(b *testing.B) {
	for _, size := range []int{10, 100, 1_000, 10_000, 100_000} {
		world := ecs.NewWorld()

		for range size {
			if _, err := ecs.Spawn(&world, &emptyComponentA{}); err != nil {
				b.FailNow()
			}
			if _, err := ecs.Spawn(&world, &emptyComponentB{}, &emptyComponentA{}); err != nil {
				b.FailNow()
			}
			if _, err := ecs.Spawn(&world, &emptyComponentC{}); err != nil {
				b.FailNow()
			}
			if _, err := ecs.Spawn(&world, &emptyComponentA{}); err != nil {
				b.FailNow()
			}
			if _, err := ecs.Spawn(&world, &componentWithValue{value: size}, &emptyComponentA{}); err != nil {
				b.FailNow()
			}
			if _, err := ecs.Spawn(&world, &emptyComponentC{}, &emptyComponentA{}, &emptyComponentB{}); err != nil {
				b.FailNow()
			}
			if _, err := ecs.Spawn(&world, &emptyComponentD{}, &emptyComponentC{}); err != nil {
				b.FailNow()
			}
			if _, err := ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{}); err != nil {
				b.FailNow()
			}
		}

		b.Run(fmt.Sprintf("Query1-Basic-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.NoFilter, ecs.NoOptional, ecs.NoReadOnly]{}

			err := query.PrepareOptions()
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query1-Optional-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.NoFilter, ecs.Optional1[emptyComponentA], ecs.NoReadOnly]{}

			err := query.PrepareOptions()
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query1-ReadOnly-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.NoFilter, ecs.NoOptional, ecs.AllReadOnly]{}

			err := query.PrepareOptions()
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query1-With1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.With[emptyComponentC], ecs.NoOptional, ecs.AllReadOnly]{}

			err := query.PrepareOptions()
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query1-Without1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.Without[emptyComponentC], ecs.NoOptional, ecs.AllReadOnly]{}

			err := query.PrepareOptions()
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query2-Basic-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.NoFilter, ecs.NoOptional, ecs.NoReadOnly]{}

			err := query.PrepareOptions()
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query2-Optional-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.NoFilter, ecs.Optional1[emptyComponentA], ecs.NoReadOnly]{}

			err := query.PrepareOptions()
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query2-ReadOnly-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.NoFilter, ecs.NoOptional, ecs.AllReadOnly]{}

			err := query.PrepareOptions()
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query2-With1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.With[emptyComponentC], ecs.NoOptional, ecs.AllReadOnly]{}

			err := query.PrepareOptions()
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query2-Without1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.Without[emptyComponentC], ecs.NoOptional, ecs.AllReadOnly]{}

			err := query.PrepareOptions()
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})
	}
}
