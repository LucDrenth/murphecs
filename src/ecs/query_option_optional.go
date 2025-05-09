package ecs

type OptionalComponents interface {
	getOptionalComponentIds() []ComponentId
}

// Entities have to have all components in order to be in the query result
type NoOptional struct{}

// Entities will be in query result even if A is not present. If A is not present,
// it will be nil in the query result.
type Optional1[A IComponent] struct{}

// Entities will be in query result even if any of these are not present. Components that are
// not present will be nil in the query result.
type Optional2[A, B IComponent] struct{}

// Entities will be in query result even if any of these are not present. Components that are
// not present will be nil in the query result.
type Optional3[A, B, C IComponent] struct{}

// Entities will be in query result even if any of these are not present. Components that are
// not present will be nil in the query result.
type Optional4[A, B, C, D IComponent] struct{}

func (o NoOptional) getOptionalComponentIds() []ComponentId {
	return []ComponentId{}
}

func (o Optional1[A]) getOptionalComponentIds() []ComponentId {
	return []ComponentId{
		ComponentIdFor[A](),
	}
}

func (o Optional2[A, B]) getOptionalComponentIds() []ComponentId {
	return []ComponentId{
		ComponentIdFor[A](),
		ComponentIdFor[B](),
	}
}

func (o Optional3[A, B, C]) getOptionalComponentIds() []ComponentId {
	return []ComponentId{
		ComponentIdFor[A](),
		ComponentIdFor[B](),
		ComponentIdFor[C](),
	}
}

func (o Optional4[A, B, C, D]) getOptionalComponentIds() []ComponentId {
	return []ComponentId{
		ComponentIdFor[A](),
		ComponentIdFor[B](),
		ComponentIdFor[C](),
		ComponentIdFor[D](),
	}
}

func (optional Optional1[A]) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional1[A], NoReadOnly, NotLazy]]()
}
func (optional Optional2[A, B]) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional2[A, B], NoReadOnly, NotLazy]]()
}
func (optional Optional3[A, B, C]) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional3[A, B, C], NoReadOnly, NotLazy]]()
}
func (optional Optional4[A, B, C, D]) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional4[A, B, C, D], NoReadOnly, NotLazy]]()
}
