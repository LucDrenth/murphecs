# About
This document states some of the design decisions that were made for this project. If you have any questions, feedback or suggestions on any of these, feel free to reach out.

# Referencing other projects
I will be comparing the approaches I've taken with other similar projects. While this at some points may come of as a negative, it is not my intent to do so. On the contrary, I am very grateful for the examples that other frameworks provide. 

I will be linking to any other project I mention for the sake of easily comparing. If you are the owner of any of such projects and have any objection to this document mentioning your project, please reach out and I will remove the reference of your project from this document.

# Design Decisions

## ECS - spawning components
For the spawning API, I have chosen to go with a variadic parameter of components without using generics. 
```Go
func (world *world) Spawn(components ...IComponent)
```
It can be used like this:
```Go
world.Spawn(Position{...}, Velocity{...})

// If we want to add another component, simply add it as a parameter
world.Spawn(Position{...}, Velocity{...}, Color{...})
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
3. We need to pass the world as a parameter, because Golang does not allow generics to be used with methods. Calling Spawn as a method of world is simpler. However, many of the rest of the ECS API does use generics, which leads to inconsistency in the API.
4. We'd have to create a separate spawn function for each number of components we might possibly add. Go does not currently implement variadic generics.

I have not looked in to the performance difference of these two approaches yet. But in my opinion, the much cleaner API makes up for any possible performance improvement that the generics approach might have.
