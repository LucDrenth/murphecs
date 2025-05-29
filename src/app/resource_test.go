package app

import (
	"reflect"
	"testing"

	"github.com/lucdrenth/murphecs/src/utils"
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

	t.Run("fails to add resource if it is not passed by reference", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.add(resourceA{})
		assert.ErrorIs(err, ErrResourceNotAPointer)
	})

	t.Run("fails to add resource if it is blacklisted", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := registerBlacklistedResource[*resourceA](&storage)
		assert.NoError(err)
		err = storage.add(&resourceA{})
		assert.ErrorIs(err, ErrResourceTypeNotAllowed)
	})

	t.Run("fails if resource is not a valid type", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()

		err := storage.add(&[]int{1, 2})
		assert.ErrorIs(err, ErrResourceTypeNotValid)
		err = storage.add(utils.PointerTo("invalid resource type"))
		assert.ErrorIs(err, ErrResourceTypeNotValid)
		err = storage.add(utils.PointerTo(100))
		assert.ErrorIs(err, ErrResourceTypeNotValid)
		err = storage.add(utils.PointerTo(func() {}))
		assert.ErrorIs(err, ErrResourceTypeNotValid)
	})

	t.Run("fails to add resource if it is untyped and nil", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.add(nil)
		assert.ErrorIs(err, ErrResourceIsNil)
	})

	t.Run("fails to add resource if it already exists", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.add(&resourceA{})
		assert.NoError(err)
		err = storage.add(&resourceA{})
		assert.ErrorIs(err, ErrResourceAlreadyPresent)
	})

	t.Run("successfully adds struct resource", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.add(&resourceA{})
		assert.NoError(err)
	})
}

func TestAddInterfaceToResourceStorage(t *testing.T) {
	t.Run("interface resource and its struct implementation are treated as the same resource type", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()

		var resource testResourceInterface = &testResourceInterfaceA{}
		err := storage.add(resource)
		assert.NoError(err)

		err = storage.add(&testResourceInterfaceA{})
		assert.ErrorIs(err, ErrResourceAlreadyPresent)
	})

	t.Run("successfully adds interface resource that is passed by value", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()

		var log Logger = &NoOpLogger{}
		err := storage.add(log)
		assert.NoError(err)
	})

	t.Run("successfully adds interface resource that is passed by reference", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()

		var log Logger = &NoOpLogger{}
		err := storage.add(&log)
		assert.NoError(err)
	})

	t.Run("reference to interface resource and its struct implementation are not treated as the same resource type", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()

		var resource testResourceInterface = &testResourceInterfaceA{}
		err := storage.add(&resource)
		assert.NoError(err)

		err = storage.add(&testResourceInterfaceA{})
		assert.NoError(err)
	})

	t.Run("successfully adds interface resource that has value nil", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()

		var resource testResourceInterface = nil
		err := storage.add(&resource)
		assert.NoError(err)
	})
}

func TestRegisterBlacklistedResource(t *testing.T) {
	type resourceA struct{}

	t.Run("fails to blacklist resource if it is not passed by reference", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := registerBlacklistedResource[resourceA](&storage)
		assert.ErrorIs(err, ErrResourceNotAPointer)
	})

	t.Run("fails to blacklist resource if it is already present", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := registerBlacklistedResource[*resourceA](&storage)
		assert.NoError(err)
		err = registerBlacklistedResource[*resourceA](&storage)
		assert.ErrorIs(err, ErrResourceAlreadyPresent)
		assert.Equal(1, len(storage.blacklistedResources))
	})

	t.Run("successfully registers blacklisted resource", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := registerBlacklistedResource[*resourceA](&storage)
		assert.NoError(err)
		assert.Equal(1, len(storage.blacklistedResources))
	})
}

func TestGetResourceFromStorage(t *testing.T) {
	type resourceA struct{ value int }

	t.Run("fails to get resource that was not added", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()

		_, err := getResourceFromStorage[resourceA](&storage)
		assert.ErrorIs(err, ErrResourceNotFound)

		_, err = getResourceFromStorage[*resourceA](&storage)
		assert.ErrorIs(err, ErrResourceNotFound)
	})

	t.Run("successfully gets struct resource copy", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.add(&resourceA{value: 10})
		assert.NoError(err)

		resource, err := getResourceFromStorage[resourceA](&storage)
		assert.NoError(err)
		assert.Equal(10, resource.value)
	})

	t.Run("struct resource copy can not be mutated", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.add(&resourceA{value: 0})
		assert.NoError(err)

		resource, err := getResourceFromStorage[resourceA](&storage)
		assert.NoError(err)
		resource.value = 10

		resource, err = getResourceFromStorage[resourceA](&storage)
		assert.NoError(err)
		assert.NotEqual(10, resource.value)
	})

	t.Run("successfully gets struct resource pointer", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.add(&resourceA{value: 10})
		assert.NoError(err)

		resource, err := getResourceFromStorage[*resourceA](&storage)
		assert.NoError(err)
		assert.Equal(10, resource.value)
	})

	t.Run("struct resource passed by reference can be mutated", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.add(&resourceA{value: 0})
		assert.NoError(err)

		resource, err := getResourceFromStorage[*resourceA](&storage)
		assert.NoError(err)
		resource.value = 10

		resource, err = getResourceFromStorage[*resourceA](&storage)
		assert.NoError(err)
		assert.Equal(10, resource.value)
	})

	t.Run("interface resource that is passed by reference can be retrieved and mutated", func(t *testing.T) {
		assert := assert.New(t)

		var resource testResourceInterface = &testResourceInterfaceA{}
		storage := newResourceStorage()
		err := storage.add(&resource)
		assert.NoError(err)

		resource, err = getResourceFromStorage[testResourceInterface](&storage)
		assert.NoError(err)
		resource.Increment()

		resource, err = getResourceFromStorage[testResourceInterface](&storage)
		assert.NoError(err)
		assert.Equal(1, resource.Get())
	})

	t.Run("interface resource that is passed by value can be retrieved by its implementation and can be mutated", func(t *testing.T) {
		assert := assert.New(t)

		var resource testResourceInterface = &testResourceInterfaceA{}
		storage := newResourceStorage()
		err := storage.add(resource)
		assert.NoError(err)

		resource, err = getResourceFromStorage[*testResourceInterfaceA](&storage)
		assert.NoError(err)
		resource.Increment()

		resource, err = getResourceFromStorage[*testResourceInterfaceA](&storage)
		assert.NoError(err)
		assert.Equal(1, resource.Get())
	})
}

func TestGetReflectResource(t *testing.T) {
	type resourceA struct{ value int }

	t.Run("fails to get resource that was not added", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		_, err := storage.getReflectResource(reflect.TypeFor[resourceA]())
		assert.ErrorIs(err, ErrResourceNotFound)

		_, err = storage.getReflectResource(reflect.TypeFor[*resourceA]())
		assert.ErrorIs(err, ErrResourceNotFound)
	})

	t.Run("successfully gets resource", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.add(&resourceA{value: 10})
		assert.NoError(err)

		reflectedResource, err := storage.getReflectResource(reflect.TypeFor[resourceA]())
		assert.NoError(err)

		resourceReference, ok := reflectedResource.Interface().(*resourceA)
		assert.True(ok)
		assert.Equal(10, resourceReference.value)

		resourceCopy, ok := reflectedResource.Elem().Interface().(resourceA)
		assert.True(ok)
		assert.Equal(10, resourceCopy.value)
	})

	t.Run("resource can be mutated", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.add(&resourceA{value: 0})
		assert.NoError(err)

		reflectedResource, err := storage.getReflectResource(reflect.TypeFor[*resourceA]())
		assert.NoError(err)
		resource, ok := reflectedResource.Interface().(*resourceA)
		assert.True(ok)
		resource.value = 10

		reflectedResource, err = storage.getReflectResource(reflect.TypeFor[*resourceA]())
		assert.NoError(err)
		resource, ok = reflectedResource.Interface().(*resourceA)
		assert.True(ok)
		assert.Equal(10, resource.value)
	})
}

func TestResourceId(t *testing.T) {
	type resourceA struct{}
	type resourceB struct{}

	t.Run("resource return the same result, regardless of wether its passed by reference", func(t *testing.T) {
		assert := assert.New(t)

		a := reflectTypeToResourceId(reflect.TypeOf(resourceA{}))
		b := reflectTypeToResourceId(reflect.TypeOf(&resourceA{}))
		assert.Equal(a, b)
	})

	t.Run("resource ids for different resources are unique", func(t *testing.T) {
		assert := assert.New(t)

		a := reflectTypeToResourceId(reflect.TypeOf(resourceA{}))
		b := reflectTypeToResourceId(reflect.TypeOf(resourceB{}))
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
