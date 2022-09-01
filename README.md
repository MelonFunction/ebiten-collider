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

// Set arbitrary data. Useful for linking a shape to another struct.
wall.SetParent("I'm a wall")
// or...
wall.SetParent(Tile{Type:"Wall"})
// GetParent during a collision, check it's type and then change how the collision is handled. 
// If it returns data which represents a wall, then you should move by the collision.SeparatingVector, but if it's a 
// floor tile, then it can be ignored.
log.Println(wall.GetParent())

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