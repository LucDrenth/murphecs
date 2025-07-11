package ecs

import "fmt"

// Due to yield only being able to return 2 params, it can not be implemented for queries that return
// more than 2 components.

type QueryResult interface {
	NumberOfResult() uint

	// Clear removes the query results but reuses the allocated slices.
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
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query1Result[A]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query2Result[A, B]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query3Result[A, B, C]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query4Result[A, B, C, D]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.componentsD)
	q.componentsD = q.componentsD[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
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

// Iter executes function f on each entity that the query returned
func (q *Query0Result) Iter(f func(entityId EntityId)) {
	for i := range q.entityIds {
		f(q.entityIds[i])
	}
}

// Iter executes function f on each entity that the query returned
func (q *Query1Result[A]) Iter(f func(entityId EntityId, a *A)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i])
	}
}

// Iter executes function f on each entity that the query returned
func (q *Query2Result[A, B]) Iter(f func(entityId EntityId, a *A, b *B)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i])
	}
}

// Iter executes function f on each entity that the query returned
func (q *Query3Result[A, B, C]) Iter(f func(entityId EntityId, a *A, b *B, c *C)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i])
	}
}

// Iter executes function f on each entity that the query returned
func (q *Query4Result[A, B, C, D]) Iter(f func(entityId EntityId, a *A, b *B, c *C, d *D)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i])
	}
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query0Result) IterUntil(f func(entityId EntityId) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i]); err != nil {
			return err
		}
	}

	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query1Result[A]) IterUntil(f func(entityId EntityId, a *A) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i]); err != nil {
			return err
		}
	}

	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query2Result[A, B]) IterUntil(f func(entityId EntityId, a *A, b *B) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i]); err != nil {
			return err
		}
	}

	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query3Result[A, B, C]) IterUntil(f func(entityId EntityId, a *A, b *B, c *C) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i]); err != nil {
			return err
		}
	}

	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
func (q *Query4Result[A, B, C, D]) IterUntil(f func(entityId EntityId, a *A, b *B, c *C, d *D) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i]); err != nil {
			return err
		}
	}

	return nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query0Result) Single() (EntityId, error) {
	if q.NumberOfResult() != 1 {
		return nonExistingEntity, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query1Result[A]) Single() (EntityId, *A, error) {
	if q.NumberOfResult() != 1 {
		return nonExistingEntity, nil, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], q.componentsA[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query2Result[A, B]) Single() (EntityId, *A, *B, error) {
	if q.NumberOfResult() != 1 {
		return nonExistingEntity, nil, nil, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], q.componentsA[0], q.componentsB[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query3Result[A, B, C]) Single() (EntityId, *A, *B, *C, error) {
	if q.NumberOfResult() != 1 {
		return nonExistingEntity, nil, nil, nil, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query4Result[A, B, C, D]) Single() (EntityId, *A, *B, *C, *D, error) {
	if q.NumberOfResult() != 1 {
		return nonExistingEntity, nil, nil, nil, nil, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], nil
}
