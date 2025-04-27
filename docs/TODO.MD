# ECS
**Must have**
- [feature] Implement query filters And / Or
- [performance] Proper archetype ECS implementation. Once all of the above is implemented, the basic user facing API is there. But the implementation is its most basic form and is not efficient at all. [See this article from flecs ECS creator](https://ajmmertens.medium.com/building-an-ecs-1-where-are-my-entities-and-components-63d07c7da742).

**Nice-to-have**
- [quality-of-life] Automatically execute queries from system params before calling the system. We could then also have a LazyQuery that only gets called when getting query results.
- [quality-of-live] Add Query structs to only get 1 or maybe 2 options, such as QueryWithOptional, QueryWithout, QueryReadOnly and such. This is so that we don't have to write out all query options in every query.
- [code-cleanup] Check if files can be split up in to 3 packages: storage, actions, filter
- [performance] Generational entities
- [performance] When removing component from entity, or deleting an entity, the place of the component(s) in the component registry are never reused. We could move the last componentRegistry component to the place of the removed component to reuse the memory and prevent memory from forever growing. When moving a component in componentRegistry, Don't forget to update the index in entityData
- [feature] Relationships (like parent/child)
- [feature] Events
- [feature] Observers
- [feature] World configs (logger, default componentEntry size)
- [feature] Add Query5..Query16 and Optional5..Optional16
- [quality-of-life] Better query validations:
    - Return error if marking a component as optional while that component is not queried
    - Return error if there are any duplicate filters
- [quality-of-life] Function to check if a specific entity has components. This is already possible by using Get1..Get8 and checking if the error is ErrComponentNotFound, but its not very user-friendly.
- [tests] More realistic ECS benchmarks. Check out [this benchmarks page for Go ECS's](https://github.com/mlange-42/go-ecs-benchmarks)

# Window
- [feature] Be able to create a window

# Graphics
- [feature] Draw a triangle

# Project
- [feature] Linter
- [feature] Pipeline in Github that automatically runs:
    - tests - fail if any fails
    - linter - fail if project is not linted
    - benchmarks - warn if there are degradations
    - run `go mod tidy` - fail if anything changed

# App
**Must have**
- [feature] Dynamically insert resource (maybe name it pullable?) and convert logger to that.
- [feature] Loop that can run at a fixed rate
- [feature] Delta resource

**Nice-to-have**
- [feature] Resource without pointer to make a copy.
- [feature] Execute systems in parallel. If systems are in the same schedule, they can be run in parallel if the systems do not touch any passed-by-ref resource (query/resource) of any of the other systems. Mind that we can probably not run a system with a Query in parallel with any other system that pulls in ecs.World, because we can not detect if it will be used to mutate anything that other systems touch system-insert time.
