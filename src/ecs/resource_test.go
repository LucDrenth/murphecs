package ecs

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testResourceInterface interface {
	Get() int
	Increment()
}

type testResourceInterfaceA struct {
	value int
}

func (a *testResourceInterfaceA) Get() int {
	return a.value
}

func (a *testResourceInterfaceA) Increment() {
	a.value += 1
}

func TestAddStructToResourceStorage(t *testing.T) {
	type resourceA struct{}

	t.Run("fails to add resource if it is blacklisted", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := RegisterBlacklistedResource[*resourceA](&storage)
		assert.NoError(err)
		err = storage.Add(&resourceA{})
		assert.ErrorIs(err, ErrResourceTypeNotAllowed)
	})

	t.Run("fails if resource is not a valid type", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()

		err := storage.Add(&[]int{1, 2})
		assert.ErrorIs(err, ErrResourceTypeNotValid)
		err = storage.Add(new("invalid resource type"))
		assert.ErrorIs(err, ErrResourceTypeNotValid)
		err = storage.Add(new(100))
		assert.ErrorIs(err, ErrResourceTypeNotValid)
		err = storage.Add(new(func() {}))
		assert.ErrorIs(err, ErrResourceTypeNotValid)
	})

	t.Run("fails to add resource if it is untyped and nil", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.Add(nil)
		assert.ErrorIs(err, ErrResourceIsNil)
	})

	t.Run("fails to add resource if it already exists", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.Add(&resourceA{})
		assert.NoError(err)
		err = storage.Add(&resourceA{})
		assert.ErrorIs(err, ErrResourceAlreadyPresent)
	})

	t.Run("successfully adds struct resource by value", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.Add(resourceA{})
		assert.NoError(err)
	})

	t.Run("successfully adds struct resource by reference", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.Add(&resourceA{})
		assert.NoError(err)
	})
}

func TestAddInterfaceToResourceStorage(t *testing.T) {
	t.Run("interface resource and its struct implementation are treated as the same resource type", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()

		var resource testResourceInterface = &testResourceInterfaceA{}
		err := storage.Add(resource)
		assert.NoError(err)

		err = storage.Add(&testResourceInterfaceA{})
		assert.ErrorIs(err, ErrResourceAlreadyPresent)
	})

	t.Run("successfully adds interface resource that is passed by reference", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()

		var resource testResourceInterface = &testResourceInterfaceA{}
		err := storage.Add(&resource)
		assert.NoError(err)
	})

	t.Run("reference to interface resource and its struct implementation are not treated as the same resource type", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()

		var resource testResourceInterface = &testResourceInterfaceA{}
		err := storage.Add(&resource)
		assert.NoError(err)

		err = storage.Add(&testResourceInterfaceA{})
		assert.NoError(err)
	})

	t.Run("successfully adds interface resource that has value nil", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()

		var resource testResourceInterface = nil
		err := storage.Add(&resource)
		assert.NoError(err)
	})
}

func TestRegisterBlacklistedResource(t *testing.T) {
	type resourceA struct{}

	t.Run("fails to blacklist resource if it is already present", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := RegisterBlacklistedResource[*resourceA](&storage)
		assert.NoError(err)
		err = RegisterBlacklistedResource[*resourceA](&storage)
		assert.ErrorIs(err, ErrResourceAlreadyPresent)
		assert.Len(storage.blacklistedResources, 1)
	})

	t.Run("successfully registers blacklisted resource", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := RegisterBlacklistedResource[resourceA](&storage)
		assert.NoError(err)
		assert.Len(storage.blacklistedResources, 1)
	})

	t.Run("successfully registers blacklisted resource pointer", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := RegisterBlacklistedResource[*resourceA](&storage)
		assert.NoError(err)
		assert.Len(storage.blacklistedResources, 1)
	})
}

func TestGetResourceFromStorage(t *testing.T) {
	type resourceA struct{ value int }

	t.Run("fails to get resource that was not added", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()

		_, err := GetResourceFromStorage[resourceA](&storage)
		assert.ErrorIs(err, ErrResourceNotFound)

		_, err = GetResourceFromStorage[*resourceA](&storage)
		assert.ErrorIs(err, ErrResourceNotFound)
	})

	t.Run("successfully gets struct resource copy", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.Add(&resourceA{value: 10})
		assert.NoError(err)

		resource, err := GetResourceFromStorage[resourceA](&storage)
		assert.NoError(err)
		assert.Equal(10, resource.value)
	})

	t.Run("struct resource copy can not be mutated", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.Add(&resourceA{value: 0})
		assert.NoError(err)

		resource, err := GetResourceFromStorage[resourceA](&storage)
		assert.NoError(err)
		resource.value = 10

		resource, err = GetResourceFromStorage[resourceA](&storage)
		assert.NoError(err)
		assert.NotEqual(10, resource.value)
	})

	t.Run("successfully gets struct resource pointer", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.Add(&resourceA{value: 10})
		assert.NoError(err)

		resource, err := GetResourceFromStorage[*resourceA](&storage)
		assert.NoError(err)
		assert.Equal(10, resource.value)
	})

	t.Run("struct resource passed by reference can be mutated", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.Add(&resourceA{value: 0})
		assert.NoError(err)

		resource, err := GetResourceFromStorage[*resourceA](&storage)
		assert.NoError(err)
		resource.value = 10

		resource, err = GetResourceFromStorage[*resourceA](&storage)
		assert.NoError(err)
		assert.Equal(10, resource.value)
	})

	t.Run("interface resource that is passed by reference can be retrieved and mutated", func(t *testing.T) {
		assert := assert.New(t)

		var resource testResourceInterface = &testResourceInterfaceA{}
		storage := newResourceStorage()
		err := storage.Add(&resource)
		assert.NoError(err)

		resource, err = GetResourceFromStorage[testResourceInterface](&storage)
		assert.NoError(err)
		resource.Increment()

		resource, err = GetResourceFromStorage[testResourceInterface](&storage)
		assert.NoError(err)
		assert.Equal(1, resource.Get())
	})

	t.Run("interface resource that is passed by value can be retrieved by its implementation and can be mutated", func(t *testing.T) {
		assert := assert.New(t)

		var resource testResourceInterface = &testResourceInterfaceA{}
		storage := newResourceStorage()
		err := storage.Add(resource)
		assert.NoError(err)

		resource, err = GetResourceFromStorage[*testResourceInterfaceA](&storage)
		assert.NoError(err)
		resource.Increment()

		resource, err = GetResourceFromStorage[*testResourceInterfaceA](&storage)
		assert.NoError(err)
		assert.Equal(1, resource.Get())
	})
}

func TestGetReflectResource(t *testing.T) {
	type resourceA struct{ value int }

	t.Run("fails to get resource that was not added", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		_, err := storage.GetReflectResource(reflect.TypeFor[resourceA]())
		assert.ErrorIs(err, ErrResourceNotFound)

		_, err = storage.GetReflectResource(reflect.TypeFor[*resourceA]())
		assert.ErrorIs(err, ErrResourceNotFound)
	})

	t.Run("successfully gets resource", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.Add(&resourceA{value: 10})
		assert.NoError(err)

		reflectedResource, err := storage.GetReflectResource(reflect.TypeFor[resourceA]())
		assert.NoError(err)

		resourceReference, ok := reflect.TypeAssert[*resourceA](reflectedResource)
		assert.True(ok)
		assert.Equal(10, resourceReference.value)

		resourceCopy, ok := reflect.TypeAssert[resourceA](reflectedResource.Elem())
		assert.True(ok)
		assert.Equal(10, resourceCopy.value)
	})

	t.Run("resource can be mutated", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.Add(&resourceA{value: 0})
		assert.NoError(err)

		reflectedResource, err := storage.GetReflectResource(reflect.TypeFor[*resourceA]())
		assert.NoError(err)
		resource, ok := reflect.TypeAssert[*resourceA](reflectedResource)
		assert.True(ok)
		resource.value = 10

		reflectedResource, err = storage.GetReflectResource(reflect.TypeFor[*resourceA]())
		assert.NoError(err)
		resource, ok = reflect.TypeAssert[*resourceA](reflectedResource)
		assert.True(ok)
		assert.Equal(10, resource.value)
	})
}

func TestResourceId(t *testing.T) {
	type resourceA struct{}
	type resourceB struct{}

	t.Run("resource return the same result, regardless of wether its passed by reference", func(t *testing.T) {
		assert := assert.New(t)

		a := reflectTypeToResourceId(reflect.TypeFor[resourceA]())
		b := reflectTypeToResourceId(reflect.TypeFor[*resourceA]())
		assert.Equal(a, b)
	})

	t.Run("resource ids for different resources are unique", func(t *testing.T) {
		assert := assert.New(t)

		a := reflectTypeToResourceId(reflect.TypeFor[resourceA]())
		b := reflectTypeToResourceId(reflect.TypeFor[resourceB]())
		assert.NotEqual(a, b)
	})

	t.Run("resource ids for different resources of the same name and structure", func(t *testing.T) {
		assert := assert.New(t)

		oldA := resourceA{}
		type resourceA struct{}
		newA := resourceA{}

		assert.NotEqual(oldA, newA)
	})
}
