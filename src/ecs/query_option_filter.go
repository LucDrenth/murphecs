package ecs

type QueryParamFilter interface {
	getComponents() []ComponentType
	getFilterType() filterType
}
type NoFilter struct{}
type With[A IComponent] struct{}
type Without[A IComponent] struct{}
type Or[A, B QueryParamFilter] struct{}
type And[A, B QueryParamFilter] struct{}

func (filter NoFilter) getComponents() []ComponentType {
	return []ComponentType{}
}

func (filter And[A, B]) getComponents() []ComponentType {
	return []ComponentType{}
}

func (filter Or[A, B]) getComponents() []ComponentType {
	return []ComponentType{}
}

func (filter With[A]) getComponents() []ComponentType {
	return []ComponentType{GetComponentType[A]()}
}

func (filter Without[A]) getComponents() []ComponentType {
	return []ComponentType{GetComponentType[A]()}
}

func (filter NoFilter) getFilterType() filterType {
	return filterTypeNone
}

func (filter With[A]) getFilterType() filterType {
	return filterTypeWith
}

func (filter Without[A]) getFilterType() filterType {
	return filterTypeWithout
}

func (filter And[A, B]) getFilterType() filterType {
	return filterTypeAnd
}

func (filter Or[A, B]) getFilterType() filterType {
	return filterTypeOr
}

type QueryFilter interface {
	// Validate that entityData satisfies the filter
	Validate(*EntityData) bool
}
type queryFilterAnd struct {
	a QueryFilter
	b QueryFilter
}
type queryFilterOr struct {
	a QueryFilter
	b QueryFilter
}
type queryFilterWith struct {
	c []ComponentType
}
type queryFilterWithout struct {
	c []ComponentType
}

func (filter queryFilterAnd) Validate(e *EntityData) bool {
	return filter.a.Validate(e) && filter.b.Validate(e)
}

func (filter queryFilterOr) Validate(e *EntityData) bool {
	return filter.a.Validate(e) || filter.b.Validate(e)
}

func (filter queryFilterWith) Validate(e *EntityData) bool {
	for _, c := range filter.c {
		if !e.hasComponent(c) {
			return false
		}
	}

	return true
}

func (filter queryFilterWithout) Validate(e *EntityData) bool {
	for _, c := range filter.c {
		if e.hasComponent(c) {
			return false
		}
	}

	return true
}
