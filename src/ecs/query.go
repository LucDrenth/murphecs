package ecs

import (
	"fmt"
	"reflect"
	"slices"
)

type Query interface {
	Exec(world *World) error

	// Prepare needs to be called before running Exec.
	//
	// This step is necessary because the way that queries are created is optimized for users, and not
	// for the computer. This method closes that gap.
	//
	// Pass the otherWorld parameter if the Query targets another world.
	Prepare(defaultWorld *World, otherWorlds *map[WorldId]*World) error

	// Validate checks if there are any unexpected or unoptimized things. It returns an error if there
	// is something that can be optimized. The error should be treated as a warning, and does not mean
	// that the query can not be executed.
	//
	// Prepare must be called without returning any errors before calling this method.
	Validate() error

	// Clear the query results that got filled when last running Exec.
	ClearResults()

	// IsLazy returns wether this query should be treated as lazy or not. Being lazy means that it should
	// not get executed automatically, and should be done by the user.
	IsLazy() bool

	// TargetWorld returns wether this query should be executed in a custom world. Returns nil if no custom
	// world should be used, in which case it defaults to the world of the SubApp it is used in.
	TargetWorld() *WorldId

	getOptions() *CombinedQueryOptions
}

type queryComponentInfo struct {
	id        ComponentId
	isPointer bool
}

func queryComponentInfoFor[T IComponent](world *World) queryComponentInfo {
	return queryComponentInfo{
		id:        ComponentIdFor[T](world),
		isPointer: reflect.TypeFor[T]().Kind() == reflect.Pointer,
	}
}

type queryOptions struct {
	options    CombinedQueryOptions
	components []ComponentId
}

func (o *queryOptions) getOptions() *CombinedQueryOptions {
	return &o.options
}

func (o *queryOptions) IsLazy() bool {
	return o.options.isLazy
}

func (o *queryOptions) TargetWorld() *WorldId {
	return o.options.TargetWorld
}

func (o *queryOptions) Validate() error {
	return o.options.validateOptions(o.components)
}

// Query0 gets the entities that match the query options.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query0[_ QueryOption] struct {
	Query0Result
	queryOptions
}

// Query1 queries 1 component.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1] to mark component as optional. When a component is optional, entities do not
//     have to have that component in order for it to return a result, as long as it has the other (not-optional)
//     components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query1[ComponentA IComponent, _ QueryOption] struct {
	Query1Result[ComponentA]
	queryOptions

	componentInfoA queryComponentInfo
}

// Query2 queries 2 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] to mark components as optional. When a component is optional, entities
//     do not have to have that component in order for it to return a result, as long as it has the other (not-optional)
//     components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query2[ComponentA, ComponentB IComponent, _ QueryOption] struct {
	Query2Result[ComponentA, ComponentB]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
}

// Query3 queries 3 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query3[ComponentA, ComponentB, ComponentC IComponent, _ QueryOption] struct {
	Query3Result[ComponentA, ComponentB, ComponentC]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
}

// Query4 queries 4 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query4[ComponentA, ComponentB, ComponentC, ComponentD IComponent, _ QueryOption] struct {
	Query4Result[ComponentA, ComponentB, ComponentC, ComponentD]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
	componentInfoD queryComponentInfo
}

// Query5 queries 5 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query5[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE IComponent, _ QueryOption] struct {
	Query5Result[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
	componentInfoD queryComponentInfo
	componentInfoE queryComponentInfo
}

// Query6 queries 6 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query6[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF IComponent, _ QueryOption] struct {
	Query6Result[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
	componentInfoD queryComponentInfo
	componentInfoE queryComponentInfo
	componentInfoF queryComponentInfo
}

// Query7 queries 7 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query7[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG IComponent, _ QueryOption] struct {
	Query7Result[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
	componentInfoD queryComponentInfo
	componentInfoE queryComponentInfo
	componentInfoF queryComponentInfo
	componentInfoG queryComponentInfo
}

// Query8 queries 8 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query8[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH IComponent, _ QueryOption] struct {
	Query8Result[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
	componentInfoD queryComponentInfo
	componentInfoE queryComponentInfo
	componentInfoF queryComponentInfo
	componentInfoG queryComponentInfo
	componentInfoH queryComponentInfo
}

// Query9 queries 9 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query9[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI IComponent, _ QueryOption] struct {
	Query9Result[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
	componentInfoD queryComponentInfo
	componentInfoE queryComponentInfo
	componentInfoF queryComponentInfo
	componentInfoG queryComponentInfo
	componentInfoH queryComponentInfo
	componentInfoI queryComponentInfo
}

// Query10 queries 10 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query10[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ IComponent, _ QueryOption] struct {
	Query10Result[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
	componentInfoD queryComponentInfo
	componentInfoE queryComponentInfo
	componentInfoF queryComponentInfo
	componentInfoG queryComponentInfo
	componentInfoH queryComponentInfo
	componentInfoI queryComponentInfo
	componentInfoJ queryComponentInfo
}

// Query11 queries 11 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query11[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ, ComponentK IComponent, _ QueryOption] struct {
	Query11Result[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ, ComponentK]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
	componentInfoD queryComponentInfo
	componentInfoE queryComponentInfo
	componentInfoF queryComponentInfo
	componentInfoG queryComponentInfo
	componentInfoH queryComponentInfo
	componentInfoI queryComponentInfo
	componentInfoJ queryComponentInfo
	componentInfoK queryComponentInfo
}

// Query12 queries 12 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query12[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ, ComponentK, ComponentL IComponent, _ QueryOption] struct {
	Query12Result[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ, ComponentK, ComponentL]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
	componentInfoD queryComponentInfo
	componentInfoE queryComponentInfo
	componentInfoF queryComponentInfo
	componentInfoG queryComponentInfo
	componentInfoH queryComponentInfo
	componentInfoI queryComponentInfo
	componentInfoJ queryComponentInfo
	componentInfoK queryComponentInfo
	componentInfoL queryComponentInfo
}

// Query13 queries 13 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query13[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ, ComponentK, ComponentL, ComponentM IComponent, _ QueryOption] struct {
	Query13Result[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ, ComponentK, ComponentL, ComponentM]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
	componentInfoD queryComponentInfo
	componentInfoE queryComponentInfo
	componentInfoF queryComponentInfo
	componentInfoG queryComponentInfo
	componentInfoH queryComponentInfo
	componentInfoI queryComponentInfo
	componentInfoJ queryComponentInfo
	componentInfoK queryComponentInfo
	componentInfoL queryComponentInfo
	componentInfoM queryComponentInfo
}

// Query14 queries 14 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query14[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ, ComponentK, ComponentL, ComponentM, ComponentN IComponent, _ QueryOption] struct {
	Query14Result[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ, ComponentK, ComponentL, ComponentM, ComponentN]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
	componentInfoD queryComponentInfo
	componentInfoE queryComponentInfo
	componentInfoF queryComponentInfo
	componentInfoG queryComponentInfo
	componentInfoH queryComponentInfo
	componentInfoI queryComponentInfo
	componentInfoJ queryComponentInfo
	componentInfoK queryComponentInfo
	componentInfoL queryComponentInfo
	componentInfoM queryComponentInfo
	componentInfoN queryComponentInfo
}

// Query15 queries 15 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query15[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ, ComponentK, ComponentL, ComponentM, ComponentN, ComponentO IComponent, _ QueryOption] struct {
	Query15Result[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ, ComponentK, ComponentL, ComponentM, ComponentN, ComponentO]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
	componentInfoD queryComponentInfo
	componentInfoE queryComponentInfo
	componentInfoF queryComponentInfo
	componentInfoG queryComponentInfo
	componentInfoH queryComponentInfo
	componentInfoI queryComponentInfo
	componentInfoJ queryComponentInfo
	componentInfoK queryComponentInfo
	componentInfoL queryComponentInfo
	componentInfoM queryComponentInfo
	componentInfoN queryComponentInfo
	componentInfoO queryComponentInfo
}

// Query16 queries 16 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoOptional], [Optional1], [Optional2] (and so on) to mark components as optional. When a
//     component is optional, entities do not have to have that component in order for it to return a result,
//     as long as it has the other (not-optional) components.
//   - use [NoFilter] to not use any filters
//   - use [With] to make the results only include entities that has a specific component.
//   - use [Without] to make the results only include entities that do not have a specific component.
//   - use [And] and [Or] to combine filters.
type Query16[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ, ComponentK, ComponentL, ComponentM, ComponentN, ComponentO, ComponentP IComponent, _ QueryOption] struct {
	Query16Result[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, ComponentI, ComponentJ, ComponentK, ComponentL, ComponentM, ComponentN, ComponentO, ComponentP]
	queryOptions

	componentInfoA queryComponentInfo
	componentInfoB queryComponentInfo
	componentInfoC queryComponentInfo
	componentInfoD queryComponentInfo
	componentInfoE queryComponentInfo
	componentInfoF queryComponentInfo
	componentInfoG queryComponentInfo
	componentInfoH queryComponentInfo
	componentInfoI queryComponentInfo
	componentInfoJ queryComponentInfo
	componentInfoK queryComponentInfo
	componentInfoL queryComponentInfo
	componentInfoM queryComponentInfo
	componentInfoN queryComponentInfo
	componentInfoO queryComponentInfo
	componentInfoP queryComponentInfo
}

func (q *Query0[QueryOptions]) Exec(world *World) error {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		q.entityIds = append(q.entityIds, archetype.entities...)
	}

	return nil
}

func (q *Query1[ComponentA, QueryOptions]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}
func (q *Query2[ComponentA, ComponentB, QueryOptions]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}
func (q *Query3[ComponentA, ComponentB, ComponentC, QueryOptions]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c ComponentC
			if fetchC {
				c, err = fetchComponentForQueryResult[ComponentC](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}
func (q *Query4[ComponentA, ComponentB, ComponentC, ComponentD, QueryOptions]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentInfoD.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c ComponentC
			if fetchC {
				c, err = fetchComponentForQueryResult[ComponentC](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var d ComponentD
			if fetchD {
				d, err = fetchComponentForQueryResult[ComponentD](q.componentInfoD, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.componentsD = append(q.componentsD, d)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}
func (q *Query5[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, QueryOptions]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentInfoD.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentInfoE.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c ComponentC
			if fetchC {
				c, err = fetchComponentForQueryResult[ComponentC](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var d ComponentD
			if fetchD {
				d, err = fetchComponentForQueryResult[ComponentD](q.componentInfoD, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var e ComponentE
			if fetchE {
				e, err = fetchComponentForQueryResult[ComponentE](q.componentInfoE, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.componentsD = append(q.componentsD, d)
			q.componentsE = append(q.componentsE, e)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}
func (q *Query6[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, QueryOptions]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentInfoD.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentInfoE.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentInfoF.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c ComponentC
			if fetchC {
				c, err = fetchComponentForQueryResult[ComponentC](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var d ComponentD
			if fetchD {
				d, err = fetchComponentForQueryResult[ComponentD](q.componentInfoD, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var e ComponentE
			if fetchE {
				e, err = fetchComponentForQueryResult[ComponentE](q.componentInfoE, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var f ComponentF
			if fetchF {
				f, err = fetchComponentForQueryResult[ComponentF](q.componentInfoF, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.componentsD = append(q.componentsD, d)
			q.componentsE = append(q.componentsE, e)
			q.componentsF = append(q.componentsF, f)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}
func (q *Query7[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, QueryOptions]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentInfoD.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentInfoE.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentInfoF.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentInfoG.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c ComponentC
			if fetchC {
				c, err = fetchComponentForQueryResult[ComponentC](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var d ComponentD
			if fetchD {
				d, err = fetchComponentForQueryResult[ComponentD](q.componentInfoD, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var e ComponentE
			if fetchE {
				e, err = fetchComponentForQueryResult[ComponentE](q.componentInfoE, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var f ComponentF
			if fetchF {
				f, err = fetchComponentForQueryResult[ComponentF](q.componentInfoF, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var g ComponentG
			if fetchG {
				g, err = fetchComponentForQueryResult[ComponentG](q.componentInfoG, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.componentsD = append(q.componentsD, d)
			q.componentsE = append(q.componentsE, e)
			q.componentsF = append(q.componentsF, f)
			q.componentsG = append(q.componentsG, g)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}
func (q *Query8[ComponentA, ComponentB, ComponentC, ComponentD, ComponentE, ComponentF, ComponentG, ComponentH, QueryOptions]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentInfoD.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentInfoE.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentInfoF.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentInfoG.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentInfoH.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c ComponentC
			if fetchC {
				c, err = fetchComponentForQueryResult[ComponentC](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var d ComponentD
			if fetchD {
				d, err = fetchComponentForQueryResult[ComponentD](q.componentInfoD, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var e ComponentE
			if fetchE {
				e, err = fetchComponentForQueryResult[ComponentE](q.componentInfoE, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var f ComponentF
			if fetchF {
				f, err = fetchComponentForQueryResult[ComponentF](q.componentInfoF, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var g ComponentG
			if fetchG {
				g, err = fetchComponentForQueryResult[ComponentG](q.componentInfoG, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var h ComponentH
			if fetchH {
				h, err = fetchComponentForQueryResult[ComponentH](q.componentInfoH, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.componentsD = append(q.componentsD, d)
			q.componentsE = append(q.componentsE, e)
			q.componentsF = append(q.componentsF, f)
			q.componentsG = append(q.componentsG, g)
			q.componentsH = append(q.componentsH, h)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}

func (q *Query9[A, B, C, D, E, F, G, H, I, Options]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentInfoD.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentInfoE.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentInfoF.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentInfoG.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentInfoH.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentInfoI.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var d D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentInfoD, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var e E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentInfoE, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var f F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentInfoF, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var g G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentInfoG, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var h H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentInfoH, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var i I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentInfoI, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.componentsD = append(q.componentsD, d)
			q.componentsE = append(q.componentsE, e)
			q.componentsF = append(q.componentsF, f)
			q.componentsG = append(q.componentsG, g)
			q.componentsH = append(q.componentsH, h)
			q.componentsI = append(q.componentsI, i)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}

func (q *Query10[A, B, C, D, E, F, G, H, I, J, Options]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentInfoD.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentInfoE.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentInfoF.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentInfoG.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentInfoH.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentInfoI.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentInfoJ.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var d D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentInfoD, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var e E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentInfoE, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var f F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentInfoF, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var g G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentInfoG, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var h H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentInfoH, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var i I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentInfoI, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var j J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentInfoJ, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.componentsD = append(q.componentsD, d)
			q.componentsE = append(q.componentsE, e)
			q.componentsF = append(q.componentsF, f)
			q.componentsG = append(q.componentsG, g)
			q.componentsH = append(q.componentsH, h)
			q.componentsI = append(q.componentsI, i)
			q.componentsJ = append(q.componentsJ, j)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}

func (q *Query11[A, B, C, D, E, F, G, H, I, J, K, Options]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentInfoD.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentInfoE.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentInfoF.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentInfoG.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentInfoH.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentInfoI.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentInfoJ.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchK, skip := shouldHandleQueryComponent(q.componentInfoK.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var d D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentInfoD, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var e E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentInfoE, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var f F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentInfoF, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var g G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentInfoG, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var h H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentInfoH, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var i I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentInfoI, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var j J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentInfoJ, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var k K
			if fetchK {
				k, err = fetchComponentForQueryResult[K](q.componentInfoK, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.componentsD = append(q.componentsD, d)
			q.componentsE = append(q.componentsE, e)
			q.componentsF = append(q.componentsF, f)
			q.componentsG = append(q.componentsG, g)
			q.componentsH = append(q.componentsH, h)
			q.componentsI = append(q.componentsI, i)
			q.componentsJ = append(q.componentsJ, j)
			q.componentsK = append(q.componentsK, k)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}

func (q *Query12[A, B, C, D, E, F, G, H, I, J, K, L, Options]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentInfoD.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentInfoE.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentInfoF.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentInfoG.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentInfoH.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentInfoI.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentInfoJ.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchK, skip := shouldHandleQueryComponent(q.componentInfoK.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchL, skip := shouldHandleQueryComponent(q.componentInfoL.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var d D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentInfoD, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var e E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentInfoE, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var f F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentInfoF, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var g G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentInfoG, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var h H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentInfoH, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var i I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentInfoI, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var j J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentInfoJ, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var k K
			if fetchK {
				k, err = fetchComponentForQueryResult[K](q.componentInfoK, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var l L
			if fetchL {
				l, err = fetchComponentForQueryResult[L](q.componentInfoL, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.componentsD = append(q.componentsD, d)
			q.componentsE = append(q.componentsE, e)
			q.componentsF = append(q.componentsF, f)
			q.componentsG = append(q.componentsG, g)
			q.componentsH = append(q.componentsH, h)
			q.componentsI = append(q.componentsI, i)
			q.componentsJ = append(q.componentsJ, j)
			q.componentsK = append(q.componentsK, k)
			q.componentsL = append(q.componentsL, l)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}

func (q *Query13[A, B, C, D, E, F, G, H, I, J, K, L, M, Options]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentInfoD.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentInfoE.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentInfoF.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentInfoG.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentInfoH.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentInfoI.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentInfoJ.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchK, skip := shouldHandleQueryComponent(q.componentInfoK.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchL, skip := shouldHandleQueryComponent(q.componentInfoL.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchM, skip := shouldHandleQueryComponent(q.componentInfoM.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var d D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentInfoD, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var e E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentInfoE, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var f F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentInfoF, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var g G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentInfoG, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var h H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentInfoH, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var i I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentInfoI, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var j J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentInfoJ, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var k K
			if fetchK {
				k, err = fetchComponentForQueryResult[K](q.componentInfoK, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var l L
			if fetchL {
				l, err = fetchComponentForQueryResult[L](q.componentInfoL, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var m M
			if fetchM {
				m, err = fetchComponentForQueryResult[M](q.componentInfoM, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.componentsD = append(q.componentsD, d)
			q.componentsE = append(q.componentsE, e)
			q.componentsF = append(q.componentsF, f)
			q.componentsG = append(q.componentsG, g)
			q.componentsH = append(q.componentsH, h)
			q.componentsI = append(q.componentsI, i)
			q.componentsJ = append(q.componentsJ, j)
			q.componentsK = append(q.componentsK, k)
			q.componentsL = append(q.componentsL, l)
			q.componentsM = append(q.componentsM, m)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}

func (q *Query14[A, B, C, D, E, F, G, H, I, J, K, L, M, N, Options]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentInfoD.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentInfoE.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentInfoF.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentInfoG.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentInfoH.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentInfoI.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentInfoJ.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchK, skip := shouldHandleQueryComponent(q.componentInfoK.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchL, skip := shouldHandleQueryComponent(q.componentInfoL.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchM, skip := shouldHandleQueryComponent(q.componentInfoM.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchN, skip := shouldHandleQueryComponent(q.componentInfoN.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var d D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentInfoD, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var e E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentInfoE, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var f F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentInfoF, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var g G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentInfoG, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var h H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentInfoH, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var i I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentInfoI, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var j J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentInfoJ, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var k K
			if fetchK {
				k, err = fetchComponentForQueryResult[K](q.componentInfoK, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var l L
			if fetchL {
				l, err = fetchComponentForQueryResult[L](q.componentInfoL, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var m M
			if fetchM {
				m, err = fetchComponentForQueryResult[M](q.componentInfoM, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var n N
			if fetchN {
				n, err = fetchComponentForQueryResult[N](q.componentInfoN, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.componentsD = append(q.componentsD, d)
			q.componentsE = append(q.componentsE, e)
			q.componentsF = append(q.componentsF, f)
			q.componentsG = append(q.componentsG, g)
			q.componentsH = append(q.componentsH, h)
			q.componentsI = append(q.componentsI, i)
			q.componentsJ = append(q.componentsJ, j)
			q.componentsK = append(q.componentsK, k)
			q.componentsL = append(q.componentsL, l)
			q.componentsM = append(q.componentsM, m)
			q.componentsN = append(q.componentsN, n)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}

func (q *Query15[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, Options]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentInfoD.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentInfoE.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentInfoF.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentInfoG.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentInfoH.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentInfoI.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentInfoJ.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchK, skip := shouldHandleQueryComponent(q.componentInfoK.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchL, skip := shouldHandleQueryComponent(q.componentInfoL.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchM, skip := shouldHandleQueryComponent(q.componentInfoM.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchN, skip := shouldHandleQueryComponent(q.componentInfoN.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchO, skip := shouldHandleQueryComponent(q.componentInfoO.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var d D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentInfoD, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var e E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentInfoE, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var f F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentInfoF, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var g G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentInfoG, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var h H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentInfoH, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var i I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentInfoI, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var j J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentInfoJ, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var k K
			if fetchK {
				k, err = fetchComponentForQueryResult[K](q.componentInfoK, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var l L
			if fetchL {
				l, err = fetchComponentForQueryResult[L](q.componentInfoL, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var m M
			if fetchM {
				m, err = fetchComponentForQueryResult[M](q.componentInfoM, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var n N
			if fetchN {
				n, err = fetchComponentForQueryResult[N](q.componentInfoN, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var o O
			if fetchO {
				o, err = fetchComponentForQueryResult[O](q.componentInfoO, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.componentsD = append(q.componentsD, d)
			q.componentsE = append(q.componentsE, e)
			q.componentsF = append(q.componentsF, f)
			q.componentsG = append(q.componentsG, g)
			q.componentsH = append(q.componentsH, h)
			q.componentsI = append(q.componentsI, i)
			q.componentsJ = append(q.componentsJ, j)
			q.componentsK = append(q.componentsK, k)
			q.componentsL = append(q.componentsL, l)
			q.componentsM = append(q.componentsM, m)
			q.componentsN = append(q.componentsN, n)
			q.componentsO = append(q.componentsO, o)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}

func (q *Query16[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Options]) Exec(world *World) (err error) {
	q.ClearResults()

	for _, archetype := range world.archetypeStorage.componentsHashToArchetype {
		if q.options.isArchetypeFilteredOut(archetype) {
			continue
		}

		fetchA, skip := shouldHandleQueryComponent(q.componentInfoA.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentInfoB.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentInfoC.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentInfoD.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentInfoE.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentInfoF.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentInfoG.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentInfoH.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentInfoI.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentInfoJ.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchK, skip := shouldHandleQueryComponent(q.componentInfoK.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchL, skip := shouldHandleQueryComponent(q.componentInfoL.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchM, skip := shouldHandleQueryComponent(q.componentInfoM.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchN, skip := shouldHandleQueryComponent(q.componentInfoN.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchO, skip := shouldHandleQueryComponent(q.componentInfoO.id, archetype, &q.options)
		if skip {
			continue
		}
		fetchP, skip := shouldHandleQueryComponent(q.componentInfoP.id, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentInfoA, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var b B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentInfoB, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var c C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentInfoC, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var d D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentInfoD, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var e E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentInfoE, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var f F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentInfoF, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var g G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentInfoG, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var h H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentInfoH, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var i I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentInfoI, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var j J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentInfoJ, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var k K
			if fetchK {
				k, err = fetchComponentForQueryResult[K](q.componentInfoK, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var l L
			if fetchL {
				l, err = fetchComponentForQueryResult[L](q.componentInfoL, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var m M
			if fetchM {
				m, err = fetchComponentForQueryResult[M](q.componentInfoM, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var n N
			if fetchN {
				n, err = fetchComponentForQueryResult[N](q.componentInfoN, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var o O
			if fetchO {
				o, err = fetchComponentForQueryResult[O](q.componentInfoO, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			var p P
			if fetchP {
				p, err = fetchComponentForQueryResult[P](q.componentInfoP, world.entities[entity].row, archetype)
				if err != nil {
					return err
				}
			}

			q.componentsA = append(q.componentsA, a)
			q.componentsB = append(q.componentsB, b)
			q.componentsC = append(q.componentsC, c)
			q.componentsD = append(q.componentsD, d)
			q.componentsE = append(q.componentsE, e)
			q.componentsF = append(q.componentsF, f)
			q.componentsG = append(q.componentsG, g)
			q.componentsH = append(q.componentsH, h)
			q.componentsI = append(q.componentsI, i)
			q.componentsJ = append(q.componentsJ, j)
			q.componentsK = append(q.componentsK, k)
			q.componentsL = append(q.componentsL, l)
			q.componentsM = append(q.componentsM, m)
			q.componentsN = append(q.componentsN, n)
			q.componentsO = append(q.componentsO, o)
			q.componentsP = append(q.componentsP, p)
			q.entityIds = append(q.entityIds, entity)
		}
	}

	return nil
}

func (q *Query0[QueryOptions]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	_, q.options, err = getQueryOptions[QueryOptions](world, otherWorlds)
	if err != nil {
		return err
	}

	q.components = []ComponentId{}
	q.options.optimize(q.components)
	return nil
}
func (q *Query1[A, QueryOptions]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[QueryOptions](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.components = []ComponentId{
		q.componentInfoA.id,
	}
	q.options.optimize(q.components)
	return nil
}
func (q *Query2[A, B, QueryOptions]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[QueryOptions](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.components = []ComponentId{
		q.componentInfoA.id,
		q.componentInfoB.id,
	}
	q.options.optimize(q.components)
	return nil
}
func (q *Query3[A, B, C, QueryOptions]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[QueryOptions](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.components = []ComponentId{
		q.componentInfoA.id,
		q.componentInfoB.id,
		q.componentInfoC.id,
	}
	q.options.optimize(q.components)
	return nil
}
func (q *Query4[A, B, C, D, QueryOptions]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[QueryOptions](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.componentInfoD = queryComponentInfoFor[D](targetWorld)
	q.components = []ComponentId{
		q.componentInfoA.id,
		q.componentInfoB.id,
		q.componentInfoC.id,
		q.componentInfoD.id,
	}
	q.options.optimize(q.components)
	return nil
}
func (q *Query5[A, B, C, D, E, QueryOptions]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[QueryOptions](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.componentInfoD = queryComponentInfoFor[D](targetWorld)
	q.componentInfoE = queryComponentInfoFor[E](targetWorld)
	q.components = []ComponentId{
		q.componentInfoA.id,
		q.componentInfoB.id,
		q.componentInfoC.id,
		q.componentInfoD.id,
		q.componentInfoE.id,
	}
	q.options.optimize(q.components)
	return nil
}
func (q *Query6[A, B, C, D, E, F, QueryOptions]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[QueryOptions](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.componentInfoD = queryComponentInfoFor[D](targetWorld)
	q.componentInfoE = queryComponentInfoFor[E](targetWorld)
	q.componentInfoF = queryComponentInfoFor[F](targetWorld)
	q.components = []ComponentId{
		q.componentInfoA.id,
		q.componentInfoB.id,
		q.componentInfoC.id,
		q.componentInfoD.id,
		q.componentInfoE.id,
		q.componentInfoF.id,
	}
	q.options.optimize(q.components)
	return nil
}
func (q *Query7[A, B, C, D, E, F, G, QueryOptions]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[QueryOptions](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.componentInfoD = queryComponentInfoFor[D](targetWorld)
	q.componentInfoE = queryComponentInfoFor[E](targetWorld)
	q.componentInfoF = queryComponentInfoFor[F](targetWorld)
	q.componentInfoG = queryComponentInfoFor[G](targetWorld)
	q.components = []ComponentId{
		q.componentInfoA.id,
		q.componentInfoB.id,
		q.componentInfoC.id,
		q.componentInfoD.id,
		q.componentInfoE.id,
		q.componentInfoF.id,
		q.componentInfoG.id,
	}
	q.options.optimize(q.components)
	return nil
}
func (q *Query8[A, B, C, D, E, F, G, H, QueryOptions]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[QueryOptions](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.componentInfoD = queryComponentInfoFor[D](targetWorld)
	q.componentInfoE = queryComponentInfoFor[E](targetWorld)
	q.componentInfoF = queryComponentInfoFor[F](targetWorld)
	q.componentInfoG = queryComponentInfoFor[G](targetWorld)
	q.componentInfoH = queryComponentInfoFor[H](targetWorld)
	q.components = []ComponentId{
		q.componentInfoA.id,
		q.componentInfoB.id,
		q.componentInfoC.id,
		q.componentInfoD.id,
		q.componentInfoE.id,
		q.componentInfoF.id,
		q.componentInfoG.id,
		q.componentInfoH.id,
	}
	q.options.optimize(q.components)
	return nil
}

func (q *Query9[A, B, C, D, E, F, G, H, I, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.componentInfoD = queryComponentInfoFor[D](targetWorld)
	q.componentInfoE = queryComponentInfoFor[E](targetWorld)
	q.componentInfoF = queryComponentInfoFor[F](targetWorld)
	q.componentInfoG = queryComponentInfoFor[G](targetWorld)
	q.componentInfoH = queryComponentInfoFor[H](targetWorld)
	q.componentInfoI = queryComponentInfoFor[I](targetWorld)
	q.components = []ComponentId{q.componentInfoA.id, q.componentInfoB.id, q.componentInfoC.id, q.componentInfoD.id, q.componentInfoE.id, q.componentInfoF.id, q.componentInfoG.id, q.componentInfoH.id, q.componentInfoI.id}
	q.options.optimize(q.components)
	return nil
}

func (q *Query10[A, B, C, D, E, F, G, H, I, J, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.componentInfoD = queryComponentInfoFor[D](targetWorld)
	q.componentInfoE = queryComponentInfoFor[E](targetWorld)
	q.componentInfoF = queryComponentInfoFor[F](targetWorld)
	q.componentInfoG = queryComponentInfoFor[G](targetWorld)
	q.componentInfoH = queryComponentInfoFor[H](targetWorld)
	q.componentInfoI = queryComponentInfoFor[I](targetWorld)
	q.componentInfoJ = queryComponentInfoFor[J](targetWorld)
	q.components = []ComponentId{q.componentInfoA.id, q.componentInfoB.id, q.componentInfoC.id, q.componentInfoD.id, q.componentInfoE.id, q.componentInfoF.id, q.componentInfoG.id, q.componentInfoH.id, q.componentInfoI.id, q.componentInfoJ.id}
	q.options.optimize(q.components)
	return nil
}

func (q *Query11[A, B, C, D, E, F, G, H, I, J, K, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.componentInfoD = queryComponentInfoFor[D](targetWorld)
	q.componentInfoE = queryComponentInfoFor[E](targetWorld)
	q.componentInfoF = queryComponentInfoFor[F](targetWorld)
	q.componentInfoG = queryComponentInfoFor[G](targetWorld)
	q.componentInfoH = queryComponentInfoFor[H](targetWorld)
	q.componentInfoI = queryComponentInfoFor[I](targetWorld)
	q.componentInfoJ = queryComponentInfoFor[J](targetWorld)
	q.componentInfoK = queryComponentInfoFor[K](targetWorld)
	q.components = []ComponentId{q.componentInfoA.id, q.componentInfoB.id, q.componentInfoC.id, q.componentInfoD.id, q.componentInfoE.id, q.componentInfoF.id, q.componentInfoG.id, q.componentInfoH.id, q.componentInfoI.id, q.componentInfoJ.id, q.componentInfoK.id}
	q.options.optimize(q.components)
	return nil
}

func (q *Query12[A, B, C, D, E, F, G, H, I, J, K, L, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.componentInfoD = queryComponentInfoFor[D](targetWorld)
	q.componentInfoE = queryComponentInfoFor[E](targetWorld)
	q.componentInfoF = queryComponentInfoFor[F](targetWorld)
	q.componentInfoG = queryComponentInfoFor[G](targetWorld)
	q.componentInfoH = queryComponentInfoFor[H](targetWorld)
	q.componentInfoI = queryComponentInfoFor[I](targetWorld)
	q.componentInfoJ = queryComponentInfoFor[J](targetWorld)
	q.componentInfoK = queryComponentInfoFor[K](targetWorld)
	q.componentInfoL = queryComponentInfoFor[L](targetWorld)
	q.components = []ComponentId{q.componentInfoA.id, q.componentInfoB.id, q.componentInfoC.id, q.componentInfoD.id, q.componentInfoE.id, q.componentInfoF.id, q.componentInfoG.id, q.componentInfoH.id, q.componentInfoI.id, q.componentInfoJ.id, q.componentInfoK.id, q.componentInfoL.id}
	q.options.optimize(q.components)
	return nil
}

func (q *Query13[A, B, C, D, E, F, G, H, I, J, K, L, M, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.componentInfoD = queryComponentInfoFor[D](targetWorld)
	q.componentInfoE = queryComponentInfoFor[E](targetWorld)
	q.componentInfoF = queryComponentInfoFor[F](targetWorld)
	q.componentInfoG = queryComponentInfoFor[G](targetWorld)
	q.componentInfoH = queryComponentInfoFor[H](targetWorld)
	q.componentInfoI = queryComponentInfoFor[I](targetWorld)
	q.componentInfoJ = queryComponentInfoFor[J](targetWorld)
	q.componentInfoK = queryComponentInfoFor[K](targetWorld)
	q.componentInfoL = queryComponentInfoFor[L](targetWorld)
	q.componentInfoM = queryComponentInfoFor[M](targetWorld)
	q.components = []ComponentId{q.componentInfoA.id, q.componentInfoB.id, q.componentInfoC.id, q.componentInfoD.id, q.componentInfoE.id, q.componentInfoF.id, q.componentInfoG.id, q.componentInfoH.id, q.componentInfoI.id, q.componentInfoJ.id, q.componentInfoK.id, q.componentInfoL.id, q.componentInfoM.id}
	q.options.optimize(q.components)
	return nil
}

func (q *Query14[A, B, C, D, E, F, G, H, I, J, K, L, M, N, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.componentInfoD = queryComponentInfoFor[D](targetWorld)
	q.componentInfoE = queryComponentInfoFor[E](targetWorld)
	q.componentInfoF = queryComponentInfoFor[F](targetWorld)
	q.componentInfoG = queryComponentInfoFor[G](targetWorld)
	q.componentInfoH = queryComponentInfoFor[H](targetWorld)
	q.componentInfoI = queryComponentInfoFor[I](targetWorld)
	q.componentInfoJ = queryComponentInfoFor[J](targetWorld)
	q.componentInfoK = queryComponentInfoFor[K](targetWorld)
	q.componentInfoL = queryComponentInfoFor[L](targetWorld)
	q.componentInfoM = queryComponentInfoFor[M](targetWorld)
	q.componentInfoN = queryComponentInfoFor[N](targetWorld)
	q.components = []ComponentId{q.componentInfoA.id, q.componentInfoB.id, q.componentInfoC.id, q.componentInfoD.id, q.componentInfoE.id, q.componentInfoF.id, q.componentInfoG.id, q.componentInfoH.id, q.componentInfoI.id, q.componentInfoJ.id, q.componentInfoK.id, q.componentInfoL.id, q.componentInfoM.id, q.componentInfoN.id}
	q.options.optimize(q.components)
	return nil
}

func (q *Query15[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.componentInfoD = queryComponentInfoFor[D](targetWorld)
	q.componentInfoE = queryComponentInfoFor[E](targetWorld)
	q.componentInfoF = queryComponentInfoFor[F](targetWorld)
	q.componentInfoG = queryComponentInfoFor[G](targetWorld)
	q.componentInfoH = queryComponentInfoFor[H](targetWorld)
	q.componentInfoI = queryComponentInfoFor[I](targetWorld)
	q.componentInfoJ = queryComponentInfoFor[J](targetWorld)
	q.componentInfoK = queryComponentInfoFor[K](targetWorld)
	q.componentInfoL = queryComponentInfoFor[L](targetWorld)
	q.componentInfoM = queryComponentInfoFor[M](targetWorld)
	q.componentInfoN = queryComponentInfoFor[N](targetWorld)
	q.componentInfoO = queryComponentInfoFor[O](targetWorld)
	q.components = []ComponentId{q.componentInfoA.id, q.componentInfoB.id, q.componentInfoC.id, q.componentInfoD.id, q.componentInfoE.id, q.componentInfoF.id, q.componentInfoG.id, q.componentInfoH.id, q.componentInfoI.id, q.componentInfoJ.id, q.componentInfoK.id, q.componentInfoL.id, q.componentInfoM.id, q.componentInfoN.id, q.componentInfoO.id}
	q.options.optimize(q.components)
	return nil
}

func (q *Query16[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentInfoA = queryComponentInfoFor[A](targetWorld)
	q.componentInfoB = queryComponentInfoFor[B](targetWorld)
	q.componentInfoC = queryComponentInfoFor[C](targetWorld)
	q.componentInfoD = queryComponentInfoFor[D](targetWorld)
	q.componentInfoE = queryComponentInfoFor[E](targetWorld)
	q.componentInfoF = queryComponentInfoFor[F](targetWorld)
	q.componentInfoG = queryComponentInfoFor[G](targetWorld)
	q.componentInfoH = queryComponentInfoFor[H](targetWorld)
	q.componentInfoI = queryComponentInfoFor[I](targetWorld)
	q.componentInfoJ = queryComponentInfoFor[J](targetWorld)
	q.componentInfoK = queryComponentInfoFor[K](targetWorld)
	q.componentInfoL = queryComponentInfoFor[L](targetWorld)
	q.componentInfoM = queryComponentInfoFor[M](targetWorld)
	q.componentInfoN = queryComponentInfoFor[N](targetWorld)
	q.componentInfoO = queryComponentInfoFor[O](targetWorld)
	q.componentInfoP = queryComponentInfoFor[P](targetWorld)
	q.components = []ComponentId{
		q.componentInfoA.id, q.componentInfoB.id, q.componentInfoC.id, q.componentInfoD.id, q.componentInfoE.id, q.componentInfoF.id, q.componentInfoG.id, q.componentInfoH.id,
		q.componentInfoI.id, q.componentInfoJ.id, q.componentInfoK.id, q.componentInfoL.id, q.componentInfoM.id, q.componentInfoN.id, q.componentInfoO.id, q.componentInfoP.id,
	}
	q.options.optimize(q.components)
	return nil
}

func getQueryOptions[QueryOptions QueryOption](world *World, otherWorlds *map[WorldId]*World) (*World, CombinedQueryOptions, error) {
	queryOptions, err := toCombinedQueryOptions[QueryOptions](world)
	if err != nil {
		return nil, queryOptions, err
	}

	if queryOptions.TargetWorld == nil {
		return world, queryOptions, nil
	}

	// If the query targets another world, the stored componentIDs are retrieved from the wrong world. Unfortunately,
	// we need to get all query options before we can know if this query targets another world. We now get the queryOptions
	// again, this time with the correct target world.
	targetWorld := (*otherWorlds)[*queryOptions.TargetWorld]
	if targetWorld == nil {
		return nil, queryOptions, ErrTargetWorldNotFound
	}

	queryOptions, err = toCombinedQueryOptions[QueryOptions](targetWorld)
	return targetWorld, queryOptions, err
}

func (q *Query0[Options]) ClearResults()                                                  { q.Clear() }
func (q *Query1[A, Options]) ClearResults()                                               { q.Clear() }
func (q *Query2[A, B, Options]) ClearResults()                                            { q.Clear() }
func (q *Query3[A, B, C, Options]) ClearResults()                                         { q.Clear() }
func (q *Query4[A, B, C, D, Options]) ClearResults()                                      { q.Clear() }
func (q *Query5[A, B, C, D, E, Options]) ClearResults()                                   { q.Clear() }
func (q *Query6[A, B, C, D, E, F, Options]) ClearResults()                                { q.Clear() }
func (q *Query7[A, B, C, D, E, F, G, Options]) ClearResults()                             { q.Clear() }
func (q *Query8[A, B, C, D, E, F, G, H, Options]) ClearResults()                          { q.Clear() }
func (q *Query9[A, B, C, D, E, F, G, H, I, Options]) ClearResults()                       { q.Clear() }
func (q *Query10[A, B, C, D, E, F, G, H, I, J, Options]) ClearResults()                   { q.Clear() }
func (q *Query11[A, B, C, D, E, F, G, H, I, J, K, Options]) ClearResults()                { q.Clear() }
func (q *Query12[A, B, C, D, E, F, G, H, I, J, K, L, Options]) ClearResults()             { q.Clear() }
func (q *Query13[A, B, C, D, E, F, G, H, I, J, K, L, M, Options]) ClearResults()          { q.Clear() }
func (q *Query14[A, B, C, D, E, F, G, H, I, J, K, L, M, N, Options]) ClearResults()       { q.Clear() }
func (q *Query15[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, Options]) ClearResults()    { q.Clear() }
func (q *Query16[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Options]) ClearResults() { q.Clear() }

// shouldHandleQueryComponent returns wether a component should be fetched and/or skipped:
//   - shouldSkip=true means that the archetype should not be included in the query results.
//   - shouldSkip=false means that the archetype should be included in the query results.
//   - shouldFetch=true means that the (entities of the) archetype should be included in the query and should be fetched from a component storage.
//   - shouldFetch=false means that the (entities of the) archetype should be included in the query, but should be nil.
func shouldHandleQueryComponent(componentId ComponentId, archetype *Archetype, queryOptions *CombinedQueryOptions) (shouldFetch bool, shouldSkip bool) {
	if archetype.HasComponent(componentId) {
		return true, false
	}

	if slices.Contains(queryOptions.OptionalComponents, componentId) {
		return false, false
	}

	return false, true
}

// fetchComponentForQueryResult fetches a component from the component storage
func fetchComponentForQueryResult[T IComponent](componentInfo queryComponentInfo, entityRow uint, archetype *Archetype) (result T, err error) {
	storage := archetype.components[componentInfo.id]
	result, err = getComponentFromComponentStorage[T](storage, entityRow, componentInfo.isPointer)
	if err != nil {
		return result, fmt.Errorf("failed to retrieve component %s from storage: %v", componentInfo.id.DebugString(), err)
	}

	return result, nil
}
