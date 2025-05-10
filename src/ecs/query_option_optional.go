package ecs

type OptionalComponents interface {
	getOptionalComponentIds(world *World) []ComponentId
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

func (o NoOptional) getOptionalComponentIds(world *World) []ComponentId {
	return []ComponentId{}
}

func (o Optional1[A]) getOptionalComponentIds(world *World) []ComponentId {
	return []ComponentId{
		ComponentIdFor[A](world),
	}
}

func (o Optional2[A, B]) getOptionalComponentIds(world *World) []ComponentId {
	return []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
	}
}

func (o Optional3[A, B, C]) getOptionalComponentIds(world *World) []ComponentId {
	return []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
		ComponentIdFor[C](world),
	}
}

func (o Optional4[A, B, C, D]) getOptionalComponentIds(world *World) []ComponentId {
	return []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
		ComponentIdFor[C](world),
		ComponentIdFor[D](world),
	}
}

func (optional Optional1[A]) getCombinedQueryOptions(world *World) (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional1[A], NoReadOnly, NotLazy]](world)
}
func (optional Optional2[A, B]) getCombinedQueryOptions(world *World) (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional2[A, B], NoReadOnly, NotLazy]](world)
}
func (optional Optional3[A, B, C]) getCombinedQueryOptions(world *World) (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional3[A, B, C], NoReadOnly, NotLazy]](world)
}
func (optional Optional4[A, B, C, D]) getCombinedQueryOptions(world *World) (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional4[A, B, C, D], NoReadOnly, NotLazy]](world)
}
