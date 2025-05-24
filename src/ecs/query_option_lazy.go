package ecs

type IsQueryLazy interface {
	isLazy() bool
}

// NotLazy makes this query automatically get executed if it is in a system param.
type NotLazy struct{}

// Lazy makes this query not get automatically executed if it is in a system param.
// You'll have to execute it yourself by calling Exec.
//
// Lazy queries in system parameters will still get their query results from previous
// runs cleared.
type Lazy struct{}

func (NotLazy) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return CombinedQueryOptions{isLazy: false}, nil
}
func (Lazy) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return CombinedQueryOptions{isLazy: true}, nil
}

func (NotLazy) isLazy() bool {
	return false
}
func (Lazy) isLazy() bool {
	return true
}
