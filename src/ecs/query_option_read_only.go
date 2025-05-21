package ecs

type ReadOnlyComponents interface {
	getReadonlyComponentIds(world *World) (readOnlyComponentIds []ComponentId, isAllReadyOnly bool)
}

// All components in the query results are mutable.
type NoReadOnly struct{}

// All components in the query results are read-only. The components in the query results
// will be temporary copies.
//
// This allows high parallelization of systems that use this query as a system parameter.
type AllReadOnly struct{}

// Component in the query results is read-only.
// The component in the query results will be a temporary copy.
//
// This allows parallelization of systems that use this query as a system parameter.
type ReadOnly1[A IComponent] struct{}

// Components in the query results are read-only.
// The components in the query results will be temporary copies.
//
// This allows parallelization of systems that use this query as a system parameter.
type ReadOnly2[A, B IComponent] struct{}

// Components in the query results are read-only.
// The components in the query results will be temporary copies.
//
// This allows parallelization of systems that use this query as a system parameter.
type ReadOnly3[A, B, C IComponent] struct{}

// Components in the query results are read-only.
// The components in the query results will be temporary copies.
//
// This allows parallelization of systems that use this query as a system parameter.
type ReadOnly4[A, B, C, D IComponent] struct{}

func (o NoReadOnly) getReadonlyComponentIds(world *World) (readOnlyComponentIds []ComponentId, isAllReadyOnly bool) {
	return []ComponentId{}, false
}

func (o AllReadOnly) getReadonlyComponentIds(world *World) (readOnlyComponentIds []ComponentId, isAllReadyOnly bool) {
	return []ComponentId{}, true
}

func (o ReadOnly1[A]) getReadonlyComponentIds(world *World) (readOnlyComponentIds []ComponentId, isAllReadyOnly bool) {
	return []ComponentId{
		ComponentIdFor[A](world),
	}, false
}

func (o ReadOnly2[A, B]) getReadonlyComponentIds(world *World) (readOnlyComponentIds []ComponentId, isAllReadyOnly bool) {
	return []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
	}, false
}

func (o ReadOnly3[A, B, C]) getReadonlyComponentIds(world *World) (readOnlyComponentIds []ComponentId, isAllReadyOnly bool) {
	return []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
		ComponentIdFor[C](world),
	}, false
}

func (o ReadOnly4[A, B, C, D]) getReadonlyComponentIds(world *World) (readOnlyComponentIds []ComponentId, isAllReadyOnly bool) {
	return []ComponentId{
		ComponentIdFor[A](world),
		ComponentIdFor[B](world),
		ComponentIdFor[C](world),
		ComponentIdFor[D](world),
	}, false
}

func (readOnly AllReadOnly) getCombinedQueryOptions(world *World) (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, NoOptional, AllReadOnly, NotLazy, DefaultWorld]](world)
}
func (readOnly ReadOnly1[A]) getCombinedQueryOptions(world *World) (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, NoOptional, ReadOnly1[A], NotLazy, DefaultWorld]](world)
}
func (readOnly ReadOnly2[A, B]) getCombinedQueryOptions(world *World) (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, NoOptional, ReadOnly2[A, B], NotLazy, DefaultWorld]](world)
}
func (readOnly ReadOnly3[A, B, C]) getCombinedQueryOptions(world *World) (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, NoOptional, ReadOnly3[A, B, C], NotLazy, DefaultWorld]](world)
}
func (readOnly ReadOnly4[A, B, C, D]) getCombinedQueryOptions(world *World) (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, NoOptional, ReadOnly4[A, B, C, D], NotLazy, DefaultWorld]](world)
}
