package ecs

import "reflect"

// Optional wraps a component type to mark it as optional directly in a query's component list,
// e.g. Query1[Optional[NPC], Default] instead of Query1[NPC, Optional1[NPC]].
//
// Entities do not have to have the wrapped component in order to be included in the query result.
// If the entity has the component, Present is true and Value holds it. If not, Present is false and
// Value is the zero value of C.
type Optional[C AnyComponent] struct {
	Value   C
	Present bool
}

func (Optional[C]) RequiredComponents() []AnyComponent {
	return []AnyComponent{}
}

func (Optional[C]) optionalComponentType() reflect.Type {
	return reflect.TypeFor[C]()
}

// isOptionalQueryComponent is implemented by every instantiation of [Optional], regardless of C. It
// lets query preparation detect an Optional[C] component type and resolve C via reflection, since C
// itself is otherwise inaccessible from a generic type parameter at runtime.
type isOptionalQueryComponent interface {
	optionalComponentType() reflect.Type
}
