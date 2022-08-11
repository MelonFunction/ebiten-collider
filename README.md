# ebiten-collider

Basic collision detection optimized with the use of spatial partitioning.  

It can detect and resolve the following collision pairs:  
- Rectangle + Rectangle
- Rectangle + Circle
- Rectangle + Point
- Circle + Circle
- Circle + Point
- Point + Point

## Usage

[📖 Docs](https://pkg.go.dev/github.com/melonfunction/ebiten-collider)  
Look at [the example](https://github.com/melonfunction/ebiten-collider/tree/master/examples) to see how to use the library.

```go
// Create the world
hash = collider.NewSpatialHash(128)

// Create a shape
wall = hash.NewRectangleShape(
    0
    16,
    128,
    128,
)

// Check for collisions
collisions := hash.CheckCollisions(player)
for _, collision := range collisions {
    sep := collision.SeparatingVector
    // Move a shape by the overlap
    player.Move(sep.X, sep.Y)

    // Or move both shapes equally
    player.Move(sep.X/2, sep.Y/2)
    collision.Other.Move(-sep.X/2, -sep.Y/2)
}
```