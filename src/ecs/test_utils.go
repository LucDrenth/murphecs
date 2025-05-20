package ecs

import (
	"fmt"
	"reflect"
)

// componentWithPointers is useful to test garbage collection. it is initialized with
type ComponentWithPointers struct {
	Component
	name     *string
	aMap     map[string]int
	intSlice []*int
}

var (
	componentWithPointersName        string = "a name"
	componentWithPointersSliceLength int    = 10
)

func CreateComponentWithPointers() *ComponentWithPointers {
	item := ComponentWithPointers{
		name: &componentWithPointersName,
		aMap: map[string]int{
			"1": 1,
			"3": 3,
		},
		intSlice: []*int{},
	}
	for i := range componentWithPointersSliceLength {
		item.intSlice = append(item.intSlice, &i)
	}

	return &item
}

// validates that ComponentWithPointers matches what is created by CreateComponentWithPointers. Returns an error if not valid.
func (item *ComponentWithPointers) Validate() error {
	if item.name == nil {
		return fmt.Errorf("name is nil")
	}

	if *item.name != componentWithPointersName {
		return fmt.Errorf("invalid name: %s", *item.name)
	}

	expectedMap := map[string]int{
		"1": 1,
		"3": 3,
	}
	if !reflect.DeepEqual(expectedMap, item.aMap) {
		return fmt.Errorf("aMap does not match")
	}

	if len(item.intSlice) != componentWithPointersSliceLength {
		return fmt.Errorf("invalid intSlice length: %d", len(item.intSlice))
	}

	for i := range componentWithPointersSliceLength {
		value := *item.intSlice[i]
		if i != value {
			return fmt.Errorf("invalid intSlice value at index %d: %d", i, value)

		}
	}

	return nil
}
