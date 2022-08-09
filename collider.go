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

// Shape base
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
	X, Y        float64
	SpatialHash *SpatialHash
}

// CircleShape shape
type CircleShape struct {
	// Center point
	Pos         *Point
	Radius      float64
	SpatialHash *SpatialHash
}

// SquareShape shape
type SquareShape struct {
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

// NewSquareShape creates, then adds a new SquareShape to the hash before returning it
func (s *SpatialHash) NewSquareShape(x, y, w, h float64) *SquareShape {
	sq := &SquareShape{
		Pos:    &Point{x, y},
		Width:  w,
		Height: h,
	}
	s.Add(sq)
	return sq
}

// GetPosition returns the Point of the SquareShape
func (sq *SquareShape) GetPosition() *Point {
	return sq.Pos
}

// GetBounds returns the Bounds of the SquareShape
func (sq *SquareShape) GetBounds() (float64, float64, float64, float64) {
	return sq.Pos.X - sq.Width/2,
		sq.Pos.Y - sq.Height/2,
		sq.Pos.X + sq.Width/2,
		sq.Pos.Y + sq.Height/2
}

// Move moves the SquareShape by x and y
func (sq *SquareShape) Move(x, y float64) {
	sq.Pos.X += x
	sq.Pos.Y += y
	hash := sq.GetHash()
	hash.Remove(sq)
	hash.Add(sq)
}

// MoveTo moves the SquareShape to x and y
func (sq *SquareShape) MoveTo(x, y float64) {
	sq.Pos.X = x
	sq.Pos.Y = y
	hash := sq.GetHash()
	hash.Remove(sq)
	hash.Add(sq)
}

// SetHash sets the hash
func (sq *SquareShape) SetHash(s *SpatialHash) {
	sq.SpatialHash = s
}

// GetHash gets the hash
func (sq *SquareShape) GetHash() *SpatialHash {
	return sq.SpatialHash
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

// GetBounds returns the Bounds of the SquareShape
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
