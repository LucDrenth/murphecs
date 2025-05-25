package app

import (
	"fmt"
	"reflect"
	"slices"
)

type resourceId reflect.Type
type Resource any

type resourceStorage struct {
	resources            map[resourceId]Resource
	blacklistedResources []resourceId // resources that may not be added to this resourceStorage
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
	if resource == nil {
		return ErrResourceIsNil
	}

	resourceType := reflect.TypeOf(resource)
	if resourceType.Kind() != reflect.Pointer {
		// resource must be passed by reference because if it is not, we can never retrieve it by reference
		return ErrResourceNotAPointer
	}

	if resourceType.Elem().Kind() != reflect.Struct {
		return ErrResourceNotAStruct
	}

	resourceId := reflectTypeToResourceId(resourceType)

	if slices.Contains(s.blacklistedResources, resourceId) {
		return fmt.Errorf("%w: blacklisted", ErrResourceTypeNotAllowed)
	}

	if _, exists := s.resources[resourceId]; exists {
		return ErrResourceAlreadyPresent
	}

	s.resources[resourceId] = resource

	return nil
}

func registerBlacklistedResource[T Resource](storage *resourceStorage) error {
	resourceType := reflect.TypeFor[T]()
	return registerBlacklistedResourceType(resourceType, storage)
}

func registerBlacklistedResourceType(resourceType reflect.Type, storage *resourceStorage) error {
	if resourceType.Kind() != reflect.Pointer {
		// resource must be passed by reference because if it is not, we can never retrieve it by reference
		return ErrResourceNotAPointer
	}

	resourceId := reflectTypeToResourceId(resourceType)

	if slices.Contains(storage.blacklistedResources, resourceId) {
		return ErrResourceAlreadyPresent
	}

	storage.blacklistedResources = append(storage.blacklistedResources, resourceId)

	return nil
}

func getResourceFromStorage[T Resource](s *resourceStorage) (result T, err error) {
	resourceType := reflect.TypeFor[T]()
	resourceId := reflectTypeToResourceId(resourceType)

	untypedResource, exists := s.resources[resourceId]
	if !exists {
		return result, ErrResourceNotFound
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
	resourceId := reflectTypeToResourceId(resourceType)

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

func reflectTypeToResourceId(resourceType reflect.Type) resourceId {
	if resourceType.Kind() == reflect.Pointer {
		resourceType = resourceType.Elem()
	}

	return resourceId(resourceType)
}

func getResourceDebugType(resource Resource) string {
	return reflect.TypeOf(resource).String()
}
