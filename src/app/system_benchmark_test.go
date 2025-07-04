package app

import (
	"reflect"
	"testing"
)

// BenchmarkCallReflectedFunction benchmarks overhead of calling a reflected function, which is how
// systems are called.
func BenchmarkCallReflectedFunction(b *testing.B) {
	sum := func(a int, b int) int {
		return a + b
	}
	valueA := 8
	valueB := 2

	b.Run("call regular function", func(b *testing.B) {
		for b.Loop() {
			sum(valueA, valueB)
		}
	})

	b.Run("call reflected function", func(b *testing.B) {
		reflectedFunction := reflect.ValueOf(sum)
		params := []reflect.Value{
			reflect.ValueOf(valueA),
			reflect.ValueOf(valueB),
		}

		for b.Loop() {
			reflectedFunction.Call(params)
		}
	})
}
