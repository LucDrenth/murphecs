# ECS
**Must have**
- Update readme with querying example
- Put the S in to the ECS (implement systems)
- Proper archetype ECS implementation. Once all of the above is implemented, the basic user facing API is there. But the implementation is its most basic form and is not efficient at all. [See this article from flecs ECS creator](https://ajmmertens.medium.com/building-an-ecs-1-where-are-my-entities-and-components-63d07c7da742).

**Nice-to-have**
- Check if files can be split up in to 3 packages: storage, actions, filter
- Relationships (like parent/child)
- Events
- Observers
- More realistic ECS benchmarks. Check out [this benchmarks page for Go ECS's](https://github.com/mlange-42/go-ecs-benchmarks)
- Better query validations:
    - Return error if marking a component as optional while that component is not queried
    - Return error if there are any duplicate filters

# Logging
**Must have**
- Logger Interface with:
    - info
    - warn
    - error
- Logger implementation that logs to console

**Nice-to-have**
- Methods to log something only once. Maybe we can use caller location for this?

# Window
- Be able to create a window

# Graphics
- Draw a triangle

# Project
- Linter
- Pipeline in Github that automatically runs:
    - tests - fail if any fails
    - linter - fail if project is not linted
    - benchmarks - warn if there are degradations
    - run `go mod tidy` - fail if anything changed
