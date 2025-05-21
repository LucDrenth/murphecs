package ecs

import (
	"fmt"
	"testing"
)

func BenchmarkHashComponentIds(b *testing.B) {
	type component1 struct{ Component }
	type component2 struct{ Component }
	type component3 struct{ Component }
	type component4 struct{ Component }
	type component5 struct{ Component }
	type component6 struct{ Component }
	type component7 struct{ Component }
	type component8 struct{ Component }
	type component9 struct{ Component }
	type component10 struct{ Component }
	type component11 struct{ Component }
	type component12 struct{ Component }
	type component13 struct{ Component }
	type component14 struct{ Component }
	type component15 struct{ Component }
	type component16 struct{ Component }
	type component17 struct{ Component }
	type component18 struct{ Component }
	type component19 struct{ Component }
	type component20 struct{ Component }
	type component21 struct{ Component }
	type component22 struct{ Component }
	type component23 struct{ Component }
	type component24 struct{ Component }
	type component25 struct{ Component }
	type component26 struct{ Component }
	type component27 struct{ Component }
	type component28 struct{ Component }
	type component29 struct{ Component }
	type component30 struct{ Component }
	type component31 struct{ Component }
	type component32 struct{ Component }

	world := NewDefaultWorld()

	allComponentIds := []ComponentId{
		ComponentIdFor[component1](&world),
		ComponentIdFor[component2](&world),
		ComponentIdFor[component3](&world),
		ComponentIdFor[component4](&world),
		ComponentIdFor[component5](&world),
		ComponentIdFor[component6](&world),
		ComponentIdFor[component7](&world),
		ComponentIdFor[component8](&world),
		ComponentIdFor[component9](&world),
		ComponentIdFor[component10](&world),
		ComponentIdFor[component11](&world),
		ComponentIdFor[component12](&world),
		ComponentIdFor[component13](&world),
		ComponentIdFor[component14](&world),
		ComponentIdFor[component15](&world),
		ComponentIdFor[component16](&world),
		ComponentIdFor[component17](&world),
		ComponentIdFor[component18](&world),
		ComponentIdFor[component19](&world),
		ComponentIdFor[component20](&world),
		ComponentIdFor[component21](&world),
		ComponentIdFor[component22](&world),
		ComponentIdFor[component23](&world),
		ComponentIdFor[component24](&world),
		ComponentIdFor[component25](&world),
		ComponentIdFor[component26](&world),
		ComponentIdFor[component27](&world),
		ComponentIdFor[component28](&world),
		ComponentIdFor[component29](&world),
		ComponentIdFor[component30](&world),
		ComponentIdFor[component31](&world),
		ComponentIdFor[component32](&world),
	}

	for i := range allComponentIds {
		b.Run(fmt.Sprintf("hash %d component ids", i+1), func(b *testing.B) {
			componentIds := allComponentIds[:i]

			for b.Loop() {
				hashComponentIds(componentIds)
			}
		})
	}
}
