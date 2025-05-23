# ECS
**Must have**
- [feature] Observers
- [feature] Relationships (like parent/child)

**Nice-to-have**
- [performance] Generational entities. EntityId would become {id: uint, generation: uint}. World its entityList would becomes a slice instead of a map, and entityId would correlate with the slice index. This provides faster lookups. We'd need an deadEntityPool to choose from when creating a new EntityId which would increment its generation and repurpose it.
- [performance] Archetype graph to speed up archetype moves, making insert/remove more efficient. Add an `edges map[componentIds]Archetype` to Archetype where `componentIds` is a hash of the components to add or remove. This gives better performance for archetypes with many components because only the components to add/remove need to be hashed instead of the whole set of new component ids. See https://ajmmertens.medium.com/building-an-ecs-2-archetypes-and-vectorization-fe21690805f9 for more explanation and examples.
- [performance] Query can create a list of ArcheTypes on Prepare so that we don't have to go through every entity. Archetypes would have to maintain a list of entities. We'd also have to have some kind of dirty flag for archetypes so that if a new archetype is created, the query updates its list of archetypes.
- [performance] Cache Queries
- [feature] Events
- [feature] Add Query5..Query16 and Optional5..Optional16
- [tests] More realistic ECS benchmarks. Check out [this benchmarks page for Go ECS's](https://github.com/mlange-42/go-ecs-benchmarks)
- [quality-of-life] handle spawning/inserting nil

# Project
- [feature] Linter
- [feature] Pipeline in Github that automatically runs:
    - tests - fail if any fails
    - linter - fail if project is not linted
    - benchmarks - warn if there are degradations
    - run `go mod tidy` - fail if anything changed

# App
**Must have**
- [feature] Query between different App worlds
    1. Add system param `QueryOtherWorld<AppId, Query>`. Apps would need to know about each others worlds, either automatically or by having a function to manually add them.
    2. During execution of QueryOtherWorld, prevent any systems from executing in the other world.

**Nice-to-have**
- [feature] Execute systems in parallel. If systems are in the same schedule, they can be run in parallel if the systems do not touch any passed-by-ref resource (query/resource) of any of the other systems. Mind that we can probably not run a system with a Query in parallel with any other system that pulls in ecs.World, because we can not detect if it will be used to mutate anything that other systems touch system-insert time.
