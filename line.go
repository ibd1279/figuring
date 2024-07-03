package figuring

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

// Coefficienter is an interface for all the equation types that depend on coefficients.
type Coefficienter interface {
	Coefficients() []float64
}

var (
	// LineXAxis is the line that represents the X axis.
	LineXAxis Line = LineAbc(0, 1, 0)

	// LineYAxis is the line that represents the Y axis.
	LineYAxis Line = LineAbc(1, 0, 0)
)

// SlopeType is the type (direction) of the slope of a line.
type SlopeType uint

const (
	LINE_DIRECTION_UNKNOWN SlopeType = iota
	LINE_DIRECTION_HORIZONTAL
	LINE_DIRECTION_VERTICAL
	LINE_DIRECTION_SAME
	LINE_DIRECTION_MIXED
)

// Line Equation that represents a line. ax+by+c=0.
// See https://scholarsarchive.byu.edu/cgi/viewcontent.cgi?article=1000&context=facpub
type Line struct {
	abc mgl64.Vec3
	s   SlopeType
}

// LineAbc creates a Line for the coefficients of a line. Expected format is
// ax+by=c. Values of a, b, or c that are close to zero are treated as zero.
func LineAbc(a, b, c Length) Line {
	return LineFromVec3(mgl64.Vec3{float64(a), float64(b), float64(c)})
}

// LineFromVec3 creates a Line based on the provided Vec3. Expected format is
// abc[0]x+abc[1]y=abc[2]. Values in abc that are close to zero are treated as
// zero.
func LineFromVec3(abc mgl64.Vec3) Line {
	if IsZero(abc[2]) {
		abc[2] = 0
	}

	var s SlopeType
	switch {
	case IsZero(abc[0]) && IsZero(abc[1]):
		abc[0], abc[1] = 0, 0
		s = LINE_DIRECTION_UNKNOWN
	case IsZero(abc[0]) && !IsZero(abc[1]):
		abc[0] = 0
		s = LINE_DIRECTION_HORIZONTAL
	case !IsZero(abc[0]) && IsZero(abc[1]):
		abc[1] = 0
		s = LINE_DIRECTION_VERTICAL
	case Signbit(abc[0]) != Signbit(abc[1]):
		s = LINE_DIRECTION_SAME
	default:
		s = LINE_DIRECTION_MIXED
	}
	return Line{
		abc: abc,
		s:   s,
	}
}

// LineFromVector creates a Line from a point and a vector. The vector is used
// to compute the slope (a and b values), and the point is used to compute the
// intercept (c value). Values of a, b, or c that are close to zero are treated
// as zero.
func LineFromVector(p1 Pt, v Vector) Line {
	b, a := v.Units()
	c := p1.X()*a - p1.Y()*b
	return LineAbc(a, -b, -c)
}

// LineFromPt create a line from two points. A line is a linear equation in the
// implicit format (ax+by=c). See \c Segment if you want to create a line that
// only exists between two points.
func LineFromPt(p1, p2 Pt) Line { return LineFromVector(p1, p1.VectorTo(p2)) }

// Abc returns the coefficients of the linear equation.
func (le Line) Abc() (Length, Length, Length) {
	return Length(le.abc[0]), Length(le.abc[1]), Length(le.abc[2])
}

// Angle returns the angle of the line, with positive X axis as being zero radians.
func (le Line) Angle() Radians {
	return le.Vector().Angle()
}

// Coefficients implements the Coefficienter interface. Provides generic access
// to the coefficients of the linear function.
func (le Line) Coefficients() []float64 { return le.abc[:] }

// IsHorizontal returns true if the line is a horizontal line (no rise).
func (le Line) IsHorizontal() bool { return le.s == LINE_DIRECTION_HORIZONTAL }

// IsVertical returns true if the line is a vertical line (no run).
func (le Line) IsVertical() bool { return le.s == LINE_DIRECTION_VERTICAL }

// IsUnknown returns true if the linear function has no slope (no rise, no
// run).
func (le Line) IsUnknown() bool { return le.s == LINE_DIRECTION_UNKNOWN }

// NormalizeX adjusts the coefficients of the linear function to have a 1 for
// the A. Effectively reducing the formula to x+by=c, while maintaining the
// same line. Will cause the line to be in error if \c IsHorizontal() is true.
func (le Line) NormalizeX() Line {
	return LineFromVec3(mgl64.Vec3{1, le.abc[1] / le.abc[0], le.abc[2] / le.abc[0]})
}

// NormalizeY adjusts the coefficients of the linear function to have a 1 for
// B. Effectively reducing the formula to ax+y=c, while maintaining the same
// line. Will cause the line to be in error if \c IsVertical() is true.
func (le Line) NormalizeY() Line {
	return LineFromVec3(mgl64.Vec3{le.abc[0] / le.abc[1], 1, le.abc[2] / le.abc[1]})
}

// NormalizeUnit adjusts the coefficients of the linear function to have a unit
// length of 1. Will cause the line to be in error if \c IsUnknown is true.
func (le Line) NormalizeUnit() Line {
	d := math.Hypot(le.abc[0], le.abc[1])
	return LineFromVec3(mgl64.Vec3{le.abc[0] / d, le.abc[1] / d, le.abc[2] / d})
}

// OrErr checks all the coefficients of the linear function and returns a
// floating point error if any of them are non-real floating point values (NaN,
// Inf, LINE_DIRECTION_UNKNOWN).
func (le Line) OrErr() (Line, *FloatingPointError) {
	if le.s == LINE_DIRECTION_UNKNOWN {
		return le, &FloatingPointError{math.Inf(1)}
	}
	a, b, c := le.Abc()
	if _, err := a.OrErr(); err != nil {
		return le, err
	} else if _, err = b.OrErr(); err != nil {
		return le, err
	} else if _, err = c.OrErr(); err != nil {
		return le, err
	}
	return le, nil
}

// String returns a human readable representation of the linear function.
func (le Line) String() string {
	var str string
	switch le.s {
	case LINE_DIRECTION_UNKNOWN:
		str = fmt.Sprintf("0x+0y=%s",
			HumanFormat(9, -le.abc[2]))
	case LINE_DIRECTION_HORIZONTAL:
		str = fmt.Sprintf("%sy=%s",
			HumanFormat(9, le.abc[1]),
			HumanFormat(9, -le.abc[2]))
	case LINE_DIRECTION_VERTICAL:
		str = fmt.Sprintf("%sx=%s",
			HumanFormat(9, le.abc[0]),
			HumanFormat(9, -le.abc[2]))
	case LINE_DIRECTION_SAME:
		fallthrough
	default:
		ab := '+'
		b := le.abc[1]
		if Signbit(b) {
			ab = '-'
			b = -b
		}
		bc := '+'
		c := le.abc[2]
		if Signbit(c) {
			bc = '-'
			c = -c
		}
		str = fmt.Sprintf("%sx%c%sy%c%s=0",
			HumanFormat(9, le.abc[0]),
			ab,
			HumanFormat(9, b),
			bc,
			HumanFormat(9, c),
		)
	}
	return str
}

// Vector returns the vector of the line, in the direction of A, normalized to
// a magnitude of 1.
func (le Line) Vector() Vector {
	le = le.NormalizeUnit()
	ij := mgl64.Vec2{-le.abc[1], le.abc[0]}
	return VectorFromVec2(ij)
}

// XForY returns the X value for a given Y. Returns \c NaN if \c IsHorizontal()
// or \c IsUnknown() are true.
func (le Line) XForY(y Length) Length {
	switch le.s {
	case LINE_DIRECTION_VERTICAL:
		return Length(-le.abc[2] / le.abc[0])
	case LINE_DIRECTION_HORIZONTAL:
		fallthrough
	case LINE_DIRECTION_UNKNOWN:
		return Length(math.NaN())
	}

	a, b, c := le.Abc()
	return b*y/-a - c/a
}

// YForX returns the Y value for a given X. Returns \c NaN if \c IsVertical()
// or \c IsUnknown() are true.
func (le Line) YForX(x Length) Length {
	switch le.s {
	case LINE_DIRECTION_HORIZONTAL:
		return Length(-le.abc[2] / le.abc[1])
	case LINE_DIRECTION_VERTICAL:
		fallthrough
	case LINE_DIRECTION_UNKNOWN:
		return Length(math.NaN())
	}

	a, b, c := le.Abc()
	return -a*x/b - c/b
}

// RotateOrTranslateToXAxis rotates (or translates) \c pts towards the X-axis,
// such that the provided line is the X-axis. Useful for intersection finding
// between lines and equations.
func RotateOrTranslateToXAxis(a Line, pts []Pt) []Pt {
	switch {
	case a.IsUnknown():
		return pts
	case a.IsHorizontal():
		y := a.YForX(0)
		if !IsZero(y) {
			trans := PtXy(0, y).VectorTo(PtOrig)
			pts = TranslatePts(trans, pts)
		}
	default:
		x := a.XForY(0)
		origin := PtXy(x, 0)
		theta := -a.Angle()
		pts = RotatePts(theta, origin, pts)
	}
	return pts
}

// Ray represents a geometric ray, with a specific starting point and a directional vector.
type Ray struct {
	b Pt
	v Vector
}

func RayFromVector(begin Pt, direction Vector) Ray {
	return Ray{
		b: begin,
		v: direction.Normalize(),
	}
}

func (r Ray) Angle() Radians { return r.v.Angle() }
func (r Ray) Begin() Pt      { return r.b }
func (r Ray) OrErr() (Ray, *FloatingPointError) {
	if _, err := r.b.OrErr(); err != nil {
		return r, err
	} else if _, err = r.v.OrErr(); err != nil {
		return r, err
	}
	return r, nil
}
func (r Ray) Invert() Ray    { return RayFromVector(r.b, r.v.Invert()) }
func (r Ray) Line() Line     { return LineFromVector(r.b, r.v) }
func (r Ray) String() string { return fmt.Sprintf("Ray(%v, %v)", r.b, r.v) }
func (r Ray) Vector() Vector { return r.v }

func FilterPtsRay(r Ray, pts []Pt) (ret []Pt) {
	for _, pp := range pts {
		if IsEqualPair(pp, r.Begin()) {
			ret = append(ret, pp)
		} else {
			v := r.Begin().VectorTo(pp).Normalize()
			if IsEqualPair(v, r.Vector()) {
				ret = append(ret, pp)
			}
		}
	}
	return ret
}

// Segment represents a line with a fixed slope between two points.
type Segment struct {
	b, e Pt
}

// SegmentPt creates a new segment using the provided points.
func SegmentPt(begin, end Pt) Segment {
	return Segment{
		b: begin,
		e: end,
	}
}

// SegmentFromVector creates a new segment using the provided orgiin and a
// vector to compute the end.
func SegmentFromVector(begin Pt, end Vector) Segment {
	return SegmentPt(begin, begin.Add(end))
}

func (s Segment) Begin() Pt              { return s.b }
func (s Segment) BoundingBox() Rectangle { return RectanglePt(s.b, s.e) }
func (s Segment) End() Pt                { return s.e }
func (s Segment) Length() Length         { return s.b.VectorTo(s.e).Magnitude() }
func (s Segment) Angle() Radians         { return s.b.VectorTo(s.e).Angle() }
func (s Segment) Points() []Pt           { return []Pt{s.b, s.e} }
func (s Segment) OrErr() (Segment, *FloatingPointError) {
	if _, err := s.b.OrErr(); err != nil {
		return s, err
	} else if _, err = s.e.OrErr(); err != nil {
		return s, err
	}
	return s, nil
}
func (s Segment) String() string {
	return fmt.Sprintf("Segment(%v, %v)", s.b, s.e)
}
func (s Segment) Reverse() Segment { return SegmentPt(s.e, s.b) }

func IsEqualPts[T OrderedPtser](a, b T) bool {
	as, bs := a.Points(), b.Points()
	if len(as) != len(bs) {
		return false
	}
	for h := 0; h < len(as); h++ {
		if !IsEqualPair(as[h], bs[h]) {
			return false
		}
	}
	return true
}
