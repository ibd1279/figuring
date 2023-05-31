package figuring

import (
	"fmt"
	"math"
	"sort"

	"github.com/go-gl/mathgl/mgl64"
)

// Pair is the interface that of all types that allow retreiving the underlying
// units as a pair.
type Pair interface {
	Units() (Length, Length)
}

var (
	// VectorZero is a Vector of zero magnitude.
	VectorZero = VectorIj(0, 0)

	// VectorUnit is a Vector of 1 unit length (1 micrometer).
	VectorUnit = VectorIj(1, 1)

	// VectorNaN is a Vector in error.
	VectorNaN = VectorIj(Length(math.NaN()), Length(math.NaN()))

	// PtOrig is the Origin Point.
	PtOrig = PtXy(0, 0)

	// PtNaN is a Pt in error.
	PtNaN = PtXy(Length(math.NaN()), Length(math.NaN()))
)

// Quadrant represents which direction a vector points towards.
type Quadrant uint

const (
	VECTOR_DIRECTION_NONE Quadrant = iota
	VECTOR_DIRECTION_Q1
	VECTOR_DIRECTION_Q2
	VECTOR_DIRECTION_Q3
	VECTOR_DIRECTION_Q4
	VECTOR_DIRECTION_XPOS
	VECTOR_DIRECTION_YPOS
	VECTOR_DIRECTION_XNEG
	VECTOR_DIRECTION_YNEG
)

// Pt represents an x,y value on a 2d plane.
type Pt struct {
	xy mgl64.Vec2
}

// PtAt with a given x and y value.
func PtXy(x, y Length) Pt {
	xy := mgl64.Vec2{float64(x), float64(y)}
	return PtFromVec2(xy)
}

// PtFromVec2 creates a points from a vec2. Mostly used internally.
func PtFromVec2(v mgl64.Vec2) Pt {
	return Pt{xy: v}
}

// X returns the X coordinate.
func (p Pt) X() Length {
	x, _ := p.Units()
	return x
}

// Y returns the Y coordinate.
func (p Pt) Y() Length {
	_, y := p.Units()
	return y
}

// XY returns the x and y coordinate. Semantic shorthand for Units().
func (p Pt) XY() (Length, Length) {
	return p.Units()
}

// Units implements Pair Interface.
func (p Pt) Units() (Length, Length) {
	return Length(p.xy[0]), Length(p.xy[1])
}

// OrErr tests if either coordinate is NaN or Inf and returns an error if one
// is. NaN errors are prioritized over Inf errors.
func (p Pt) OrErr() (Pt, *FloatingPointError) {
	x, y := p.Units()
	_, xerr := x.OrErr()
	_, yerr := y.OrErr()
	if xerr != nil && xerr.IsNaN() {
		return p, xerr
	} else if yerr != nil && yerr.IsNaN() {
		return p, yerr
	} else if xerr != nil {
		return p, xerr
	} else if yerr != nil {
		return p, yerr
	}
	return p, nil
}

// String outputs the points coordinates.
func (p Pt) String() string {
	return fmt.Sprintf("Point({%s, %s})",
		HumanFormat(9, p.xy[0]),
		HumanFormat(9, p.xy[1]))
}

// Add \c b to \c p to get a new Pt.
func (p Pt) Add(b Vector) Pt {
	xy := mgl64.Vec2{p.xy[0] + b.ij[0], p.xy[1] + b.ij[1]}
	return PtFromVec2(xy)
}

// VectorTo creates the vector from \c p to \c b. Use PtOrig.VectorTo(p) in
// order to get the vector for an arbitrary Pt.
func (p Pt) VectorTo(b Pt) Vector {
	ij := mgl64.Vec2{b.xy[0] - p.xy[0], b.xy[1] - p.xy[1]}
	return VectorFromVec2(ij)
}

type ptSlice []Pt

func (x ptSlice) Len() int { return len(x) }
func (x ptSlice) Less(i, j int) bool {
	if x[i].X() < x[j].X() || (x[i].X() != x[i].X() && x[j].X() == x[j].X()) {
		return true
	} else if x[i].Y() < x[j].Y() || (x[i].Y() != x[i].Y() && x[j].Y() == x[j].Y()) {
		return true
	}
	return false
}
func (x ptSlice) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

// Vector represents a direction and a magnitude.
// See https://scholarsarchive.byu.edu/cgi/viewcontent.cgi?article=1000&context=facpub
type Vector struct {
	ij        mgl64.Vec2
	direction Quadrant
}

func VectorIj(i, j Length) Vector {
	ij := mgl64.Vec2{float64(i), float64(j)}
	return VectorFromVec2(ij)
}

func VectorFromVec2(ij mgl64.Vec2) Vector {
	// Treat numbers really close to zero as zero.
	if IsZero(ij[0]) {
		ij[0] = 0
	}
	if IsZero(ij[1]) {
		ij[1] = 0
	}

	// detect the direction.
	var d Quadrant
	switch {
	case ij[0] > 0 && ij[1] > 0:
		d = VECTOR_DIRECTION_Q1
	case ij[0] < 0 && ij[1] > 0:
		d = VECTOR_DIRECTION_Q2
	case ij[0] < 0 && ij[1] < 0:
		d = VECTOR_DIRECTION_Q3
	case ij[0] > 0 && ij[1] < 0:
		d = VECTOR_DIRECTION_Q4
	case ij[0] > 0 && ij[1] == 0:
		d = VECTOR_DIRECTION_XPOS
	case ij[0] < 0 && ij[1] == 0:
		d = VECTOR_DIRECTION_XNEG
	case ij[0] == 0 && ij[1] > 0:
		d = VECTOR_DIRECTION_YPOS
	case ij[0] == 0 && ij[1] < 0:
		d = VECTOR_DIRECTION_YNEG
	case ij[0] == 0 && ij[1] == 0:
		d = VECTOR_DIRECTION_NONE
	default:
		d = VECTOR_DIRECTION_NONE
	}

	// precalculate the magnitude
	return Vector{
		ij:        ij,
		direction: d,
	}
}

// VectorFromTheta returns a unit vector pointed in the direction of the provided theta.
func VectorFromTheta(theta Radians) Vector {
	ij := mgl64.Vec2{math.Cos(float64(theta)), math.Sin(float64(theta))}
	return VectorFromVec2(ij)
}

// Magnitude returns the combined distance of this vector.
func (v Vector) Magnitude() Length {
	return Length(math.Hypot(v.ij[0], v.ij[1]))
}

// Angle returns the angle of the vector with positive x-axis as 0,
// increasing anti-clockwise.
func (v Vector) Angle() Radians {
	// Handles the axis directions to avoid any divide by zero
	// or invalid values for arctan
	switch v.direction {
	case VECTOR_DIRECTION_XPOS:
		return 0.
	case VECTOR_DIRECTION_YPOS:
		return math.Pi * 0.5
	case VECTOR_DIRECTION_XNEG:
		return math.Pi
	case VECTOR_DIRECTION_YNEG:
		return math.Pi * 1.5
	}

	rads := math.Atan(v.ij[0] / v.ij[1])

	// Convert the tan angle into the positive x-axis anti-clockwise.
	switch v.direction {
	case VECTOR_DIRECTION_Q1:
		fallthrough
	case VECTOR_DIRECTION_Q2:
		rads = math.Pi*0.5 - rads
	case VECTOR_DIRECTION_Q3:
		fallthrough
	case VECTOR_DIRECTION_Q4:
		rads = math.Pi*1.5 - rads
	}

	return Radians(rads)
}

// Units returns the units of the vector.
func (v Vector) Units() (Length, Length) {
	return Length(v.ij[0]), Length(v.ij[1])
}

// OrErr tests if either unit is NaN or Inf and returns an error if one is. NaN
// errors are prioritized over Inf errors.
func (v Vector) OrErr() (Vector, *FloatingPointError) {
	i, j := v.Units()
	_, ierr := i.OrErr()
	_, jerr := j.OrErr()
	if ierr != nil && ierr.IsNaN() {
		return v, ierr
	} else if jerr != nil && jerr.IsNaN() {
		return v, jerr
	} else if ierr != nil {
		return v, ierr
	} else if jerr != nil {
		return v, jerr
	}
	return v, nil
}

// String outputs the units.
func (v Vector) String() string {
	return fmt.Sprintf("Vector(Point({%s, %s}))",
		HumanFormat(9, v.ij[0]),
		HumanFormat(9, v.ij[1]))
}

// Rotate creates a new vector that has been rotated \c theta radians
// anti-clockwise.
func (v Vector) Rotate(rads Radians) Vector {
	a := mgl64.Mat2{
		math.Cos(float64(rads)), math.Sin(float64(rads)),
		-math.Sin(float64(rads)), math.Cos(float64(rads)),
	}
	ij := a.Mul2x1(v.ij)
	return VectorFromVec2(ij)
}

// Scale does scalar multiplication of the Vector
func (v Vector) Scale(m Length) Vector {
	return v.ScaleUnits(m, m)
}

// ScaleUnits scales the units of the vector independently.
func (v Vector) ScaleUnits(mx, my Length) Vector {
	a := mgl64.Mat2{
		float64(mx), 0,
		0, float64(my),
	}
	ij := a.Mul2x1(v.ij)
	return VectorFromVec2(ij)
}

// SkewUnits skews the vector by the provided units.
func (v Vector) SkewUnits(sx, sy Length) Vector {
	a := mgl64.Mat2{
		1, float64(sy),
		float64(sx), 1,
	}
	ij := a.Mul2x1(v.ij)
	return VectorFromVec2(ij)
}

// Invert the values of this vector. Return (-i, -j)
func (v Vector) Invert() Vector {
	a := mgl64.Mat2{
		-1, 0,
		0, -1,
	}
	ij := a.Mul2x1(v.ij)
	return VectorFromVec2(ij)
}

// Normalize the vector to be a unit length
func (v Vector) Normalize() Vector {
	m := v.Magnitude()
	if v.direction == VECTOR_DIRECTION_NONE || IsZero(m) {
		return VectorNaN
	}
	return v.Scale(1 / m)
}

// Add the units of the vectors. Returns (v.i+n.i, v.j+n.j)
func (v Vector) Add(n Vector) Vector {
	ij := mgl64.Vec2{v.ij[0] + n.ij[0], v.ij[1] + n.ij[1]}
	return VectorFromVec2(ij)
}

// Dot product of the vector. Returns (v.i*n.i + v.j*n.j)
func (v Vector) Dot(n Vector) Length {
	return Length(v.ij[0]*n.ij[0] + v.ij[1]*n.ij[1])
}

// RotatePts rotates \c pts by \c theta around \c origin.
func RotatePts(theta Radians, origin Pt, pts []Pt) []Pt {
	ret := make([]Pt, len(pts))
	for h, p := range pts {
		v := origin.VectorTo(p)
		v = v.Rotate(theta)
		ret[h] = origin.Add(v)
	}
	return ret
}

// TranslatePts translates \c pts by \c v.
func TranslatePts(v Vector, pts []Pt) []Pt {
	tm := mgl64.Mat3{
		1, 0, 0,
		0, 1, 0,
		v.ij[0], v.ij[1], 1,
	}
	ret := make([]Pt, len(pts))
	for h, p := range pts {
		xyz := tm.Mul3x1(p.xy.Vec3(1))
		ret[h] = PtFromVec2(xyz.Vec2())
	}
	return ret
}

// ShearPts performs a shear rotation on \c pts by \c v.
func ShearPts(v Vector, pts []Pt) []Pt {
	tm := mgl64.Mat2{
		1, v.ij[1],
		v.ij[0], 1,
	}
	ret := make([]Pt, len(pts))
	for h, p := range pts {
		xy := tm.Mul2x1(p.xy)
		ret[h] = PtFromVec2(xy)
	}
	return ret
}

// ScalePts scales the coordinates of \c pts by \c v.
func ScalePts(v Vector, pts []Pt) []Pt {
	tm := mgl64.Diag2(v.ij)
	ret := make([]Pt, len(pts))
	for h, p := range pts {
		xy := tm.Mul2x1(p.xy)
		ret[h] = PtFromVec2(xy)
	}
	return ret
}

// Limits returns the min-x, max-x, min-y, max-y in that order.
func LimitsPts(pts []Pt) (Length, Length, Length, Length) {
	xs := make([]Length, len(pts))
	ys := make([]Length, len(pts))
	for h, p := range pts {
		xs[h], ys[h] = p.X(), p.Y()
	}
	return Minimum(xs...), Maximum(xs...), Minimum(ys...), Maximum(ys...)
}

// SortPts sorts \c pts.
func SortPts(pts []Pt) []Pt {
	ptslice := ptSlice(pts)
	sort.Sort(ptslice)
	return []Pt(ptslice)
}

// IsEqualPair takes two objects that implement the pair interface and compares
// that they are equal.
func IsEqualPair[T Pair](a, b T) bool {
	ax, ay := a.Units()
	bx, by := b.Units()
	if IsEqual(ax, bx) && IsEqual(ay, by) {
		return true
	}
	return false
}

// IsZeropair checks if both units of a Pair are really close to zero.
func IsZeroPair[T Pair](a T) bool {
	ax, ay := a.Units()
	if IsZero(ax) && IsZero(ay) {
		return true
	}
	return false
}
