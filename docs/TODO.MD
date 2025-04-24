# ECS
**Must have**
- [feature] implement query filters And / Or
- [performance] Proper archetype ECS implementation. Once all of the above is implemented, the basic user facing API is there. But the implementation is its most basic form and is not efficient at all. [See this article from flecs ECS creator](https://ajmmertens.medium.com/building-an-ecs-1-where-are-my-entities-and-components-63d07c7da742).

**Nice-to-have**
- [quality-of-life] Automatically execute queries from system params before calling the system. We could then also have a LazyQuery that only gets called when getting query results.
- [code-cleanup] Check if files can be split up in to 3 packages: storage, actions, filter
- [performance] Generational entities
- [performance] When removing component from entity, or deleting an entity, the place of the component(s) in the component registry are never reused. We could move the last componentRegistry component to the place of the removed component to reuse the memory and prevent memory from forever growing. When moving a component in componentRegistry, Don't forget to update the index in entityData
- [performance] Run schedule systems in parallel. We could have the query params use Q[C] for immutable read and Q[*C] for mutable read. We can then parallelize all systems in the same schedule that do not query a component that is mutably queried by another system.
- [feature] Relationships (like parent/child)
- [feature] Events
- [feature] Observers
- [feature] World configs (logger, default componentEntry size)
- [tests] More realistic ECS benchmarks. Check out [this benchmarks page for Go ECS's](https://github.com/mlange-42/go-ecs-benchmarks)
- [quality-of-life] Better query validations:
    - Return error if marking a component as optional while that component is not queried
    - Return error if there are any duplicate filters

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

# General
- [code-cleanup] use reflect.TypeFor instead of utils.TypeOf. But make sure to benchmark it to see how performant it is 
