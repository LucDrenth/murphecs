package app

import (
	"fmt"
	"reflect"
)

type resourceId reflect.Type
type Resource any

type resourceStorage struct {
	resources map[resourceId]Resource
}

func newResourceStorage() resourceStorage {
	return resourceStorage{
		resources: map[resourceId]Resource{},
	}
}

// Return an error if:
//   - resource is not passed by reference
//   - resource is already present
func (s *resourceStorage) add(resource Resource) error {
	resourceType := reflect.TypeOf(resource)
	if resourceType.Kind() != reflect.Pointer {
		// resource must be passed by reference because if it is not, we can never get it by reference
		return ErrResourceNotAPointer
	}

	resourceId := reflectTypeToComponentId(resourceType)
	if _, exists := s.resources[resourceId]; exists {
		return fmt.Errorf("%w: %s", ErrResourceAlreadyPresent, resourceId.String())
	}

	s.resources[resourceId] = resource

	return nil
}

func getResourceFromStorage[T Resource](s *resourceStorage) (result T, err error) {
	resourceType := reflect.TypeFor[T]()
	resourceId := reflectTypeToComponentId(resourceType)

	untypedResource, exists := s.resources[resourceId]
	if !exists {
		return result, fmt.Errorf("%w: %s", ErrResourceNotFound, resourceId.String())
	}

	if untypedResource == nil {
		return result, ErrResourceTypeNotValid
	}

	ok := false

	if resourceType.Kind() != reflect.Pointer {
		untypedResource = reflect.ValueOf(untypedResource).Elem().Interface()
	}

	result, ok = untypedResource.(T)
	if !ok {
		return result, fmt.Errorf("%w: failed to cast resource to %s", ErrResourceTypeNotValid, resourceId.String())
	}

	return result, nil
}

// getReflectResource returns a pointer to the resource, regardless if wether resourceType is a pointer or not.
func (s *resourceStorage) getReflectResource(resourceType reflect.Type) (result reflect.Value, err error) {
	resourceId := reflectTypeToComponentId(resourceType)

	untypedResource, exists := s.resources[resourceId]
	if !exists {
		return result, ErrResourceNotFound
	}

	if untypedResource == nil {
		return result, ErrResourceTypeNotValid
	}

	result = reflect.ValueOf(untypedResource)

	return result, nil
}

func reflectTypeToComponentId(resourceType reflect.Type) resourceId {
	if resourceType.Kind() == reflect.Pointer {
		resourceType = resourceType.Elem()
	}

	return resourceId(resourceType)
}
