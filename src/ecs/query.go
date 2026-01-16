package ecs

import (
	"fmt"
	"slices"

	"github.com/lucdrenth/murphecs/src/utils"
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
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1] to specify if components are read-only. Marking components as
//     read-only allows systems with queries as system-params to be run in parallel with other systems.
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

	componentIdA ComponentId
}

// Query2 queries 2 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] to specify if components are read-only. Marking
//     components as read-only allows systems with queries as system-params to be run in parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
}

// Query3 queries 3 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
}

// Query4 queries 4 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
	componentIdD ComponentId
}

// Query5 queries 5 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
	componentIdD ComponentId
	componentIdE ComponentId
}

// Query6 queries 6 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
	componentIdD ComponentId
	componentIdE ComponentId
	componentIdF ComponentId
}

// Query7 queries 7 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
	componentIdD ComponentId
	componentIdE ComponentId
	componentIdF ComponentId
	componentIdG ComponentId
}

// Query8 queries 8 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
	componentIdD ComponentId
	componentIdE ComponentId
	componentIdF ComponentId
	componentIdG ComponentId
	componentIdH ComponentId
}

// Query9 queries 9 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
	componentIdD ComponentId
	componentIdE ComponentId
	componentIdF ComponentId
	componentIdG ComponentId
	componentIdH ComponentId
	componentIdI ComponentId
}

// Query10 queries 10 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
	componentIdD ComponentId
	componentIdE ComponentId
	componentIdF ComponentId
	componentIdG ComponentId
	componentIdH ComponentId
	componentIdI ComponentId
	componentIdJ ComponentId
}

// Query11 queries 11 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
	componentIdD ComponentId
	componentIdE ComponentId
	componentIdF ComponentId
	componentIdG ComponentId
	componentIdH ComponentId
	componentIdI ComponentId
	componentIdJ ComponentId
	componentIdK ComponentId
}

// Query12 queries 12 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
	componentIdD ComponentId
	componentIdE ComponentId
	componentIdF ComponentId
	componentIdG ComponentId
	componentIdH ComponentId
	componentIdI ComponentId
	componentIdJ ComponentId
	componentIdK ComponentId
	componentIdL ComponentId
}

// Query13 queries 13 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
	componentIdD ComponentId
	componentIdE ComponentId
	componentIdF ComponentId
	componentIdG ComponentId
	componentIdH ComponentId
	componentIdI ComponentId
	componentIdJ ComponentId
	componentIdK ComponentId
	componentIdL ComponentId
	componentIdM ComponentId
}

// Query14 queries 14 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
	componentIdD ComponentId
	componentIdE ComponentId
	componentIdF ComponentId
	componentIdG ComponentId
	componentIdH ComponentId
	componentIdI ComponentId
	componentIdJ ComponentId
	componentIdK ComponentId
	componentIdL ComponentId
	componentIdM ComponentId
	componentIdN ComponentId
}

// Query15 queries 15 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
	componentIdD ComponentId
	componentIdE ComponentId
	componentIdF ComponentId
	componentIdG ComponentId
	componentIdH ComponentId
	componentIdI ComponentId
	componentIdJ ComponentId
	componentIdK ComponentId
	componentIdL ComponentId
	componentIdM ComponentId
	componentIdN ComponentId
	componentIdO ComponentId
}

// Query16 queries 16 components.
//
// Prepare must be called once before calling Execute.
//
// The following query options are available:
//   - use [NoReadOnly], [AllReadOnly], [ReadOnly1], [ReadOnly2] (and so on) to specify if components
//     are read-only. Marking components as read-only allows systems with queries as system-params to be run in
//     parallel with other systems.
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

	componentIdA ComponentId
	componentIdB ComponentId
	componentIdC ComponentId
	componentIdD ComponentId
	componentIdE ComponentId
	componentIdF ComponentId
	componentIdG ComponentId
	componentIdH ComponentId
	componentIdI ComponentId
	componentIdJ ComponentId
	componentIdK ComponentId
	componentIdL ComponentId
	componentIdM ComponentId
	componentIdN ComponentId
	componentIdO ComponentId
	componentIdP ComponentId
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentIdA, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentIdB, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *ComponentC
			if fetchC {
				c, err = fetchComponentForQueryResult[ComponentC](q.componentIdC, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentIdD, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *ComponentC
			if fetchC {
				c, err = fetchComponentForQueryResult[ComponentC](q.componentIdC, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var d *ComponentD
			if fetchD {
				d, err = fetchComponentForQueryResult[ComponentD](q.componentIdD, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentIdD, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentIdE, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *ComponentC
			if fetchC {
				c, err = fetchComponentForQueryResult[ComponentC](q.componentIdC, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var d *ComponentD
			if fetchD {
				d, err = fetchComponentForQueryResult[ComponentD](q.componentIdD, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var e *ComponentE
			if fetchE {
				e, err = fetchComponentForQueryResult[ComponentE](q.componentIdE, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentIdD, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentIdE, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentIdF, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *ComponentC
			if fetchC {
				c, err = fetchComponentForQueryResult[ComponentC](q.componentIdC, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var d *ComponentD
			if fetchD {
				d, err = fetchComponentForQueryResult[ComponentD](q.componentIdD, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var e *ComponentE
			if fetchE {
				e, err = fetchComponentForQueryResult[ComponentE](q.componentIdE, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var f *ComponentF
			if fetchF {
				f, err = fetchComponentForQueryResult[ComponentF](q.componentIdF, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentIdD, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentIdE, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentIdF, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentIdG, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *ComponentC
			if fetchC {
				c, err = fetchComponentForQueryResult[ComponentC](q.componentIdC, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var d *ComponentD
			if fetchD {
				d, err = fetchComponentForQueryResult[ComponentD](q.componentIdD, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var e *ComponentE
			if fetchE {
				e, err = fetchComponentForQueryResult[ComponentE](q.componentIdE, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var f *ComponentF
			if fetchF {
				f, err = fetchComponentForQueryResult[ComponentF](q.componentIdF, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var g *ComponentG
			if fetchG {
				g, err = fetchComponentForQueryResult[ComponentG](q.componentIdG, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentIdD, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentIdE, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentIdF, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentIdG, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentIdH, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *ComponentA
			if fetchA {
				a, err = fetchComponentForQueryResult[ComponentA](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *ComponentB
			if fetchB {
				b, err = fetchComponentForQueryResult[ComponentB](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *ComponentC
			if fetchC {
				c, err = fetchComponentForQueryResult[ComponentC](q.componentIdC, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var d *ComponentD
			if fetchD {
				d, err = fetchComponentForQueryResult[ComponentD](q.componentIdD, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var e *ComponentE
			if fetchE {
				e, err = fetchComponentForQueryResult[ComponentE](q.componentIdE, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var f *ComponentF
			if fetchF {
				f, err = fetchComponentForQueryResult[ComponentF](q.componentIdF, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var g *ComponentG
			if fetchG {
				g, err = fetchComponentForQueryResult[ComponentG](q.componentIdG, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var h *ComponentH
			if fetchH {
				h, err = fetchComponentForQueryResult[ComponentH](q.componentIdH, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentIdD, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentIdE, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentIdF, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentIdG, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentIdH, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentIdI, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentIdC, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var d *D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentIdD, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var e *E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentIdE, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var f *F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentIdF, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var g *G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentIdG, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var h *H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentIdH, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var i *I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentIdI, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentIdD, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentIdE, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentIdF, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentIdG, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentIdH, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentIdI, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentIdJ, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentIdC, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var d *D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentIdD, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var e *E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentIdE, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var f *F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentIdF, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var g *G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentIdG, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var h *H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentIdH, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var i *I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentIdI, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var j *J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentIdJ, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentIdD, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentIdE, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentIdF, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentIdG, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentIdH, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentIdI, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentIdJ, archetype, &q.options)
		if skip {
			continue
		}
		fetchK, skip := shouldHandleQueryComponent(q.componentIdK, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentIdC, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var d *D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentIdD, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var e *E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentIdE, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var f *F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentIdF, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var g *G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentIdG, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var h *H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentIdH, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var i *I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentIdI, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var j *J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentIdJ, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var k *K
			if fetchK {
				k, err = fetchComponentForQueryResult[K](q.componentIdK, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentIdD, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentIdE, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentIdF, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentIdG, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentIdH, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentIdI, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentIdJ, archetype, &q.options)
		if skip {
			continue
		}
		fetchK, skip := shouldHandleQueryComponent(q.componentIdK, archetype, &q.options)
		if skip {
			continue
		}
		fetchL, skip := shouldHandleQueryComponent(q.componentIdL, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentIdC, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var d *D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentIdD, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var e *E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentIdE, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var f *F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentIdF, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var g *G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentIdG, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var h *H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentIdH, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var i *I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentIdI, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var j *J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentIdJ, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var k *K
			if fetchK {
				k, err = fetchComponentForQueryResult[K](q.componentIdK, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var l *L
			if fetchL {
				l, err = fetchComponentForQueryResult[L](q.componentIdL, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentIdD, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentIdE, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentIdF, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentIdG, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentIdH, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentIdI, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentIdJ, archetype, &q.options)
		if skip {
			continue
		}
		fetchK, skip := shouldHandleQueryComponent(q.componentIdK, archetype, &q.options)
		if skip {
			continue
		}
		fetchL, skip := shouldHandleQueryComponent(q.componentIdL, archetype, &q.options)
		if skip {
			continue
		}
		fetchM, skip := shouldHandleQueryComponent(q.componentIdM, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentIdC, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var d *D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentIdD, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var e *E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentIdE, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var f *F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentIdF, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var g *G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentIdG, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var h *H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentIdH, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var i *I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentIdI, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var j *J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentIdJ, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var k *K
			if fetchK {
				k, err = fetchComponentForQueryResult[K](q.componentIdK, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var l *L
			if fetchL {
				l, err = fetchComponentForQueryResult[L](q.componentIdL, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var m *M
			if fetchM {
				m, err = fetchComponentForQueryResult[M](q.componentIdM, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentIdD, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentIdE, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentIdF, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentIdG, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentIdH, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentIdI, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentIdJ, archetype, &q.options)
		if skip {
			continue
		}
		fetchK, skip := shouldHandleQueryComponent(q.componentIdK, archetype, &q.options)
		if skip {
			continue
		}
		fetchL, skip := shouldHandleQueryComponent(q.componentIdL, archetype, &q.options)
		if skip {
			continue
		}
		fetchM, skip := shouldHandleQueryComponent(q.componentIdM, archetype, &q.options)
		if skip {
			continue
		}
		fetchN, skip := shouldHandleQueryComponent(q.componentIdN, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentIdC, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var d *D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentIdD, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var e *E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentIdE, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var f *F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentIdF, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var g *G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentIdG, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var h *H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentIdH, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var i *I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentIdI, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var j *J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentIdJ, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var k *K
			if fetchK {
				k, err = fetchComponentForQueryResult[K](q.componentIdK, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var l *L
			if fetchL {
				l, err = fetchComponentForQueryResult[L](q.componentIdL, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var m *M
			if fetchM {
				m, err = fetchComponentForQueryResult[M](q.componentIdM, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var n *N
			if fetchN {
				n, err = fetchComponentForQueryResult[N](q.componentIdN, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentIdD, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentIdE, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentIdF, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentIdG, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentIdH, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentIdI, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentIdJ, archetype, &q.options)
		if skip {
			continue
		}
		fetchK, skip := shouldHandleQueryComponent(q.componentIdK, archetype, &q.options)
		if skip {
			continue
		}
		fetchL, skip := shouldHandleQueryComponent(q.componentIdL, archetype, &q.options)
		if skip {
			continue
		}
		fetchM, skip := shouldHandleQueryComponent(q.componentIdM, archetype, &q.options)
		if skip {
			continue
		}
		fetchN, skip := shouldHandleQueryComponent(q.componentIdN, archetype, &q.options)
		if skip {
			continue
		}
		fetchO, skip := shouldHandleQueryComponent(q.componentIdO, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentIdC, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var d *D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentIdD, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var e *E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentIdE, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var f *F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentIdF, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var g *G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentIdG, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var h *H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentIdH, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var i *I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentIdI, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var j *J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentIdJ, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var k *K
			if fetchK {
				k, err = fetchComponentForQueryResult[K](q.componentIdK, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var l *L
			if fetchL {
				l, err = fetchComponentForQueryResult[L](q.componentIdL, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var m *M
			if fetchM {
				m, err = fetchComponentForQueryResult[M](q.componentIdM, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var n *N
			if fetchN {
				n, err = fetchComponentForQueryResult[N](q.componentIdN, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var o *O
			if fetchO {
				o, err = fetchComponentForQueryResult[O](q.componentIdO, world.entities[entity].row, archetype, &q.options)
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

		fetchA, skip := shouldHandleQueryComponent(q.componentIdA, archetype, &q.options)
		if skip {
			continue
		}
		fetchB, skip := shouldHandleQueryComponent(q.componentIdB, archetype, &q.options)
		if skip {
			continue
		}
		fetchC, skip := shouldHandleQueryComponent(q.componentIdC, archetype, &q.options)
		if skip {
			continue
		}
		fetchD, skip := shouldHandleQueryComponent(q.componentIdD, archetype, &q.options)
		if skip {
			continue
		}
		fetchE, skip := shouldHandleQueryComponent(q.componentIdE, archetype, &q.options)
		if skip {
			continue
		}
		fetchF, skip := shouldHandleQueryComponent(q.componentIdF, archetype, &q.options)
		if skip {
			continue
		}
		fetchG, skip := shouldHandleQueryComponent(q.componentIdG, archetype, &q.options)
		if skip {
			continue
		}
		fetchH, skip := shouldHandleQueryComponent(q.componentIdH, archetype, &q.options)
		if skip {
			continue
		}
		fetchI, skip := shouldHandleQueryComponent(q.componentIdI, archetype, &q.options)
		if skip {
			continue
		}
		fetchJ, skip := shouldHandleQueryComponent(q.componentIdJ, archetype, &q.options)
		if skip {
			continue
		}
		fetchK, skip := shouldHandleQueryComponent(q.componentIdK, archetype, &q.options)
		if skip {
			continue
		}
		fetchL, skip := shouldHandleQueryComponent(q.componentIdL, archetype, &q.options)
		if skip {
			continue
		}
		fetchM, skip := shouldHandleQueryComponent(q.componentIdM, archetype, &q.options)
		if skip {
			continue
		}
		fetchN, skip := shouldHandleQueryComponent(q.componentIdN, archetype, &q.options)
		if skip {
			continue
		}
		fetchO, skip := shouldHandleQueryComponent(q.componentIdO, archetype, &q.options)
		if skip {
			continue
		}
		fetchP, skip := shouldHandleQueryComponent(q.componentIdP, archetype, &q.options)
		if skip {
			continue
		}

		for _, entity := range archetype.entities {
			var a *A
			if fetchA {
				a, err = fetchComponentForQueryResult[A](q.componentIdA, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var b *B
			if fetchB {
				b, err = fetchComponentForQueryResult[B](q.componentIdB, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var c *C
			if fetchC {
				c, err = fetchComponentForQueryResult[C](q.componentIdC, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var d *D
			if fetchD {
				d, err = fetchComponentForQueryResult[D](q.componentIdD, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var e *E
			if fetchE {
				e, err = fetchComponentForQueryResult[E](q.componentIdE, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var f *F
			if fetchF {
				f, err = fetchComponentForQueryResult[F](q.componentIdF, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var g *G
			if fetchG {
				g, err = fetchComponentForQueryResult[G](q.componentIdG, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var h *H
			if fetchH {
				h, err = fetchComponentForQueryResult[H](q.componentIdH, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var i *I
			if fetchI {
				i, err = fetchComponentForQueryResult[I](q.componentIdI, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var j *J
			if fetchJ {
				j, err = fetchComponentForQueryResult[J](q.componentIdJ, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var k *K
			if fetchK {
				k, err = fetchComponentForQueryResult[K](q.componentIdK, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var l *L
			if fetchL {
				l, err = fetchComponentForQueryResult[L](q.componentIdL, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var m *M
			if fetchM {
				m, err = fetchComponentForQueryResult[M](q.componentIdM, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var n *N
			if fetchN {
				n, err = fetchComponentForQueryResult[N](q.componentIdN, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var o *O
			if fetchO {
				o, err = fetchComponentForQueryResult[O](q.componentIdO, world.entities[entity].row, archetype, &q.options)
				if err != nil {
					return err
				}
			}

			var p *P
			if fetchP {
				p, err = fetchComponentForQueryResult[P](q.componentIdP, world.entities[entity].row, archetype, &q.options)
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

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.components = []ComponentId{
		q.componentIdA,
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

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.components = []ComponentId{
		q.componentIdA,
		q.componentIdB,
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

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.components = []ComponentId{
		q.componentIdA,
		q.componentIdB,
		q.componentIdC,
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

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.componentIdD = ComponentIdFor[D](targetWorld)
	q.components = []ComponentId{
		q.componentIdA,
		q.componentIdB,
		q.componentIdC,
		q.componentIdD,
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

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.componentIdD = ComponentIdFor[D](targetWorld)
	q.componentIdE = ComponentIdFor[E](targetWorld)
	q.components = []ComponentId{
		q.componentIdA,
		q.componentIdB,
		q.componentIdC,
		q.componentIdD,
		q.componentIdE,
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

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.componentIdD = ComponentIdFor[D](targetWorld)
	q.componentIdE = ComponentIdFor[E](targetWorld)
	q.componentIdF = ComponentIdFor[F](targetWorld)
	q.components = []ComponentId{
		q.componentIdA,
		q.componentIdB,
		q.componentIdC,
		q.componentIdD,
		q.componentIdE,
		q.componentIdF,
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

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.componentIdD = ComponentIdFor[D](targetWorld)
	q.componentIdE = ComponentIdFor[E](targetWorld)
	q.componentIdF = ComponentIdFor[F](targetWorld)
	q.componentIdG = ComponentIdFor[G](targetWorld)
	q.components = []ComponentId{
		q.componentIdA,
		q.componentIdB,
		q.componentIdC,
		q.componentIdD,
		q.componentIdE,
		q.componentIdF,
		q.componentIdG,
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

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.componentIdD = ComponentIdFor[D](targetWorld)
	q.componentIdE = ComponentIdFor[E](targetWorld)
	q.componentIdF = ComponentIdFor[F](targetWorld)
	q.componentIdG = ComponentIdFor[G](targetWorld)
	q.componentIdH = ComponentIdFor[H](targetWorld)
	q.components = []ComponentId{
		q.componentIdA,
		q.componentIdB,
		q.componentIdC,
		q.componentIdD,
		q.componentIdE,
		q.componentIdF,
		q.componentIdG,
		q.componentIdH,
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

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.componentIdD = ComponentIdFor[D](targetWorld)
	q.componentIdE = ComponentIdFor[E](targetWorld)
	q.componentIdF = ComponentIdFor[F](targetWorld)
	q.componentIdG = ComponentIdFor[G](targetWorld)
	q.componentIdH = ComponentIdFor[H](targetWorld)
	q.componentIdI = ComponentIdFor[I](targetWorld)
	q.components = []ComponentId{q.componentIdA, q.componentIdB, q.componentIdC, q.componentIdD, q.componentIdE, q.componentIdF, q.componentIdG, q.componentIdH, q.componentIdI}
	q.options.optimize(q.components)
	return nil
}

func (q *Query10[A, B, C, D, E, F, G, H, I, J, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.componentIdD = ComponentIdFor[D](targetWorld)
	q.componentIdE = ComponentIdFor[E](targetWorld)
	q.componentIdF = ComponentIdFor[F](targetWorld)
	q.componentIdG = ComponentIdFor[G](targetWorld)
	q.componentIdH = ComponentIdFor[H](targetWorld)
	q.componentIdI = ComponentIdFor[I](targetWorld)
	q.componentIdJ = ComponentIdFor[J](targetWorld)
	q.components = []ComponentId{q.componentIdA, q.componentIdB, q.componentIdC, q.componentIdD, q.componentIdE, q.componentIdF, q.componentIdG, q.componentIdH, q.componentIdI, q.componentIdJ}
	q.options.optimize(q.components)
	return nil
}

func (q *Query11[A, B, C, D, E, F, G, H, I, J, K, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.componentIdD = ComponentIdFor[D](targetWorld)
	q.componentIdE = ComponentIdFor[E](targetWorld)
	q.componentIdF = ComponentIdFor[F](targetWorld)
	q.componentIdG = ComponentIdFor[G](targetWorld)
	q.componentIdH = ComponentIdFor[H](targetWorld)
	q.componentIdI = ComponentIdFor[I](targetWorld)
	q.componentIdJ = ComponentIdFor[J](targetWorld)
	q.componentIdK = ComponentIdFor[K](targetWorld)
	q.components = []ComponentId{q.componentIdA, q.componentIdB, q.componentIdC, q.componentIdD, q.componentIdE, q.componentIdF, q.componentIdG, q.componentIdH, q.componentIdI, q.componentIdJ, q.componentIdK}
	q.options.optimize(q.components)
	return nil
}

func (q *Query12[A, B, C, D, E, F, G, H, I, J, K, L, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.componentIdD = ComponentIdFor[D](targetWorld)
	q.componentIdE = ComponentIdFor[E](targetWorld)
	q.componentIdF = ComponentIdFor[F](targetWorld)
	q.componentIdG = ComponentIdFor[G](targetWorld)
	q.componentIdH = ComponentIdFor[H](targetWorld)
	q.componentIdI = ComponentIdFor[I](targetWorld)
	q.componentIdJ = ComponentIdFor[J](targetWorld)
	q.componentIdK = ComponentIdFor[K](targetWorld)
	q.componentIdL = ComponentIdFor[L](targetWorld)
	q.components = []ComponentId{q.componentIdA, q.componentIdB, q.componentIdC, q.componentIdD, q.componentIdE, q.componentIdF, q.componentIdG, q.componentIdH, q.componentIdI, q.componentIdJ, q.componentIdK, q.componentIdL}
	q.options.optimize(q.components)
	return nil
}

func (q *Query13[A, B, C, D, E, F, G, H, I, J, K, L, M, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.componentIdD = ComponentIdFor[D](targetWorld)
	q.componentIdE = ComponentIdFor[E](targetWorld)
	q.componentIdF = ComponentIdFor[F](targetWorld)
	q.componentIdG = ComponentIdFor[G](targetWorld)
	q.componentIdH = ComponentIdFor[H](targetWorld)
	q.componentIdI = ComponentIdFor[I](targetWorld)
	q.componentIdJ = ComponentIdFor[J](targetWorld)
	q.componentIdK = ComponentIdFor[K](targetWorld)
	q.componentIdL = ComponentIdFor[L](targetWorld)
	q.componentIdM = ComponentIdFor[M](targetWorld)
	q.components = []ComponentId{q.componentIdA, q.componentIdB, q.componentIdC, q.componentIdD, q.componentIdE, q.componentIdF, q.componentIdG, q.componentIdH, q.componentIdI, q.componentIdJ, q.componentIdK, q.componentIdL, q.componentIdM}
	q.options.optimize(q.components)
	return nil
}

func (q *Query14[A, B, C, D, E, F, G, H, I, J, K, L, M, N, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.componentIdD = ComponentIdFor[D](targetWorld)
	q.componentIdE = ComponentIdFor[E](targetWorld)
	q.componentIdF = ComponentIdFor[F](targetWorld)
	q.componentIdG = ComponentIdFor[G](targetWorld)
	q.componentIdH = ComponentIdFor[H](targetWorld)
	q.componentIdI = ComponentIdFor[I](targetWorld)
	q.componentIdJ = ComponentIdFor[J](targetWorld)
	q.componentIdK = ComponentIdFor[K](targetWorld)
	q.componentIdL = ComponentIdFor[L](targetWorld)
	q.componentIdM = ComponentIdFor[M](targetWorld)
	q.componentIdN = ComponentIdFor[N](targetWorld)
	q.components = []ComponentId{q.componentIdA, q.componentIdB, q.componentIdC, q.componentIdD, q.componentIdE, q.componentIdF, q.componentIdG, q.componentIdH, q.componentIdI, q.componentIdJ, q.componentIdK, q.componentIdL, q.componentIdM, q.componentIdN}
	q.options.optimize(q.components)
	return nil
}

func (q *Query15[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.componentIdD = ComponentIdFor[D](targetWorld)
	q.componentIdE = ComponentIdFor[E](targetWorld)
	q.componentIdF = ComponentIdFor[F](targetWorld)
	q.componentIdG = ComponentIdFor[G](targetWorld)
	q.componentIdH = ComponentIdFor[H](targetWorld)
	q.componentIdI = ComponentIdFor[I](targetWorld)
	q.componentIdJ = ComponentIdFor[J](targetWorld)
	q.componentIdK = ComponentIdFor[K](targetWorld)
	q.componentIdL = ComponentIdFor[L](targetWorld)
	q.componentIdM = ComponentIdFor[M](targetWorld)
	q.componentIdN = ComponentIdFor[N](targetWorld)
	q.componentIdO = ComponentIdFor[O](targetWorld)
	q.components = []ComponentId{q.componentIdA, q.componentIdB, q.componentIdC, q.componentIdD, q.componentIdE, q.componentIdF, q.componentIdG, q.componentIdH, q.componentIdI, q.componentIdJ, q.componentIdK, q.componentIdL, q.componentIdM, q.componentIdN, q.componentIdO}
	q.options.optimize(q.components)
	return nil
}

func (q *Query16[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Options]) Prepare(world *World, otherWorlds *map[WorldId]*World) (err error) {
	var targetWorld *World
	targetWorld, q.options, err = getQueryOptions[Options](world, otherWorlds)
	if err != nil {
		return err
	}

	q.componentIdA = ComponentIdFor[A](targetWorld)
	q.componentIdB = ComponentIdFor[B](targetWorld)
	q.componentIdC = ComponentIdFor[C](targetWorld)
	q.componentIdD = ComponentIdFor[D](targetWorld)
	q.componentIdE = ComponentIdFor[E](targetWorld)
	q.componentIdF = ComponentIdFor[F](targetWorld)
	q.componentIdG = ComponentIdFor[G](targetWorld)
	q.componentIdH = ComponentIdFor[H](targetWorld)
	q.componentIdI = ComponentIdFor[I](targetWorld)
	q.componentIdJ = ComponentIdFor[J](targetWorld)
	q.componentIdK = ComponentIdFor[K](targetWorld)
	q.componentIdL = ComponentIdFor[L](targetWorld)
	q.componentIdM = ComponentIdFor[M](targetWorld)
	q.componentIdN = ComponentIdFor[N](targetWorld)
	q.componentIdO = ComponentIdFor[O](targetWorld)
	q.componentIdP = ComponentIdFor[P](targetWorld)
	q.components = []ComponentId{
		q.componentIdA, q.componentIdB, q.componentIdC, q.componentIdD, q.componentIdE, q.componentIdF, q.componentIdG, q.componentIdH,
		q.componentIdI, q.componentIdJ, q.componentIdK, q.componentIdL, q.componentIdM, q.componentIdN, q.componentIdO, q.componentIdP,
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
func fetchComponentForQueryResult[T IComponent](componentId ComponentId, entityRow uint, archetype *Archetype, queryOptions *CombinedQueryOptions) (result *T, err error) {
	storage := archetype.components[componentId]
	result, err = getComponentFromComponentStorage[T](storage, entityRow)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve component %s from storage: %v", componentId.DebugString(), err)
	}

	if result != nil && (queryOptions.ReadOnlyComponents.IsAllReadOnly || slices.Contains(queryOptions.ReadOnlyComponents.ComponentIds, componentId)) {
		result = utils.ClonePointerValue(result)
	}

	return result, nil
}
