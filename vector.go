package collider

import "math"

// Vector represents a point in space
type Vector struct {
	X, Y float64
}

// NewVector returns a new *Vector
func NewVector(x, y float64) *Vector {
	return &Vector{x, y}
}

// Mult multiplies v by scaler and returns a new Vector for chaining
func (v *Vector) Mult(scalar float64) *Vector {
	return &Vector{
		v.X * scalar,
		v.Y * scalar,
	}
}

// Add adds o to v and returns a new Vector
func (v *Vector) Add(o *Vector) *Vector {
	return &Vector{
		v.X + o.X,
		v.Y + o.Y,
	}
}

// Sub subtracts o from v and returns a new Vector
func (v *Vector) Sub(o *Vector) *Vector {
	return &Vector{
		v.X - o.X,
		v.Y - o.Y,
	}
}

// Length returns the length of the vector
func (v *Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize returns a new Vector representing the normal of v
func (v *Vector) Normalize() *Vector {
	l := v.Length()
	if l > 0 {
		return &Vector{
			v.X / l,
			v.Y / l,
		}
	}
	return NewVector(v.X, v.Y)
}

// Rotate rotates a point about 0,0 and returns v for chaining
func (v *Vector) Rotate(phi float64) *Vector {
	c, s := math.Cos(phi), math.Sin(phi)
	return &Vector{
		X: c*v.X - s*v.Y,
		Y: s*v.X + c*v.Y,
	}
}

// RotateAround rotates a Vector about another Vector and returns v for chaining
func (v *Vector) RotateAround(phi float64, o *Vector) *Vector {
	c, s := math.Cos(phi), math.Sin(phi)
	n := NewVector(v.X, v.Y).Sub(o)
	return NewVector(
		c*n.X-s*n.Y,
		s*n.X+c*n.Y,
	).Add(o)
}

// AngleTo returns the angle between v to other
func (v *Vector) AngleTo(other *Vector) float64 {
	return math.Atan2(v.Y-other.Y, v.X-other.X)
}
