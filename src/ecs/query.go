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

func (q *Query0[QueryOptions]) ClearResults() {
	q.Clear()
}
func (q *Query1[ComponentA, QueryOptions]) ClearResults() {
	q.Clear()
}
func (q *Query2[ComponentA, ComponentB, QueryOptions]) ClearResults() {
	q.Clear()
}
func (q *Query3[ComponentA, ComponentB, ComponentC, QueryOptions]) ClearResults() {
	q.Clear()
}
func (q *Query4[ComponentA, ComponentB, ComponentC, ComponentD, QueryOptions]) ClearResults() {
	q.Clear()
}

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
