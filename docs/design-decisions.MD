# About
This document states some of the design decisions that were made for this project. If you have any questions, feedback or suggestions on any of these, feel free to reach out.

# Referencing other projects
I will be comparing the approaches I've taken with other similar projects. While this at some points may come of as a negative, it is not my intent to do so. On the contrary, I am very grateful for the examples that other frameworks provide. 

I will be linking to any other project I mention for the sake of easily comparing. If you are the owner of any of such projects and have any objection to this document mentioning your project, please reach out and I will remove the reference of your project from this document.

# Design Decisions

## ECS - spawning components with variadic
For the spawning API, I have chosen to go with a variadic parameter of components without using generics. 
```Go
func (world *world) Spawn(components ...IComponent)
```
It can be used like this:
```Go
ecs.Spawn(&world, Position{...}, Velocity{...})

// If we want to add another component, simply add it as a parameter
ecs.Spawn(&world, Position{...}, Velocity{...}, Color{...})
```

I think this provides the cleanest spawning API that is currently possible in Go (1.24.2). When using generics, like implemented in [Ark](https://github.com/mlange-42/ark), it would look something like this:
```Go
ecs.Spawn2[Position, Velocity](&world, Position{...}, Velocity{...})

// If we now want to add another entity, we'd have to write:
ecs.Spawn3[Position, Velocity, Color](&world, Position{...}, Velocity{...}, Color{...})
```

Comparing this to the variadic approach, the generic approach has a couple of major downsides:
1. We have to pass both the component type (as generic) and the component itself. With the variadic approach, we only have to add the component.
2. If we'd want to add another component (which is frequently needed during development), we'd have to update 3 things: the function name, the component types (as generics) and finally adding the actual component. With the variadic implementation, we can simply add the component.
3. We'd have to create a separate spawn function for each number of components we might possibly add. Go does not currently implement variadic generics.

I have not looked in to the performance difference of these two approaches yet. But in my opinion, the much cleaner API makes up for any possible performance improvement that the generics approach might have.

## ECS - functions instead of world methods
Some world-related functions, like `Spawn`, could have been implemented as methods of `world`. However, some world-related functions can not be implemented as methods because they involve generics, and Go (1.24.2) does not support this methods with generics.

To keep the ECS API consistent, all major user-facing ECS functionality is implemented with functions.

## ECS - no subdirectories
There are no subdirectories in the ECS package, and it has many files. I have decided on doing it this way so that the API is simple to use: the ECS API can only be used through `ecs.*`. 

If we were to split the ecs package up in to multiple directories (something like ecs/storage, ecs/action, ecs/query etc.), the ECS functions would be scattered. This complicates the API. Users would have to remember which package contains which part, like `storage.World`, `query.Filter` and so on. We could get around this by exposing everything in the ecs package using type aliases, but then the API could still be used in two ways, and seeing what's underneath the exposed/aliased type would require an extra click.
