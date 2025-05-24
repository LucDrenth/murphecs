package ecs

type TargetWorld interface {
	GetWorldId() *WorldId
}

// DefaultWorld lets this query use the world of the SubApp it is used in.
type DefaultWorld struct{}

func (DefaultWorld) GetWorldId() *WorldId {
	return nil
}

func (DefaultWorld) GetCombinedQueryOptions(world *World) (CombinedQueryOptions, error) {
	return CombinedQueryOptions{}, nil
}
