package benchmark

import (
	"fmt"
	"testing"

	"github.com/lucdrenth/murphecs/src/ecs"
)

func BenchmarkSpawn(b *testing.B) {
	b.Run("VariadicOneComponent-ByReference", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewDefaultWorld()
			ecs.Spawn(world, &emptyComponentA{})
		}
	})

	b.Run("VariadicOneComponent-ByValue", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewDefaultWorld()
			ecs.Spawn(world, emptyComponentA{})
		}
	})

	b.Run("VariadicTwoComponents-ByReference", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewDefaultWorld()
			ecs.Spawn(world, &emptyComponentA{}, &emptyComponentB{})
		}
	})

	b.Run("VariadicTwoComponents-ByValue", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewDefaultWorld()
			ecs.Spawn(world, emptyComponentA{}, emptyComponentB{})
		}
	})

	b.Run("VariadicThreeComponents-ByReference", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewDefaultWorld()
			ecs.Spawn(world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{})
		}
	})

	b.Run("VariadicThreeComponents-ByValue", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewDefaultWorld()
			ecs.Spawn(world, emptyComponentA{}, emptyComponentB{}, emptyComponentC{})
		}
	})

	b.Run("VariadicFourComponents-ByReference", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewDefaultWorld()
			ecs.Spawn(world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
		}
	})

	b.Run("VariadicFourComponents-ByValue", func(b *testing.B) {
		for b.Loop() {
			world := ecs.NewDefaultWorld()
			ecs.Spawn(world, emptyComponentA{}, emptyComponentB{}, emptyComponentC{}, emptyComponentD{})
		}
	})
}

func BenchmarkInsert(b *testing.B) {
	for _, size := range []int{10, 100, 1_000, 10_000} {
		setupWorld := func() *ecs.World {
			world := ecs.NewDefaultWorld()

			for range size {
				if err := fillWorld(world); err != nil {
					b.FailNow()
				}
			}

			return world
		}

		b.Run(fmt.Sprintf("OneComponent-ByReference-Size-%d", size), func(b *testing.B) {
			world := setupWorld()

			for b.Loop() {
				entity, _ := ecs.Spawn(world)
				ecs.Insert(world, entity, &emptyComponentA{})
			}
		})

		b.Run(fmt.Sprintf("OneComponent-ByValue-Size-%d", size), func(b *testing.B) {
			world := setupWorld()

			for b.Loop() {
				entity, _ := ecs.Spawn(world)
				ecs.Insert(world, entity, emptyComponentA{})
			}
		})

		b.Run(fmt.Sprintf("TwoComponent-ByReference-Size-%d", size), func(b *testing.B) {
			world := setupWorld()

			for b.Loop() {
				entity, _ := ecs.Spawn(world)
				ecs.Insert(world, entity, &emptyComponentA{}, &emptyComponentB{})
			}
		})

		b.Run(fmt.Sprintf("TwoComponent-ByValue-Size-%d", size), func(b *testing.B) {
			world := setupWorld()

			for b.Loop() {
				entity, _ := ecs.Spawn(world)
				ecs.Insert(world, entity, emptyComponentA{}, emptyComponentB{})
			}
		})

		b.Run(fmt.Sprintf("ThreeComponent-ByReference-Size-%d", size), func(b *testing.B) {
			world := setupWorld()

			for b.Loop() {
				entity, _ := ecs.Spawn(world)
				ecs.Insert(world, entity, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{})
			}
		})

		b.Run(fmt.Sprintf("ThreeComponent-ByValue-Size-%d", size), func(b *testing.B) {
			world := setupWorld()

			for b.Loop() {
				entity, _ := ecs.Spawn(world)
				ecs.Insert(world, entity, emptyComponentA{}, emptyComponentB{}, emptyComponentC{})
			}
		})

		b.Run(fmt.Sprintf("FourComponent-ByReference-Size-%d", size), func(b *testing.B) {
			world := setupWorld()

			for b.Loop() {
				entity, _ := ecs.Spawn(world)
				ecs.Insert(world, entity, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
			}
		})

		b.Run(fmt.Sprintf("FourComponent-ByValue-Size-%d", size), func(b *testing.B) {
			world := setupWorld()

			for b.Loop() {
				entity, _ := ecs.Spawn(world)
				ecs.Insert(world, entity, emptyComponentA{}, emptyComponentB{}, emptyComponentC{}, emptyComponentD{})
			}
		})
	}
}

func BenchmarkRemove(b *testing.B) {
	for _, size := range []int{10, 100, 1_000, 10_000} {
		setupWorld := func() *ecs.World {
			world := ecs.NewDefaultWorld()

			for range size {
				if err := fillWorld(world); err != nil {
					b.FailNow()
				}
			}

			return world
		}

		b.Run(fmt.Sprintf("OneComponent-Size-%d", size), func(b *testing.B) {
			world := setupWorld()

			for b.Loop() {
				entity, _ := ecs.Spawn(world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
				ecs.Remove1[emptyComponentA](world, entity)
			}
		})

		b.Run(fmt.Sprintf("TwoComponents-Size-%d", size), func(b *testing.B) {
			world := setupWorld()

			for b.Loop() {
				entity, _ := ecs.Spawn(world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
				ecs.Remove2[emptyComponentA, emptyComponentB](world, entity)
			}
		})

		b.Run(fmt.Sprintf("ThreeComponents-Size-%d", size), func(b *testing.B) {
			world := setupWorld()

			for b.Loop() {
				entity, _ := ecs.Spawn(world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
				ecs.Remove3[emptyComponentA, emptyComponentB, emptyComponentC](world, entity)
			}
		})

		b.Run(fmt.Sprintf("FourComponents-Size-%d", size), func(b *testing.B) {
			world := setupWorld()

			for b.Loop() {
				entity, _ := ecs.Spawn(world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
				ecs.Remove4[emptyComponentA, emptyComponentB, emptyComponentC, emptyComponentD](world, entity)
			}
		})
	}
}

func BenchmarkDelete(b *testing.B) {
	world := ecs.NewDefaultWorld()
	if err := fillWorld(world); err != nil {
		b.FailNow()
	}

	for b.Loop() {
		entity, _ := ecs.Spawn(world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})
		ecs.Delete(world, entity)
	}
}

func BenchmarkGet(b *testing.B) {
	world := ecs.NewDefaultWorld()
	if err := fillWorld(world); err != nil {
		b.FailNow()
	}

	target, _ := ecs.Spawn(world, &emptyComponentA{}, &emptyComponentB{}, &emptyComponentC{}, &emptyComponentD{})

	if err := fillWorld(world); err != nil {
		b.FailNow()
	}

	b.Run("Get1", func(b *testing.B) {
		for b.Loop() {
			ecs.Get1[emptyComponentA](world, target)
		}
	})

	b.Run("Get2", func(b *testing.B) {
		for b.Loop() {
			ecs.Get2[emptyComponentA, emptyComponentB](world, target)
		}
	})

	b.Run("Get3", func(b *testing.B) {
		for b.Loop() {
			ecs.Get3[emptyComponentA, emptyComponentB, emptyComponentC](world, target)
		}
	})

	b.Run("Get4", func(b *testing.B) {
		for b.Loop() {
			ecs.Get4[emptyComponentA, emptyComponentB, emptyComponentC, emptyComponentD](world, target)
		}
	})
}

func BenchmarkHasComponent(b *testing.B) {
	for _, numberOfEntities := range []int{10, 100, 1_000, 10_000} {
		world := ecs.NewDefaultWorld()

		for range numberOfEntities {
			if err := fillWorld(world); err != nil {
				b.FailNow()
			}

			if _, err := ecs.Spawn(world, &emptyComponentA{}); err != nil {
				b.Fatal(err)
			}

			if _, err := ecs.Spawn(world); err != nil {
				b.Fatal(err)
			}
		}

		entity, err := ecs.Spawn(world, &emptyComponentA{})
		if err != nil {
			b.Fatal(err)
		}

		b.Run(fmt.Sprintf("ComponentFound-%d-Entities", numberOfEntities), func(b *testing.B) {
			var componentFound bool

			for b.Loop() {
				componentFound, err = ecs.HasComponent[emptyComponentA](world, entity)
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
				componentFound, err = ecs.HasComponent[emptyComponentB](world, entity)
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
		world := ecs.NewDefaultWorld()
		for range numberOfEntities {
			if err := fillWorld(world); err != nil {
				b.FailNow()
			}

			if _, err := ecs.Spawn(world, &emptyComponentA{}); err != nil {
				b.Fatal(err)
			}

			if _, err := ecs.Spawn(world); err != nil {
				b.Fatal(err)
			}
		}

		entity, err := ecs.Spawn(world, &emptyComponentA{})
		if err != nil {
			b.Fatal(err)
		}

		componentIdA := ecs.ComponentIdFor[emptyComponentA](world)
		componentIdB := ecs.ComponentIdFor[emptyComponentB](world)

		b.Run(fmt.Sprintf("ComponentFound-%d-Entities", numberOfEntities), func(b *testing.B) {
			var componentFound bool

			for b.Loop() {
				componentFound, err = ecs.HasComponentId(world, entity, componentIdA)
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
				componentFound, err = ecs.HasComponentId(world, entity, componentIdB)
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
		world := ecs.NewDefaultWorld()

		for range size {
			if err := fillWorld(world); err != nil {
				b.FailNow()
			}
		}

		b.Run(fmt.Sprintf("Query0-Basic-Size-%d", size), func(b *testing.B) {
			query := ecs.Query0[ecs.Default]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query0-With1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query0[ecs.With[emptyComponentC]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query0-Without1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query0[ecs.Without[emptyComponentC]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query1-Basic-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.Default]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query1-Optional-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.Optional1[emptyComponentA]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query1-ReadOnly-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.AllReadOnly]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query1-With1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.With[emptyComponentC]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query1-Without1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query1[emptyComponentA, ecs.Without[emptyComponentC]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query2-Basic-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.Default]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query2-Optional-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.Optional1[emptyComponentA]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query2-ReadOnly-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.AllReadOnly]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query2-With1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.With[emptyComponentC]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query2-Without1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query2[emptyComponentA, emptyComponentD, ecs.Without[emptyComponentC]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query3-Basic-Size-%d", size), func(b *testing.B) {
			query := ecs.Query3[emptyComponentA, emptyComponentD, emptyComponentC, ecs.Default]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query3-Optional-Size-%d", size), func(b *testing.B) {
			query := ecs.Query3[emptyComponentA, emptyComponentD, emptyComponentC, ecs.Optional1[emptyComponentA]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query3-ReadOnly-Size-%d", size), func(b *testing.B) {
			query := ecs.Query3[emptyComponentA, emptyComponentD, emptyComponentC, ecs.AllReadOnly]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query3-With1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query3[emptyComponentA, emptyComponentD, emptyComponentC, ecs.With[componentWithValue]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query3-Without1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query3[emptyComponentA, emptyComponentD, emptyComponentC, ecs.Without[componentWithValue]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query4-Basic-Size-%d", size), func(b *testing.B) {
			query := ecs.Query4[emptyComponentA, emptyComponentD, emptyComponentB, emptyComponentC, ecs.Default]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query4-Optional-Size-%d", size), func(b *testing.B) {
			query := ecs.Query4[emptyComponentA, emptyComponentD, emptyComponentB, emptyComponentC, ecs.Optional1[emptyComponentA]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query4-ReadOnly-Size-%d", size), func(b *testing.B) {
			query := ecs.Query4[emptyComponentA, emptyComponentD, emptyComponentB, emptyComponentC, ecs.AllReadOnly]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query4-With1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query4[emptyComponentA, emptyComponentD, emptyComponentB, emptyComponentC, ecs.With[componentWithValue]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})

		b.Run(fmt.Sprintf("Query4-Without1-Size-%d", size), func(b *testing.B) {
			query := ecs.Query4[emptyComponentA, emptyComponentD, emptyComponentB, emptyComponentC, ecs.Without[componentWithValue]]{}

			err := query.Prepare(world, nil)
			if err != nil {
				b.FailNow()
			}

			for b.Loop() {
				query.Exec(world)
			}
		})
	}
}

// Fills the world with 7 different archetypes
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
