// Benchmark ECS using its public (user facing) functions
package ecs_test

import (
	"fmt"
	"testing"

	"github.com/lucdrenth/murphecs/src/ecs"
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
			world := ecs.DefaultWorld()
			ecs.Spawn(&world, &emptyComponentA{})
		}
	})

	b.Run("VariadicTwoComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.DefaultWorld()
			ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{})
		}
	})

	b.Run("VariadicThreeComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.DefaultWorld()
			ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{})
		}
	})

	b.Run("VariadicFourComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.DefaultWorld()
			ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
		}
	})
}

func BenchmarkInsert(b *testing.B) {
	b.Run("VariadicOneComponent", func(b *testing.B) {
		for b.Loop() {
			world := ecs.DefaultWorld()
			entity, _ := ecs.Spawn(&world)

			ecs.Insert(&world, entity, &emptyComponentA{})
		}
	})

	b.Run("VariadicTwoComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.DefaultWorld()
			entity, _ := ecs.Spawn(&world)

			ecs.Insert(&world, entity, &emptyComponentA{}, &emptyComponentB{})
		}
	})

	b.Run("VariadicThreeComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.DefaultWorld()
			entity, _ := ecs.Spawn(&world)

			ecs.Insert(&world, entity, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{})
		}
	})

	b.Run("VariadicFourComponents", func(b *testing.B) {
		for b.Loop() {
			world := ecs.DefaultWorld()
			entity, _ := ecs.Spawn(&world)

			ecs.Insert(&world, entity, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
		}
	})
}

func BenchmarkRemove(b *testing.B) {
	b.Run("OneComponent", func(b *testing.B) {
		world := ecs.DefaultWorld()
		if err := fillWorld(&world); err != nil {
			b.FailNow()
		}

		for b.Loop() {
			entity, _ := ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
			ecs.Remove[emptyComponentA](&world, entity)
		}
	})

	b.Run("TwoComponents", func(b *testing.B) {
		world := ecs.DefaultWorld()
		if err := fillWorld(&world); err != nil {
			b.FailNow()
		}

		for b.Loop() {
			entity, _ := ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
			ecs.Remove2[emptyComponentA, emptyComponentB](&world, entity)
		}
	})

	b.Run("ThreeComponents", func(b *testing.B) {
		world := ecs.DefaultWorld()
		if err := fillWorld(&world); err != nil {
			b.FailNow()
		}

		for b.Loop() {
			entity, _ := ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
			ecs.Remove3[emptyComponentA, emptyComponentB, emptyComponentC](&world, entity)
		}
	})

	b.Run("FourComponents", func(b *testing.B) {
		world := ecs.DefaultWorld()
		if err := fillWorld(&world); err != nil {
			b.FailNow()
		}

		for b.Loop() {
			entity, _ := ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})

			ecs.Remove4[emptyComponentA, emptyComponentB, emptyComponentC, emptyComponentD](&world, entity)
		}
	})
}

func BenchmarkDelete(b *testing.B) {
	world := ecs.DefaultWorld()
	if err := fillWorld(&world); err != nil {
		b.FailNow()
	}

	for b.Loop() {
		entity, _ := ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
		ecs.Delete(&world, entity)
	}
}

func BenchmarkGet(b *testing.B) {
	world := ecs.DefaultWorld()
	if err := fillWorld(&world); err != nil {
		b.FailNow()
	}

	target, _ := ecs.Spawn(&world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})

	if err := fillWorld(&world); err != nil {
		b.FailNow()
	}

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

func BenchmarkHasComponent(b *testing.B) {
	for _, numberOfEntities := range []int{10, 100, 1_000, 10_000} {
		world := ecs.DefaultWorld()

		for range numberOfEntities {
			if err := fillWorld(&world); err != nil {
				b.FailNow()
			}

			if _, err := ecs.Spawn(&world, &emptyComponentA{}); err != nil {
				b.Fatal(err)
			}

			if _, err := ecs.Spawn(&world); err != nil {
				b.Fatal(err)
			}
		}

		entity, err := ecs.Spawn(&world, &emptyComponentA{})
		if err != nil {
			b.Fatal(err)
		}

		b.Run(fmt.Sprintf("ComponentFound-%d-Entities", numberOfEntities), func(b *testing.B) {
			var componentFound bool

			for b.Loop() {
				componentFound, err = ecs.HasComponent[emptyComponentA](&world, entity)
			}

			if err != nil {
				b.Fatal(err)
			}
			if !componentFound {
				b.Fatal("component should have been found")
			}
		})

		b.Run(fmt.Sprintf("ComponentNotFound-%d-Entities", numberOfEntities), func(b *testing.B) {
			var componentFound bool

			for b.Loop() {
				componentFound, err = ecs.HasComponent[emptyComponentB](&world, entity)
			}

			if err != nil {
				b.Fatal(err)
			}
			if componentFound {
				b.Fatal("component should not have been found")
			}
		})
	}
}

func BenchmarkHasComponentId(b *testing.B) {
	for _, numberOfEntities := range []int{10, 100, 1_000, 10_000} {
		world := ecs.DefaultWorld()
		for range numberOfEntities {
			if err := fillWorld(&world); err != nil {
				b.FailNow()
			}

			if _, err := ecs.Spawn(&world, &emptyComponentA{}); err != nil {
				b.Fatal(err)
			}

			if _, err := ecs.Spawn(&world); err != nil {
				b.Fatal(err)
			}
		}

		entity, err := ecs.Spawn(&world, &emptyComponentA{})
		if err != nil {
			b.Fatal(err)
		}

		componentIdA := ecs.ComponentIdFor[emptyComponentA](&world)
		componentIdB := ecs.ComponentIdFor[emptyComponentB](&world)

		b.Run(fmt.Sprintf("ComponentFound-%d-Entities", numberOfEntities), func(b *testing.B) {
			var componentFound bool

			for b.Loop() {
				componentFound, err = ecs.HasComponentId(&world, entity, componentIdA)
			}

			if err != nil {
				b.Fatal(err)
			}
			if !componentFound {
				b.Fatal("component should have been found")
			}
		})

		b.Run(fmt.Sprintf("ComponentNotFound-%d-Entities", numberOfEntities), func(b *testing.B) {
			var componentFound bool

			for b.Loop() {
				componentFound, err = ecs.HasComponentId(&world, entity, componentIdB)
			}

			if err != nil {
				b.Fatal(err)
			}
			if componentFound {
				b.Fatal("component should not have been found")
			}
		})
	}
}

func BenchmarkQuery(b *testing.B) {
	for _, size := range []int{10, 100, 1_000, 10_000} {
		world := ecs.DefaultWorld()

		for range size {
			if err := fillWorld(&world); err != nil {
				b.FailNow()
			}
		}

		b.Run(fmt.Sprintf("Query0-Basic-Size-%d", size), func(b *testing.B) {
			query := ecs.Query0[ecs.Default]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query0-With1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query0[ecs.With[emptyComponentC]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query0-Without1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query0[ecs.Without[emptyComponentC]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query1-Basic-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.Default]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query1-Optional-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.Optional1[emptyComponentA]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query1-ReadOnly-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.AllReadOnly]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query1-With1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.With[emptyComponentC]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query1-Without1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.Without[emptyComponentC]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query2-Basic-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.Default]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query2-Optional-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.Optional1[emptyComponentA]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query2-ReadOnly-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.AllReadOnly]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query2-With1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.With[emptyComponentC]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query2-Without1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.Without[emptyComponentC]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query3-Basic-Size-%d", size), func(b *testing.B) {
			query := ecs.Query3[emptyComponentA, emptyComponentD, emptyComponentC, ecs.Default]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query3-Optional-Size-%d", size), func(b *testing.B) {
			query := ecs.Query3[emptyComponentA, emptyComponentD, emptyComponentC, ecs.Optional1[emptyComponentA]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query3-ReadOnly-Size-%d", size), func(b *testing.B) {
			query := ecs.Query3[emptyComponentA, emptyComponentD, emptyComponentC, ecs.AllReadOnly]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query3-With1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query3[emptyComponentA, emptyComponentD, emptyComponentC, ecs.With[componentWithValue]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query3-Without1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query3[emptyComponentA, emptyComponentD, emptyComponentC, ecs.Without[componentWithValue]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query4-Basic-Size-%d", size), func(b *testing.B) {
			query := ecs.Query4[emptyComponentA, emptyComponentD, emptyComponentB, emptyComponentC, ecs.Default]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query4-Optional-Size-%d", size), func(b *testing.B) {
			query := ecs.Query4[emptyComponentA, emptyComponentD, emptyComponentB, emptyComponentC, ecs.Optional1[emptyComponentA]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query4-ReadOnly-Size-%d", size), func(b *testing.B) {
			query := ecs.Query4[emptyComponentA, emptyComponentD, emptyComponentB, emptyComponentC, ecs.AllReadOnly]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query4-With1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query4[emptyComponentA, emptyComponentD, emptyComponentB, emptyComponentC, ecs.With[componentWithValue]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})

		b.Run(fmt.Sprintf("Query4-Without1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query4[emptyComponentA, emptyComponentD, emptyComponentB, emptyComponentC, ecs.Without[componentWithValue]]{}

			err := query.Prepare(&world)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(&world)
			}
		})
	}
}

func fillWorld(world *ecs.World) error {
	if _, err := ecs.Spawn(world, &emptyComponentA{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &emptyComponentB{}, &emptyComponentA{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &emptyComponentC{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &emptyComponentA{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &componentWithValue{value: 456}, &emptyComponentA{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &emptyComponentC{}, &emptyComponentA{}, &emptyComponentB{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &emptyComponentD{}, &emptyComponentC{}); err != nil {
		return err
	}
	if _, err := ecs.Spawn(world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{}); err != nil {
		return err
	}

	return nil
}
