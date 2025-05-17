# ECS
**Must have**
- [bug] Fix GC issue when copying pointer values. See https://github.com/mlange-42/ark/issues/246
- [feature] Observers
- [feature] Relationships (like parent/child)
- [performance] When removing component from entity, or deleting an entity, the place of the component(s) in the component storage are never reused. To prevent memory from forever growing, We could either move the last componentRegistry component to the place of the removed component, but mind that we'd need to update it for all storages in the archetype

**Nice-to-have**
- [performance] Generational entities
- [performance] Archetype graph to speed up insert/remove. Add an `edges map[componentIds]Archetype` to Archetype where `componentIds` is a hash of the components to add or remove. This gives better performance for archetypes with many components because only the components to add/remove need to be hashed instead of the whole set of new component ids. See https://ajmmertens.medium.com/building-an-ecs-2-archetypes-and-vectorization-fe21690805f9 for more explanation and examples.
- [performance] Query can create a list of ArcheTypes on Prepare so that we don't have to go through every entity. Archetypes would have to maintain a list of entities.
- [feature] Events
- [feature] Add Query5..Query16 and Optional5..Optional16
- [feature] World config on how to grow component storage capacity
- [quality-of-life] Function to check if a specific entity has components. This is already possible by using Get1..Get8 and checking if the error is ErrComponentNotFound, but its not very user-friendly.
- [quality-of-life] Add Query0 to only get entities. This is already possible be specifying some random component and marking it as optional, but this is not very user-friendly an not as performant as it could be. 
- [quality-of-life] Set function to either overwrite or insert a component. It is already possible by removing and then inserting a component but its not very user friendly.
- [tests] More realistic ECS benchmarks. Check out [this benchmarks page for Go ECS's](https://github.com/mlange-42/go-ecs-benchmarks)

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

**Nice-to-have**
- [feature] Execute systems in parallel. If systems are in the same schedule, they can be run in parallel if the systems do not touch any passed-by-ref resource (query/resource) of any of the other systems. Mind that we can probably not run a system with a Query in parallel with any other system that pulls in ecs.World, because we can not detect if it will be used to mutate anything that other systems touch system-insert time.
