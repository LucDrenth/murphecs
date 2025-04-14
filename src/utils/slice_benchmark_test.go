package utils

import (
	"fmt"
	"testing"
)

var sizes = []int{10, 100, 1_000, 10_000}

func BenchmarkGetFirstDuplicate(b *testing.B) {
	for _, size := range sizes {
		setup := func() []int {
			data := make([]int, size)

			for i := range size {
				data[i] = i
			}

			return data
		}

		b.Run(fmt.Sprintf("AtTheStart-Size-%d", size), func(b *testing.B) {
			data := setup()
			data[5] = 1

			for b.Loop() {
				GetFirstDuplicate(data)
			}
		})

		b.Run(fmt.Sprintf("InTheMiddle-Size-%d", size), func(b *testing.B) {
			data := setup()
			data[size/2] = 1

			for b.Loop() {
				GetFirstDuplicate(data)
			}
		})

		b.Run(fmt.Sprintf("AtTheEnd-Size-%d", size), func(b *testing.B) {
			data := setup()
			data[size-5] = 1

			for b.Loop() {
				GetFirstDuplicate(data)
			}
		})

		b.Run(fmt.Sprintf("NoDuplicates-Size-%d", size), func(b *testing.B) {
			data := setup()

			for b.Loop() {
				GetFirstDuplicate(data)
			}
		})
	}
}

func BenchmarkRemoveFromSlice(b *testing.B) {
	setup := func(size int) []int {
		return make([]int, size)
	}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("AtTheStart-Size-%d", size), func(b *testing.B) {
			for b.Loop() {
				data := setup(size)
				// We'd want to reset the timer here, but that makes this benchmark hang

				RemoveFromSlice(&data, 5)
			}
		})

		b.Run(fmt.Sprintf("InTheMiddle-Size-%d", size), func(b *testing.B) {
			for b.Loop() {
				data := setup(size)
				// We'd want to reset the timer here, but that makes this benchmark hang

				RemoveFromSlice(&data, size/2)
			}
		})

		b.Run(fmt.Sprintf("AtTheEnd-Size-%d", size), func(b *testing.B) {
			for b.Loop() {
				data := setup(size)
				// We'd want to reset the timer here, but that makes this benchmark hang

				RemoveFromSlice(&data, size-5)
			}
		})
	}
}

func BenchmarkRemoveFromSliceAndMaintainOrder(b *testing.B) {
	for _, size := range sizes {
		setup := func() []int {
			return make([]int, size)
		}

		b.Run(fmt.Sprintf("AtTheStart-Size-%d", size), func(b *testing.B) {
			for b.Loop() {
				data := setup()
				// We'd want to reset the timer here, but that makes this benchmark hang

				RemoveFromSliceAndMaintainOrder(&data, 5)
			}
		})

		b.Run(fmt.Sprintf("InTheMiddle-Size-%d", size), func(b *testing.B) {
			for b.Loop() {
				data := setup()
				// We'd want to reset the timer here, but that makes this benchmark hang

				RemoveFromSliceAndMaintainOrder(&data, size/2)
			}
		})

		b.Run(fmt.Sprintf("AtTheEnd-Size-%d", size), func(b *testing.B) {
			for b.Loop() {
				data := setup()
				// We'd want to reset the timer here, but that makes this benchmark hang

				RemoveFromSliceAndMaintainOrder(&data, size-5)
			}
		})
	}
}
