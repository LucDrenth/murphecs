package ecs

type ReadOnlyComponents interface {
	getReadonlyComponentIds() (readOnlyComponentIds []ComponentId, isAllReadyOnly bool)
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

func (o NoReadOnly) getReadonlyComponentIds() (readOnlyComponentIds []ComponentId, isAllReadyOnly bool) {
	return []ComponentId{}, false
}

func (o AllReadOnly) getReadonlyComponentIds() (readOnlyComponentIds []ComponentId, isAllReadyOnly bool) {
	return []ComponentId{}, true
}

func (o ReadOnly1[A]) getReadonlyComponentIds() (readOnlyComponentIds []ComponentId, isAllReadyOnly bool) {
	return []ComponentId{
		ComponentIdFor[A](),
	}, false
}

func (o ReadOnly2[A, B]) getReadonlyComponentIds() (readOnlyComponentIds []ComponentId, isAllReadyOnly bool) {
	return []ComponentId{
		ComponentIdFor[A](),
		ComponentIdFor[B](),
	}, false
}

func (o ReadOnly3[A, B, C]) getReadonlyComponentIds() (readOnlyComponentIds []ComponentId, isAllReadyOnly bool) {
	return []ComponentId{
		ComponentIdFor[A](),
		ComponentIdFor[B](),
		ComponentIdFor[C](),
	}, false
}

func (o ReadOnly4[A, B, C, D]) getReadonlyComponentIds() (readOnlyComponentIds []ComponentId, isAllReadyOnly bool) {
	return []ComponentId{
		ComponentIdFor[A](),
		ComponentIdFor[B](),
		ComponentIdFor[C](),
		ComponentIdFor[D](),
	}, false
}

func (readOnly AllReadOnly) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, NoOptional, AllReadOnly, NotLazy]]()
}
func (readOnly ReadOnly1[A]) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, NoOptional, ReadOnly1[A], NotLazy]]()
}
func (readOnly ReadOnly2[A, B]) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, NoOptional, ReadOnly2[A, B], NotLazy]]()
}
func (readOnly ReadOnly3[A, B, C]) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, NoOptional, ReadOnly3[A, B, C], NotLazy]]()
}
func (readOnly ReadOnly4[A, B, C, D]) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, NoOptional, ReadOnly4[A, B, C, D], NotLazy]]()
}
