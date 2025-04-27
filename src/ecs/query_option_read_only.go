package ecs

type ReadOnlyComponents interface {
	getReadonlyComponentTypes() (readOnlyComponentTypes []ComponentType, isAllReadyOnly bool)
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

func (o NoReadOnly) getReadonlyComponentTypes() (readOnlyComponentTypes []ComponentType, isAllReadyOnly bool) {
	return []ComponentType{}, false
}

func (o AllReadOnly) getReadonlyComponentTypes() (readOnlyComponentTypes []ComponentType, isAllReadyOnly bool) {
	return []ComponentType{}, true
}

func (o ReadOnly1[A]) getReadonlyComponentTypes() (readOnlyComponentTypes []ComponentType, isAllReadyOnly bool) {
	return []ComponentType{
		GetComponentType[A](),
	}, false
}

func (o ReadOnly2[A, B]) getReadonlyComponentTypes() (readOnlyComponentTypes []ComponentType, isAllReadyOnly bool) {
	return []ComponentType{
		GetComponentType[A](),
		GetComponentType[B](),
	}, false
}

func (o ReadOnly3[A, B, C]) getReadonlyComponentTypes() (readOnlyComponentTypes []ComponentType, isAllReadyOnly bool) {
	return []ComponentType{
		GetComponentType[A](),
		GetComponentType[B](),
		GetComponentType[C](),
	}, false
}

func (o ReadOnly4[A, B, C, D]) getReadonlyComponentTypes() (readOnlyComponentTypes []ComponentType, isAllReadyOnly bool) {
	return []ComponentType{
		GetComponentType[A](),
		GetComponentType[B](),
		GetComponentType[C](),
		GetComponentType[D](),
	}, false
}
