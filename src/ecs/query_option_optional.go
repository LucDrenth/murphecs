package ecs

type OptionalComponents interface {
	getOptionalComponentTypes() []ComponentType
}

// Entities have to have all components in order to be in the query result
type AllRequired struct{}

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

func (o AllRequired) getOptionalComponentTypes() []ComponentType {
	return []ComponentType{}
}

func (o Optional1[A]) getOptionalComponentTypes() []ComponentType {
	return []ComponentType{
		GetComponentType[A](),
	}
}

func (o Optional2[A, B]) getOptionalComponentTypes() []ComponentType {
	return []ComponentType{
		GetComponentType[A](),
		GetComponentType[B](),
	}
}

func (o Optional3[A, B, C]) getOptionalComponentTypes() []ComponentType {
	return []ComponentType{
		GetComponentType[A](),
		GetComponentType[B](),
		GetComponentType[C](),
	}
}

func (o Optional4[A, B, C, D]) getOptionalComponentTypes() []ComponentType {
	return []ComponentType{
		GetComponentType[A](),
		GetComponentType[B](),
		GetComponentType[C](),
		GetComponentType[D](),
	}
}
