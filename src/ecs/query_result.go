package ecs

// Due to yield only being able to return 2 params, it can not be implemented for queries that return
// more than 2 components.

type QueryResult interface {
	NumberOfResult() uint
	Clear()
}
type Query0Result struct {
	entityIds []EntityId
}
type Query1Result[A IComponent] struct {
	componentsA []*A
	entityIds   []EntityId
}
type Query2Result[A, B IComponent] struct {
	componentsA []*A
	componentsB []*B
	entityIds   []EntityId
}
type Query3Result[A, B, C IComponent] struct {
	componentsA []*A
	componentsB []*B
	componentsC []*C
	entityIds   []EntityId
}
type Query4Result[A, B, C, D IComponent] struct {
	componentsA []*A
	componentsB []*B
	componentsC []*C
	componentsD []*D
	entityIds   []EntityId
}

func (q *Query0Result) Clear() {
	q.entityIds = []EntityId{}
}
func (q *Query1Result[A]) Clear() {
	q.componentsA = []*A{}
	q.entityIds = []EntityId{}
}
func (q *Query2Result[A, B]) Clear() {
	q.componentsA = []*A{}
	q.componentsB = []*B{}
	q.entityIds = []EntityId{}
}
func (q *Query3Result[A, B, C]) Clear() {
	q.componentsA = []*A{}
	q.componentsB = []*B{}
	q.componentsC = []*C{}
	q.entityIds = []EntityId{}
}
func (q *Query4Result[A, B, C, D]) Clear() {
	q.componentsA = []*A{}
	q.componentsB = []*B{}
	q.componentsC = []*C{}
	q.componentsD = []*D{}
	q.entityIds = []EntityId{}
}

func (q *Query0Result) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query1Result[A]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query2Result[A, B]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query3Result[A, B, C]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query4Result[A, B, C, D]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}

// Range lets you range over the query result
//
// for component := range queryResult.Range() { ... }
func (q *Query1Result[A]) Range() func(yield func(*A) bool) {
	return func(yield func(*A) bool) {
		for i := range q.entityIds {
			if !yield(q.componentsA[i]) {
				return
			}
		}
	}
}

// Range lets you range over the query result
//
// for component := range queryResult.Range() { ... }
func (q *Query2Result[A, B]) Range() func(yield func(*A, *B) bool) {
	return func(yield func(*A, *B) bool) {
		for i := range q.entityIds {
			if !yield(q.componentsA[i], q.componentsB[i]) {
				return
			}
		}
	}
}

// Iter executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query0Result) Iter(f func(entityId EntityId) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i]); err != nil {
			return err
		}
	}

	return nil
}

// Iter executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query1Result[A]) Iter(f func(entityId EntityId, a *A) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i]); err != nil {
			return err
		}
	}

	return nil
}

// Iter executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query2Result[A, B]) Iter(f func(entityId EntityId, a *A, b *B) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i]); err != nil {
			return err
		}
	}

	return nil
}

// Iter executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query3Result[A, B, C]) Iter(f func(entityId EntityId, a *A, b *B, c *C) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i]); err != nil {
			return err
		}
	}

	return nil
}

// Iter executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query4Result[A, B, C, D]) Iter(f func(entityId EntityId, a *A, b *B, c *C, d *D) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i]); err != nil {
			return err
		}
	}

	return nil
}
