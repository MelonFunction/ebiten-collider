// Package collider 👍
package collider

import (
	"errors"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Vars
var (
	ErrShapeNotFound = errors.New("Couldn't remove shape from SpatialHash; not found")
)

// Shape interface. It's probably not needed but it keeps code more readable.
type Shape interface {
	GetPosition() *Point // get the position
	GetBounds() (float64, float64, float64, float64)
	Move(x, y float64)      // move by amount
	MoveTo(x, y float64)    // move to position
	SetHash(s *SpatialHash) // sets ref to hash
	GetHash() *SpatialHash  // gets	 ref to hash
}

// Point is a point in space with X and Y positions
type Point struct {
	X, Y float64
}

// PointShape shape, also used for the hash's cell coordinates
type PointShape struct {
	Pos         *Point
	SpatialHash *SpatialHash
}

// CircleShape shape
type CircleShape struct {
	// Center point
	Pos         *Point
	Radius      float64
	SpatialHash *SpatialHash
}

// RectangleShape shape
type RectangleShape struct {
	// Center point
	Pos           *Point
	Width, Height float64
	SpatialHash   *SpatialHash
}

// Cell contains shapes
type Cell struct {
	Shapes map[Shape]Shape
}

// CellCoord is a coordinate used to index the hash
type CellCoord struct {
	X, Y int
}

// SpatialHash contains cells
type SpatialHash struct {
	// Size of the grid/cell/partition
	CellSize int
	// Store shapes in a cell depending on their bounds
	Hash map[CellCoord]*Cell
	// Backref for shapes to find its containing cells
	Backref map[Shape][]*Cell
}

// NewSpatialHash returns a new *SpatialHash
func NewSpatialHash(cellSize int) *SpatialHash {
	return &SpatialHash{
		CellSize: cellSize,
		Hash:     make(map[CellCoord]*Cell),
		Backref:  make(map[Shape][]*Cell),
	}
}

// Add adds a shape to the spatial hash
func (s *SpatialHash) Add(shape Shape) {
	x1, y1, x2, y2 := shape.GetBounds()

	// make sure big shapes are constrained properly
	xStep := x2 - x1
	if xStep > float64(s.CellSize) {
		xStep = xStep / float64(s.CellSize)
	}
	yStep := y2 - y1
	if yStep > float64(s.CellSize) {
		yStep = yStep / float64(s.CellSize)
	}
	for x := x1; x <= x2; x += xStep {
		for y := y1; y <= y2; y += yStep {
			hashPos := CellCoord{
				int(math.Floor(x / float64(s.CellSize))),
				int(math.Floor(y / float64(s.CellSize))),
			}
			if _, ok := s.Hash[hashPos]; !ok {
				s.Hash[hashPos] = &Cell{Shapes: make(map[Shape]Shape)}
			}
			s.Hash[hashPos].Shapes[shape] = shape                        // add shape to cell
			s.Backref[shape] = append(s.Backref[shape], s.Hash[hashPos]) // add cell to backref
		}
	}

	shape.SetHash(s)
}

// Remove removes a shape from the spatial hash
func (s *SpatialHash) Remove(shape Shape) error {
	if cells, ok := s.Backref[shape]; ok {
		for _, cell := range cells {
			delete(cell.Shapes, shape)
		}
		s.Backref[shape] = nil
	}

	return ErrShapeNotFound
}

// GetCollisionCandidates returns a list of all shapes in the same cells as shape
func (s *SpatialHash) GetCollisionCandidates(shape Shape) []Shape {
	shapesMap := make(map[Shape]bool)
	if cells, ok := s.Backref[shape]; ok {
		for _, cell := range cells {
			for _, sh := range cell.Shapes {
				shapesMap[sh] = true
			}
		}
	}
	delete(shapesMap, shape)
	shapes := make([]Shape, len(shapesMap))
	for k := range shapesMap {
		shapes = append(shapes, k)
	}
	return shapes
}

// CollisionData contains information about the collision
type CollisionData struct {
	Other            Shape
	SeparatingVector Point
}

// CheckCollisions returns a list of all shapes and their separating vector
func (s *SpatialHash) CheckCollisions(shape Shape) []CollisionData {
	collisions := make([]CollisionData, 0)
	candidates := s.GetCollisionCandidates(shape)

	for _, candidate := range candidates {
		switch other := candidate.(type) {
		case *PointShape:
			_ = other
		case *RectangleShape:
		case *CircleShape:
		default:
			// TODO error
		}

	}

	return collisions
}

// Draw is a debug function. It draws a rectangle for every cell which has had a shape in it at some point.
func (s *SpatialHash) Draw(surface *ebiten.Image) {
	for pos, cell := range s.Hash {
		x, y, w := float64(pos.X*s.CellSize), float64(pos.Y*s.CellSize), float64(s.CellSize)
		color := color.RGBA{255, 255, 255, 255}
		ebitenutil.DrawLine(surface, x, y, x+w, y, color)
		ebitenutil.DrawLine(surface, x, y, x, y+w, color)
		ebitenutil.DrawLine(surface, x, y+w, x+w, y+w, color)
		ebitenutil.DrawLine(surface, x+w, y, x+w, y+w, color)

		ebitenutil.DebugPrintAt(surface, fmt.Sprintf("%d", len(cell.Shapes)), pos.X*s.CellSize, pos.Y*s.CellSize)
	}
}

// NewPointShape creates, then adds a new PointShape to the hash before returning it
func (s *SpatialHash) NewPointShape(x, y float64) *PointShape {
	po := &PointShape{
		Pos: &Point{x, y},
	}
	s.Add(po)
	return po
}

// GetPosition returns the Point of the PointShape
func (po *PointShape) GetPosition() *Point {
	return po.Pos
}

// GetBounds returns the Bounds of the PointShape
func (po *PointShape) GetBounds() (float64, float64, float64, float64) {
	return po.Pos.X - 0.5, po.Pos.Y - 0.5, po.Pos.X + 0.5, po.Pos.Y + 0.5
}

// Move moves the PointShape by x and y
func (po *PointShape) Move(x, y float64) {
	po.Pos.X += x
	po.Pos.Y += y
	hash := po.GetHash()
	hash.Remove(po)
	hash.Add(po)
}

// MoveTo moves the PointShape to x and y
func (po *PointShape) MoveTo(x, y float64) {
	po.Pos.X = x
	po.Pos.Y = y
	hash := po.GetHash()
	hash.Remove(po)
	hash.Add(po)
}

// SetHash sets the hash
func (po *PointShape) SetHash(s *SpatialHash) {
	po.SpatialHash = s
}

// GetHash gets the hash
func (po *PointShape) GetHash() *SpatialHash {
	return po.SpatialHash
}

// NewCircleShape creates, then adds a new CircleShape to the hash before returning it
func (s *SpatialHash) NewCircleShape(x, y, r float64) *CircleShape {
	ci := &CircleShape{
		Pos:    &Point{x, y},
		Radius: r,
	}
	s.Add(ci)
	return ci
}

// GetPosition returns the Point of the CircleShape
func (ci *CircleShape) GetPosition() *Point {
	return ci.Pos
}

// GetBounds returns the Bounds of the CircleShape
func (ci *CircleShape) GetBounds() (float64, float64, float64, float64) {
	return ci.Pos.X - ci.Radius,
		ci.Pos.Y - ci.Radius,
		ci.Pos.X + ci.Radius,
		ci.Pos.Y + ci.Radius
}

// Move moves the CircleShape by x and y
func (ci *CircleShape) Move(x, y float64) {
	ci.Pos.X += x
	ci.Pos.Y += y
	hash := ci.GetHash()
	hash.Remove(ci)
	hash.Add(ci)
}

// MoveTo moves the CircleShape to x and y
func (ci *CircleShape) MoveTo(x, y float64) {
	ci.Pos.X = x
	ci.Pos.Y = y
	hash := ci.GetHash()
	hash.Remove(ci)
	hash.Add(ci)
}

// SetHash sets the hash
func (ci *CircleShape) SetHash(s *SpatialHash) {
	ci.SpatialHash = s
}

// GetHash gets the hash
func (ci *CircleShape) GetHash() *SpatialHash {
	return ci.SpatialHash
}

// NewRectangleShape creates, then adds a new RectangleShape to the hash before returning it
func (s *SpatialHash) NewRectangleShape(x, y, w, h float64) *RectangleShape {
	sq := &RectangleShape{
		Pos:    &Point{x, y},
		Width:  w,
		Height: h,
	}
	s.Add(sq)
	return sq
}

// GetPosition returns the Point of the RectangleShape
func (sq *RectangleShape) GetPosition() *Point {
	return sq.Pos
}

// GetBounds returns the Bounds of the RectangleShape
func (sq *RectangleShape) GetBounds() (float64, float64, float64, float64) {
	return sq.Pos.X - sq.Width/2,
		sq.Pos.Y - sq.Height/2,
		sq.Pos.X + sq.Width/2,
		sq.Pos.Y + sq.Height/2
}

// Move moves the RectangleShape by x and y
func (sq *RectangleShape) Move(x, y float64) {
	sq.Pos.X += x
	sq.Pos.Y += y
	hash := sq.GetHash()
	hash.Remove(sq)
	hash.Add(sq)
}

// MoveTo moves the RectangleShape to x and y
func (sq *RectangleShape) MoveTo(x, y float64) {
	sq.Pos.X = x
	sq.Pos.Y = y
	hash := sq.GetHash()
	hash.Remove(sq)
	hash.Add(sq)
}

// SetHash sets the hash
func (sq *RectangleShape) SetHash(s *SpatialHash) {
	sq.SpatialHash = s
}

// GetHash gets the hash
func (sq *RectangleShape) GetHash() *SpatialHash {
	return sq.SpatialHash
}
