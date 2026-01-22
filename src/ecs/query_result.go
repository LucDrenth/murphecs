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
	componentsA []A
	entityIds   []EntityId
}
type Query2Result[A, B IComponent] struct {
	componentsA []A
	componentsB []B
	entityIds   []EntityId
}
type Query3Result[A, B, C IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	entityIds   []EntityId
}
type Query4Result[A, B, C, D IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	componentsD []D
	entityIds   []EntityId
}
type Query5Result[A, B, C, D, E IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	componentsD []D
	componentsE []E
	entityIds   []EntityId
}
type Query6Result[A, B, C, D, E, F IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	componentsD []D
	componentsE []E
	componentsF []F
	entityIds   []EntityId
}
type Query7Result[A, B, C, D, E, F, G IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	componentsD []D
	componentsE []E
	componentsF []F
	componentsG []G
	entityIds   []EntityId
}
type Query8Result[A, B, C, D, E, F, G, H IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	componentsD []D
	componentsE []E
	componentsF []F
	componentsG []G
	componentsH []H
	entityIds   []EntityId
}
type Query9Result[A, B, C, D, E, F, G, H, I IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	componentsD []D
	componentsE []E
	componentsF []F
	componentsG []G
	componentsH []H
	componentsI []I
	entityIds   []EntityId
}
type Query10Result[A, B, C, D, E, F, G, H, I, J IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	componentsD []D
	componentsE []E
	componentsF []F
	componentsG []G
	componentsH []H
	componentsI []I
	componentsJ []J
	entityIds   []EntityId
}
type Query11Result[A, B, C, D, E, F, G, H, I, J, K IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	componentsD []D
	componentsE []E
	componentsF []F
	componentsG []G
	componentsH []H
	componentsI []I
	componentsJ []J
	componentsK []K
	entityIds   []EntityId
}
type Query12Result[A, B, C, D, E, F, G, H, I, J, K, L IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	componentsD []D
	componentsE []E
	componentsF []F
	componentsG []G
	componentsH []H
	componentsI []I
	componentsJ []J
	componentsK []K
	componentsL []L
	entityIds   []EntityId
}
type Query13Result[A, B, C, D, E, F, G, H, I, J, K, L, M IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	componentsD []D
	componentsE []E
	componentsF []F
	componentsG []G
	componentsH []H
	componentsI []I
	componentsJ []J
	componentsK []K
	componentsL []L
	componentsM []M
	entityIds   []EntityId
}
type Query14Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	componentsD []D
	componentsE []E
	componentsF []F
	componentsG []G
	componentsH []H
	componentsI []I
	componentsJ []J
	componentsK []K
	componentsL []L
	componentsM []M
	componentsN []N
	entityIds   []EntityId
}
type Query15Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	componentsD []D
	componentsE []E
	componentsF []F
	componentsG []G
	componentsH []H
	componentsI []I
	componentsJ []J
	componentsK []K
	componentsL []L
	componentsM []M
	componentsN []N
	componentsO []O
	entityIds   []EntityId
}
type Query16Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P IComponent] struct {
	componentsA []A
	componentsB []B
	componentsC []C
	componentsD []D
	componentsE []E
	componentsF []F
	componentsG []G
	componentsH []H
	componentsI []I
	componentsJ []J
	componentsK []K
	componentsL []L
	componentsM []M
	componentsN []N
	componentsO []O
	componentsP []P
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
func (q *Query5Result[A, B, C, D, E]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.componentsD)
	q.componentsD = q.componentsD[:0]
	clear(q.componentsE)
	q.componentsE = q.componentsE[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query6Result[A, B, C, D, E, F]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.componentsD)
	q.componentsD = q.componentsD[:0]
	clear(q.componentsE)
	q.componentsE = q.componentsE[:0]
	clear(q.componentsF)
	q.componentsF = q.componentsF[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query7Result[A, B, C, D, E, F, G]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.componentsD)
	q.componentsD = q.componentsD[:0]
	clear(q.componentsE)
	q.componentsE = q.componentsE[:0]
	clear(q.componentsF)
	q.componentsF = q.componentsF[:0]
	clear(q.componentsG)
	q.componentsG = q.componentsG[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query8Result[A, B, C, D, E, F, G, H]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.componentsD)
	q.componentsD = q.componentsD[:0]
	clear(q.componentsE)
	q.componentsE = q.componentsE[:0]
	clear(q.componentsF)
	q.componentsF = q.componentsF[:0]
	clear(q.componentsG)
	q.componentsG = q.componentsG[:0]
	clear(q.componentsH)
	q.componentsH = q.componentsH[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query9Result[A, B, C, D, E, F, G, H, I]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.componentsD)
	q.componentsD = q.componentsD[:0]
	clear(q.componentsE)
	q.componentsE = q.componentsE[:0]
	clear(q.componentsF)
	q.componentsF = q.componentsF[:0]
	clear(q.componentsG)
	q.componentsG = q.componentsG[:0]
	clear(q.componentsH)
	q.componentsH = q.componentsH[:0]
	clear(q.componentsI)
	q.componentsI = q.componentsI[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query10Result[A, B, C, D, E, F, G, H, I, J]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.componentsD)
	q.componentsD = q.componentsD[:0]
	clear(q.componentsE)
	q.componentsE = q.componentsE[:0]
	clear(q.componentsF)
	q.componentsF = q.componentsF[:0]
	clear(q.componentsG)
	q.componentsG = q.componentsG[:0]
	clear(q.componentsH)
	q.componentsH = q.componentsH[:0]
	clear(q.componentsI)
	q.componentsI = q.componentsI[:0]
	clear(q.componentsJ)
	q.componentsJ = q.componentsJ[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query11Result[A, B, C, D, E, F, G, H, I, J, K]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.componentsD)
	q.componentsD = q.componentsD[:0]
	clear(q.componentsE)
	q.componentsE = q.componentsE[:0]
	clear(q.componentsF)
	q.componentsF = q.componentsF[:0]
	clear(q.componentsG)
	q.componentsG = q.componentsG[:0]
	clear(q.componentsH)
	q.componentsH = q.componentsH[:0]
	clear(q.componentsI)
	q.componentsI = q.componentsI[:0]
	clear(q.componentsJ)
	q.componentsJ = q.componentsJ[:0]
	clear(q.componentsK)
	q.componentsK = q.componentsK[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query12Result[A, B, C, D, E, F, G, H, I, J, K, L]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.componentsD)
	q.componentsD = q.componentsD[:0]
	clear(q.componentsE)
	q.componentsE = q.componentsE[:0]
	clear(q.componentsF)
	q.componentsF = q.componentsF[:0]
	clear(q.componentsG)
	q.componentsG = q.componentsG[:0]
	clear(q.componentsH)
	q.componentsH = q.componentsH[:0]
	clear(q.componentsI)
	q.componentsI = q.componentsI[:0]
	clear(q.componentsJ)
	q.componentsJ = q.componentsJ[:0]
	clear(q.componentsK)
	q.componentsK = q.componentsK[:0]
	clear(q.componentsL)
	q.componentsL = q.componentsL[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query13Result[A, B, C, D, E, F, G, H, I, J, K, L, M]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.componentsD)
	q.componentsD = q.componentsD[:0]
	clear(q.componentsE)
	q.componentsE = q.componentsE[:0]
	clear(q.componentsF)
	q.componentsF = q.componentsF[:0]
	clear(q.componentsG)
	q.componentsG = q.componentsG[:0]
	clear(q.componentsH)
	q.componentsH = q.componentsH[:0]
	clear(q.componentsI)
	q.componentsI = q.componentsI[:0]
	clear(q.componentsJ)
	q.componentsJ = q.componentsJ[:0]
	clear(q.componentsK)
	q.componentsK = q.componentsK[:0]
	clear(q.componentsL)
	q.componentsL = q.componentsL[:0]
	clear(q.componentsM)
	q.componentsM = q.componentsM[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query14Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.componentsD)
	q.componentsD = q.componentsD[:0]
	clear(q.componentsE)
	q.componentsE = q.componentsE[:0]
	clear(q.componentsF)
	q.componentsF = q.componentsF[:0]
	clear(q.componentsG)
	q.componentsG = q.componentsG[:0]
	clear(q.componentsH)
	q.componentsH = q.componentsH[:0]
	clear(q.componentsI)
	q.componentsI = q.componentsI[:0]
	clear(q.componentsJ)
	q.componentsJ = q.componentsJ[:0]
	clear(q.componentsK)
	q.componentsK = q.componentsK[:0]
	clear(q.componentsL)
	q.componentsL = q.componentsL[:0]
	clear(q.componentsM)
	q.componentsM = q.componentsM[:0]
	clear(q.componentsN)
	q.componentsN = q.componentsN[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query15Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.componentsD)
	q.componentsD = q.componentsD[:0]
	clear(q.componentsE)
	q.componentsE = q.componentsE[:0]
	clear(q.componentsF)
	q.componentsF = q.componentsF[:0]
	clear(q.componentsG)
	q.componentsG = q.componentsG[:0]
	clear(q.componentsH)
	q.componentsH = q.componentsH[:0]
	clear(q.componentsI)
	q.componentsI = q.componentsI[:0]
	clear(q.componentsJ)
	q.componentsJ = q.componentsJ[:0]
	clear(q.componentsK)
	q.componentsK = q.componentsK[:0]
	clear(q.componentsL)
	q.componentsL = q.componentsL[:0]
	clear(q.componentsM)
	q.componentsM = q.componentsM[:0]
	clear(q.componentsN)
	q.componentsN = q.componentsN[:0]
	clear(q.componentsO)
	q.componentsO = q.componentsO[:0]
	clear(q.entityIds)
	q.entityIds = q.entityIds[:0]
}
func (q *Query16Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P]) Clear() {
	clear(q.componentsA)
	q.componentsA = q.componentsA[:0]
	clear(q.componentsB)
	q.componentsB = q.componentsB[:0]
	clear(q.componentsC)
	q.componentsC = q.componentsC[:0]
	clear(q.componentsD)
	q.componentsD = q.componentsD[:0]
	clear(q.componentsE)
	q.componentsE = q.componentsE[:0]
	clear(q.componentsF)
	q.componentsF = q.componentsF[:0]
	clear(q.componentsG)
	q.componentsG = q.componentsG[:0]
	clear(q.componentsH)
	q.componentsH = q.componentsH[:0]
	clear(q.componentsI)
	q.componentsI = q.componentsI[:0]
	clear(q.componentsJ)
	q.componentsJ = q.componentsJ[:0]
	clear(q.componentsK)
	q.componentsK = q.componentsK[:0]
	clear(q.componentsL)
	q.componentsL = q.componentsL[:0]
	clear(q.componentsM)
	q.componentsM = q.componentsM[:0]
	clear(q.componentsN)
	q.componentsN = q.componentsN[:0]
	clear(q.componentsO)
	q.componentsO = q.componentsO[:0]
	clear(q.componentsP)
	q.componentsP = q.componentsP[:0]
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
func (q *Query5Result[A, B, C, D, E]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query6Result[A, B, C, D, E, F]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query7Result[A, B, C, D, E, F, G]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query8Result[A, B, C, D, E, F, G, H]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query9Result[A, B, C, D, E, F, G, H, I]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query10Result[A, B, C, D, E, F, G, H, I, J]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query11Result[A, B, C, D, E, F, G, H, I, J, K]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query12Result[A, B, C, D, E, F, G, H, I, J, K, L]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query13Result[A, B, C, D, E, F, G, H, I, J, K, L, M]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query14Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query15Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}
func (q *Query16Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P]) NumberOfResult() uint {
	return uint(len(q.entityIds))
}

// Range lets you range over the query result
//
// for component := range queryResult.Range() { ... }
func (q *Query1Result[A]) Range() func(yield func(A) bool) {
	return func(yield func(A) bool) {
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
func (q *Query2Result[A, B]) Range() func(yield func(A, B) bool) {
	return func(yield func(A, B) bool) {
		for i := range q.entityIds {
			if !yield(q.componentsA[i], q.componentsB[i]) {
				return
			}
		}
	}
}

// Iter executes function f on each entity that the query returned.
func (q *Query0Result) Iter(f func(entityId EntityId)) {
	for i := range q.entityIds {
		f(q.entityIds[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointer.
func (q *Query1Result[A]) Iter(f func(entityId EntityId, a A)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query2Result[A, B]) Iter(f func(entityId EntityId, a A, b B)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query3Result[A, B, C]) Iter(f func(entityId EntityId, a A, b B, c C)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query4Result[A, B, C, D]) Iter(f func(entityId EntityId, a A, b B, c C, d D)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query5Result[A, B, C, D, E]) Iter(f func(entityId EntityId, a A, b B, c C, d D, e E)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query6Result[A, B, C, D, E, F]) Iter(f func(entityId EntityId, a A, b B, c C, d D, e E, f F)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query7Result[A, B, C, D, E, F, G]) Iter(f func(entityId EntityId, a A, b B, c C, d D, e E, f F, g G)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query8Result[A, B, C, D, E, F, G, H]) Iter(f func(entityId EntityId, a A, b B, c C, d D, e E, f F, g G, h H)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query9Result[A, B, C, D, E, F, G, H, I]) Iter(f func(EntityId, A, B, C, D, E, F, G, H, I)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query10Result[A, B, C, D, E, F, G, H, I, J]) Iter(f func(EntityId, A, B, C, D, E, F, G, H, I, J)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query11Result[A, B, C, D, E, F, G, H, I, J, K]) Iter(f func(EntityId, A, B, C, D, E, F, G, H, I, J, K)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i], q.componentsK[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query12Result[A, B, C, D, E, F, G, H, I, J, K, L]) Iter(f func(EntityId, A, B, C, D, E, F, G, H, I, J, K, L)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i], q.componentsK[i], q.componentsL[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query13Result[A, B, C, D, E, F, G, H, I, J, K, L, M]) Iter(f func(EntityId, A, B, C, D, E, F, G, H, I, J, K, L, M)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i], q.componentsK[i], q.componentsL[i], q.componentsM[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query14Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N]) Iter(f func(EntityId, A, B, C, D, E, F, G, H, I, J, K, L, M, N)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i], q.componentsK[i], q.componentsL[i], q.componentsM[i], q.componentsN[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query15Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O]) Iter(f func(EntityId, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i], q.componentsK[i], q.componentsL[i], q.componentsM[i], q.componentsN[i], q.componentsO[i])
	}
}

// Iter executes function f on each entity that the query returned.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query16Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P]) Iter(f func(EntityId, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P)) {
	for i := range q.entityIds {
		f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i], q.componentsK[i], q.componentsL[i], q.componentsM[i], q.componentsN[i], q.componentsO[i], q.componentsP[i])
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
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointer.
func (q *Query1Result[A]) IterUntil(f func(entityId EntityId, a A) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i]); err != nil {
			return err
		}
	}

	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query2Result[A, B]) IterUntil(f func(entityId EntityId, a A, b B) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i]); err != nil {
			return err
		}
	}

	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query3Result[A, B, C]) IterUntil(f func(entityId EntityId, a A, b B, c C) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i]); err != nil {
			return err
		}
	}

	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query4Result[A, B, C, D]) IterUntil(f func(entityId EntityId, a A, b B, c C, d D) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i]); err != nil {
			return err
		}
	}

	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query5Result[A, B, C, D, E]) IterUntil(f func(entityId EntityId, a A, b B, c C, d D, e E) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i]); err != nil {
			return err
		}
	}

	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query6Result[A, B, C, D, E, F]) IterUntil(f func(entityId EntityId, a A, b B, c C, d D, e E, f F) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i]); err != nil {
			return err
		}
	}

	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query7Result[A, B, C, D, E, F, G]) IterUntil(f func(entityId EntityId, a A, b B, c C, d D, e E, f F, g G) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i]); err != nil {
			return err
		}
	}

	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query8Result[A, B, C, D, E, F, G, H]) IterUntil(f func(entityId EntityId, a A, b B, c C, d D, e E, f F, g G, h H) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i]); err != nil {
			return err
		}
	}

	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query9Result[A, B, C, D, E, F, G, H, I]) IterUntil(f func(EntityId, A, B, C, D, E, F, G, H, I) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i]); err != nil {
			return err
		}
	}
	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query10Result[A, B, C, D, E, F, G, H, I, J]) IterUntil(f func(EntityId, A, B, C, D, E, F, G, H, I, J) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i]); err != nil {
			return err
		}
	}
	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query11Result[A, B, C, D, E, F, G, H, I, J, K]) IterUntil(f func(EntityId, A, B, C, D, E, F, G, H, I, J, K) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i], q.componentsK[i]); err != nil {
			return err
		}
	}
	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query12Result[A, B, C, D, E, F, G, H, I, J, K, L]) IterUntil(f func(EntityId, A, B, C, D, E, F, G, H, I, J, K, L) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i], q.componentsK[i], q.componentsL[i]); err != nil {
			return err
		}
	}
	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query13Result[A, B, C, D, E, F, G, H, I, J, K, L, M]) IterUntil(f func(EntityId, A, B, C, D, E, F, G, H, I, J, K, L, M) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i], q.componentsK[i], q.componentsL[i], q.componentsM[i]); err != nil {
			return err
		}
	}
	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query14Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N]) IterUntil(f func(EntityId, A, B, C, D, E, F, G, H, I, J, K, L, M, N) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i], q.componentsK[i], q.componentsL[i], q.componentsM[i], q.componentsN[i]); err != nil {
			return err
		}
	}
	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query15Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O]) IterUntil(f func(EntityId, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i], q.componentsK[i], q.componentsL[i], q.componentsM[i], q.componentsN[i], q.componentsO[i]); err != nil {
			return err
		}
	}
	return nil
}

// IterUntil executes function f on each entity that the query returned, until f returns an error.
// If any of the calls to f returned an error, this function returns that error.
//
// During iteration, do not make changes that update the archetype
// of the entity. This invalidates the component pointers.
func (q *Query16Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P]) IterUntil(f func(EntityId, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P) error) error {
	for i := range q.entityIds {
		if err := f(q.entityIds[i], q.componentsA[i], q.componentsB[i], q.componentsC[i], q.componentsD[i], q.componentsE[i], q.componentsF[i], q.componentsG[i], q.componentsH[i], q.componentsI[i], q.componentsJ[i], q.componentsK[i], q.componentsL[i], q.componentsM[i], q.componentsN[i], q.componentsO[i], q.componentsP[i]); err != nil {
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
func (q *Query1Result[A]) Single() (EntityId, A, error) {
	if q.NumberOfResult() != 1 {
		var a A
		return nonExistingEntity, a, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], q.componentsA[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query2Result[A, B]) Single() (EntityId, A, B, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		return nonExistingEntity, a, b, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], q.componentsA[0], q.componentsB[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query3Result[A, B, C]) Single() (EntityId, A, B, C, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		return nonExistingEntity, a, b, c, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query4Result[A, B, C, D]) Single() (EntityId, A, B, C, D, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		var d D
		return nonExistingEntity, a, b, c, d, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query5Result[A, B, C, D, E]) Single() (EntityId, A, B, C, D, E, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		var d D
		var e E
		return nonExistingEntity, a, b, c, d, e, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], q.componentsE[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query6Result[A, B, C, D, E, F]) Single() (EntityId, A, B, C, D, E, F, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		var d D
		var e E
		var f F
		return nonExistingEntity, a, b, c, d, e, f, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], q.componentsE[0], q.componentsF[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query7Result[A, B, C, D, E, F, G]) Single() (EntityId, A, B, C, D, E, F, G, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		var d D
		var e E
		var f F
		var g G
		return nonExistingEntity, a, b, c, d, e, f, g, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], q.componentsE[0], q.componentsF[0], q.componentsG[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query8Result[A, B, C, D, E, F, G, H]) Single() (EntityId, A, B, C, D, E, F, G, H, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		var d D
		var e E
		var f F
		var g G
		var h H
		return nonExistingEntity, a, b, c, d, e, f, g, h, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}

	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], q.componentsE[0], q.componentsF[0], q.componentsG[0], q.componentsH[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query9Result[A, B, C, D, E, F, G, H, I]) Single() (EntityId, A, B, C, D, E, F, G, H, I, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		var d D
		var e E
		var f F
		var g G
		var h H
		var i I
		return nonExistingEntity, a, b, c, d, e, f, g, h, i, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}
	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], q.componentsE[0], q.componentsF[0], q.componentsG[0], q.componentsH[0], q.componentsI[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query10Result[A, B, C, D, E, F, G, H, I, J]) Single() (EntityId, A, B, C, D, E, F, G, H, I, J, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		var d D
		var e E
		var f F
		var g G
		var h H
		var i I
		var j J
		return nonExistingEntity, a, b, c, d, e, f, g, h, i, j, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}
	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], q.componentsE[0], q.componentsF[0], q.componentsG[0], q.componentsH[0], q.componentsI[0], q.componentsJ[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query11Result[A, B, C, D, E, F, G, H, I, J, K]) Single() (EntityId, A, B, C, D, E, F, G, H, I, J, K, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		var d D
		var e E
		var f F
		var g G
		var h H
		var i I
		var j J
		var k K
		return nonExistingEntity, a, b, c, d, e, f, g, h, i, j, k, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}
	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], q.componentsE[0], q.componentsF[0], q.componentsG[0], q.componentsH[0], q.componentsI[0], q.componentsJ[0], q.componentsK[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query12Result[A, B, C, D, E, F, G, H, I, J, K, L]) Single() (EntityId, A, B, C, D, E, F, G, H, I, J, K, L, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		var d D
		var e E
		var f F
		var g G
		var h H
		var i I
		var j J
		var k K
		var l L
		return nonExistingEntity, a, b, c, d, e, f, g, h, i, j, k, l, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}
	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], q.componentsE[0], q.componentsF[0], q.componentsG[0], q.componentsH[0], q.componentsI[0], q.componentsJ[0], q.componentsK[0], q.componentsL[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query13Result[A, B, C, D, E, F, G, H, I, J, K, L, M]) Single() (EntityId, A, B, C, D, E, F, G, H, I, J, K, L, M, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		var d D
		var e E
		var f F
		var g G
		var h H
		var i I
		var j J
		var k K
		var l L
		var m M
		return nonExistingEntity, a, b, c, d, e, f, g, h, i, j, k, l, m, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}
	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], q.componentsE[0], q.componentsF[0], q.componentsG[0], q.componentsH[0], q.componentsI[0], q.componentsJ[0], q.componentsK[0], q.componentsL[0], q.componentsM[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query14Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N]) Single() (EntityId, A, B, C, D, E, F, G, H, I, J, K, L, M, N, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		var d D
		var e E
		var f F
		var g G
		var h H
		var i I
		var j J
		var k K
		var l L
		var m M
		var n N
		return nonExistingEntity, a, b, c, d, e, f, g, h, i, j, k, l, m, n, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}
	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], q.componentsE[0], q.componentsF[0], q.componentsG[0], q.componentsH[0], q.componentsI[0], q.componentsJ[0], q.componentsK[0], q.componentsL[0], q.componentsM[0], q.componentsN[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query15Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O]) Single() (EntityId, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		var d D
		var e E
		var f F
		var g G
		var h H
		var i I
		var j J
		var k K
		var l L
		var m M
		var n N
		var o O
		return nonExistingEntity, a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}
	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], q.componentsE[0], q.componentsF[0], q.componentsG[0], q.componentsH[0], q.componentsI[0], q.componentsJ[0], q.componentsK[0], q.componentsL[0], q.componentsM[0], q.componentsN[0], q.componentsO[0], nil
}

// Single returns the only query result, or an [ErrUnexpectedNumberOfQueryResults] error if there is not exactly 1 result.
func (q *Query16Result[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P]) Single() (EntityId, A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, error) {
	if q.NumberOfResult() != 1 {
		var a A
		var b B
		var c C
		var d D
		var e E
		var f F
		var g G
		var h H
		var i I
		var j J
		var k K
		var l L
		var m M
		var n N
		var o O
		var p P
		return nonExistingEntity, a, b, c, d, e, f, g, h, i, j, k, l, m, n, o, p, fmt.Errorf("%w: expected 1, got %d", ErrUnexpectedNumberOfQueryResults, q.NumberOfResult())
	}
	return q.entityIds[0], q.componentsA[0], q.componentsB[0], q.componentsC[0], q.componentsD[0], q.componentsE[0], q.componentsF[0], q.componentsG[0], q.componentsH[0], q.componentsI[0], q.componentsJ[0], q.componentsK[0], q.componentsL[0], q.componentsM[0], q.componentsN[0], q.componentsO[0], q.componentsP[0], nil
}
