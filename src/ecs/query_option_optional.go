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

// Entities will be in query result even if any of these are not present. Components that are
// not present will be nil in the query result.
type Optional5[A, B, C, D, E IComponent] struct{}

// Entities will be in query result even if any of these are not present. Components that are
// not present will be nil in the query result.
type Optional6[A, B, C, D, E, F IComponent] struct{}

// Entities will be in query result even if any of these are not present. Components that are
// not present will be nil in the query result.
type Optional7[A, B, C, D, E, F, G IComponent] struct{}

// Entities will be in query result even if any of these are not present. Components that are
// not present will be nil in the query result.
type Optional8[A, B, C, D, E, F, G, H IComponent] struct{}

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

func (o Optional5[A, B, C, D, E]) getOptionalComponentIds(world *World) []ComponentId {
	return []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
		ComponentIdFor[C](world),
		ComponentIdFor[D](world),
		ComponentIdFor[E](world),
	}
}

func (o Optional6[A, B, C, D, E, F]) getOptionalComponentIds(world *World) []ComponentId {
	return []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
		ComponentIdFor[C](world),
		ComponentIdFor[D](world),
		ComponentIdFor[E](world),
		ComponentIdFor[F](world),
	}
}

func (o Optional7[A, B, C, D, E, F, G]) getOptionalComponentIds(world *World) []ComponentId {
	return []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
		ComponentIdFor[C](world),
		ComponentIdFor[D](world),
		ComponentIdFor[E](world),
		ComponentIdFor[F](world),
		ComponentIdFor[G](world),
	}
}

func (o Optional8[A, B, C, D, E, F, G, H]) getOptionalComponentIds(world *World) []ComponentId {
	return []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
		ComponentIdFor[C](world),
		ComponentIdFor[D](world),
		ComponentIdFor[E](world),
		ComponentIdFor[F](world),
		ComponentIdFor[G](world),
		ComponentIdFor[H](world),
	}
}

func (optional Optional1[A]) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional1[A], NotLazy, DefaultWorld]](world)
}
func (optional Optional2[A, B]) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional2[A, B], NotLazy, DefaultWorld]](world)
}
func (optional Optional3[A, B, C]) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional3[A, B, C], NotLazy, DefaultWorld]](world)
}
func (optional Optional4[A, B, C, D]) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional4[A, B, C, D], NotLazy, DefaultWorld]](world)
}
func (optional Optional5[A, B, C, D, E]) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional5[A, B, C, D, E], NotLazy, DefaultWorld]](world)
}
func (optional Optional6[A, B, C, D, E, F]) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional6[A, B, C, D, E, F], NotLazy, DefaultWorld]](world)
}
func (optional Optional7[A, B, C, D, E, F, G]) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional7[A, B, C, D, E, F, G], NotLazy, DefaultWorld]](world)
}
func (optional Optional8[A, B, C, D, E, F, G, H]) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, Optional8[A, B, C, D, E, F, G, H], NotLazy, DefaultWorld]](world)
}
