# ECS
**Must have**
- [feature] Observers

**Nice-to-have**
- [performance] Generational entities. EntityId would become {id: uint, generation: uint}. World its entityList would becomes a slice instead of a map, and entityId would correlate with the slice index. This provides faster lookups. We'd need an deadEntityPool to choose from when creating a new EntityId which would increment its generation and repurpose it.
- [performance] Archetype graph to speed up archetype moves, making insert/remove more efficient. Add an `edges map[componentIds]Archetype` to Archetype where `componentIds` is a hash of the components to add or remove. This gives better performance for archetypes with many components because only the components to add/remove need to be hashed instead of the whole set of new component ids. See https://ajmmertens.medium.com/building-an-ecs-2-archetypes-and-vectorization-fe21690805f9 for more explanation and examples.
- [performance] Reduce the number of archetypes that queries go through on calls to Exec. There are multiple possible approaches:
    1. Each query stores a list of archetypes. We'd have to have some kind of dirty flag for archetypes so that if a new archetype is created, the query updates its list of archetypes. The downside of this approach is that complex applications that have archetype moves every frame will not benefit from this. We could get around this by having a 'smart' dirty flag system that only marks certain components as dirty.
    2. Each query stores a list of possible combinations of component hashes and only checks archetypes with those hashes (using world.archetypeStorage.componentsHashToArchetype). The downside of this approach is that this list could potentially get very large.
    3. Each query stores a list of all involved components (queried + from filters) and gets a union of archetypes from `world.archetypeStorage.componentIdToArchetypes`. We could initially do this every Exec call and add cache it in a later step (using dirty flags for archetypes when new archetypes are created).
- [performance] Cache Queries
- [feature] Relationships (like parent/child)
- [tests] More realistic ECS benchmarks. Check out [this benchmarks page for Go ECS's](https://github.com/mlange-42/go-ecs-benchmarks)

# Project
- [feature] Pipeline in Github that automatically runs:
    - tests - fail if any fails
    - linter - fail if project is not linted
    - benchmarks - warn if there are degradations
    - run `go mod tidy` - fail if anything changed

# App
**Nice-to-have**
- [feature] Execute systems in parallel. If systems are in the same schedule, they can be run in parallel if the systems do not touch any passed-by-ref resource (query/resource) of any of the other systems. Mind that we can probably not run a system with a Query in parallel with any other system that pulls in ecs.World, because we can not detect if it will be used to mutate anything that other systems touch system-insert time.
- [performance] When executing systems, not all outer worlds need to be locked for the whole run. We actually only need to lock some components (or the archetypes of those components?) for some worlds.
