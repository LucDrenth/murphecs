package app

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddToResourceStorage(t *testing.T) {
	type resourceA struct{}

	t.Run("fails to add resource if it already exists", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.add(&resourceA{})
		assert.NoError(err)
		err = storage.add(&resourceA{})
		assert.ErrorIs(err, ErrResourceAlreadyPresent)
	})

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

	t.Run("successfully gets resource copy", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.add(&resourceA{value: 10})
		assert.NoError(err)

		resource, err := getResourceFromStorage[resourceA](&storage)
		assert.NoError(err)
		assert.Equal(10, resource.value)
	})

	t.Run("resource copy can not be mutated", func(t *testing.T) {
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

	t.Run("successfully gets resource pointer", func(t *testing.T) {
		assert := assert.New(t)

		storage := newResourceStorage()
		err := storage.add(&resourceA{value: 10})
		assert.NoError(err)

		resource, err := getResourceFromStorage[*resourceA](&storage)
		assert.NoError(err)
		assert.Equal(10, resource.value)
	})

	t.Run("pointer resource can be mutated", func(t *testing.T) {
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

func TestComponentId(t *testing.T) {
	type resourceA struct{}
	type resourceB struct{}

	t.Run("resource return the same result, regardless of wether its passed by reference", func(t *testing.T) {
		assert := assert.New(t)

		a := reflectTypeToComponentId(reflect.TypeOf(resourceA{}))
		b := reflectTypeToComponentId(reflect.TypeOf(&resourceA{}))
		assert.Equal(a, b)
	})

	t.Run("resource ids for different resources are unique", func(t *testing.T) {
		assert := assert.New(t)

		a := reflectTypeToComponentId(reflect.TypeOf(resourceA{}))
		b := reflectTypeToComponentId(reflect.TypeOf(resourceB{}))
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
