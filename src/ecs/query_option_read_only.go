package ecs

type ReadOnlyComponents interface {
	isAllReadOnly() bool
}

// All components in the query results are mutable.
type NotAllReadOnly struct{}

// All components in the query results are read-only. The components in the query results
// will be temporary copies.
//
// This allows high parallelization of systems that use this query as a system parameter.
type AllReadOnly struct{}

func (o NotAllReadOnly) isAllReadOnly() bool {
	return false
}

func (o AllReadOnly) isAllReadOnly() bool {
	return true
}

func (readOnly AllReadOnly) getCombinedQueryOptions() (combinedQueryOptions, error) {
	return toCombinedQueryOptions[QueryOptions[NoFilter, AllReadOnly]]()
}
