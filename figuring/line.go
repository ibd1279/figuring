package figuring

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Coefficienter interface {
	Coefficients() []float64
}

var (
	LineXAxis Line = LineAbc(0, 1, 0)
	LineYAxis Line = LineAbc(1, 0, 0)
)

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
	return LineAbc(a, -b, c)
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
// B. Effectively reducing the formula to ax+b=c, while maintaining the same
// line. Will cause the line to be in error if \c IsVertical() is true.
func (le Line) NormalizeY() Line {
	return LineFromVec3(mgl64.Vec3{le.abc[0] / le.abc[1], 1, le.abc[2] / le.abc[1]})
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
			HumanFormat(9, le.abc[2]))
	case LINE_DIRECTION_HORIZONTAL:
		str = fmt.Sprintf("%sy=%s",
			HumanFormat(9, le.abc[1]),
			HumanFormat(9, le.abc[2]))
	case LINE_DIRECTION_VERTICAL:
		str = fmt.Sprintf("%sx=%s",
			HumanFormat(9, le.abc[0]),
			HumanFormat(9, le.abc[2]))
	case LINE_DIRECTION_SAME:
		fallthrough
	default:
		sign := '+'
		b := le.abc[1]
		if Signbit(b) {
			sign = '-'
			b = -b
		}
		str = fmt.Sprintf("%sx%c%sy=%s",
			HumanFormat(9, le.abc[0]),
			sign,
			HumanFormat(9, b),
			HumanFormat(9, le.abc[2]),
		)
	}
	return str
}

// Vector returns the vector of the line, in the direction of A, normalized to
// a magnitude of 1.
func (le Line) Vector() Vector {
	ij := mgl64.Vec2{-le.abc[1], le.abc[0]}
	return VectorFromVec2(ij).Normalize()
}

// XForY returns the X value for a given Y. Returns \c NaN if \c IsHorizontal()
// or \c IsUnknown() are true.
func (le Line) XForY(y Length) Length {
	switch le.s {
	case LINE_DIRECTION_VERTICAL:
		return Length(le.abc[2] / le.abc[0])
	case LINE_DIRECTION_HORIZONTAL:
		fallthrough
	case LINE_DIRECTION_UNKNOWN:
		return Length(math.NaN())
	}

	a, b, c := le.Abc()
	return -b*y/a + c/a
}

// YForX returns the Y value for a given X. Returns \c NaN if \c IsVertical()
// or \c IsUnknown() are true.
func (le Line) YForX(x Length) Length {
	switch le.s {
	case LINE_DIRECTION_HORIZONTAL:
		return Length(le.abc[2] / le.abc[1])
	case LINE_DIRECTION_VERTICAL:
		fallthrough
	case LINE_DIRECTION_UNKNOWN:
		return Length(math.NaN())
	}

	a, b, c := le.Abc()
	return -a*x/b + c/b
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

// IntersectionLineLine returns the intersection points of two lines. returns
// an empty slice if the lines do not intersect.
func IntersectionLineLine(a, b Line) []Pt {
	aTheta, bTheta := a.Angle(), b.Angle()
	if IsEqual(aTheta, bTheta) {
		// Parallel lines cannot meet in this geometry.
		// also catches the same line passed twice
		return nil
	}

	var p Pt
	switch {
	case a.IsUnknown():
		fallthrough
	case b.IsUnknown():
		return nil
	case a.IsVertical():
		b, a = a, b
		fallthrough
	case b.IsVertical():
		x := b.XForY(0)
		y := a.YForX(x)
		p = PtXy(x, y)
	case a.IsHorizontal():
		b, a = a, b
		fallthrough
	case b.IsHorizontal():
		y := b.YForX(0)
		x := a.XForY(y)
		p = PtXy(x, y)
	default:
		na, nb := a.NormalizeY(), b.NormalizeY()
		ma, _, ba := na.Abc()
		mb, _, bb := nb.Abc()

		x := Length((bb - ba) / (mb - ma))
		y := b.YForX(x)

		p = PtXy(x, y)
	}

	return []Pt{p}
}

func IntersectionLineBezier(a Line, b Bezier) []Pt {
	bb := b.BoundingBox()
	grossIntersections := IntersectionRectangleLine(bb, a)
	if len(grossIntersections) == 0 {
		return nil
	}

	var pts []Pt = RotateOrTranslateToXAxis(a, b.Points())

	// At this point, the line is now the X axis. Find the roots of the curve.
	b2 := BezierPt(pts[0], pts[1], pts[2], pts[3])
	yr := b2.y.Roots()
	roots := make([]Pt, 0, len(yr))
	for h := 0; h < len(yr); h++ {
		if 0 <= yr[h] && yr[h] <= 1.0 {
			roots = append(roots, b.PtAtT(yr[h]))
		}
	}

	return roots
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

func IntersectionRayLine(a Ray, b Line) []Pt {
	aLine := a.Line()
	pts := FilterPtsRay(a, IntersectionLineLine(aLine, b))
	if len(pts) == 0 {
		return nil
	}

	return pts
}
func IntersectionRayRay(a Ray, b Ray) []Pt {
	aLine := a.Line()
	bLine := b.Line()
	pts := FilterPtsRay(a, FilterPtsRay(b, IntersectionLineLine(aLine, bLine)))
	if len(pts) == 0 {
		return nil
	}

	return pts
}
func IntersectionRaySegment(a Ray, b Segment) []Pt {
	aLine := a.Line()
	pts := FilterPtsRay(a, IntersectionSegmentLine(b, aLine))
	if len(pts) == 0 {
		return nil
	}

	return pts
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

func IntersectionSegmentSegment(a, b Segment) []Pt {
	a1 := a.End().Y() - a.Begin().Y()
	b1 := a.Begin().X() - a.End().X()
	c1 := a1*a.Begin().X() + b1*a.Begin().Y()

	a2 := b.End().Y() - b.Begin().Y()
	b2 := b.Begin().X() - b.End().X()
	c2 := a2*b.Begin().X() + b2*b.Begin().Y()

	det := a1*b2 - a2*b1
	if IsZero(det) {
		return nil
	}
	x := (b2*c1 - b1*c2) / det
	y := (a1*c2 - a2*c1) / det

	alx, amx, aly, amy := LimitsPts(a.Points())
	blx, bmx, bly, bmy := LimitsPts(b.Points())

	lx, mx := Maximum(alx, blx), Minimum(amx, bmx)
	ly, my := Maximum(aly, bly), Minimum(amy, bmy)

	if lx <= x && x <= mx && ly <= y && y <= my {
		return []Pt{PtXy(x, y)}
	}
	return nil
}
func IntersectionSegmentLine(a Segment, b Line) []Pt {
	aLine := LineFromPt(a.Begin(), a.End())
	potentialPoints := IntersectionLineLine(aLine, b)
	if len(potentialPoints) == 0 {
		return nil
	}

	lx, mx, ly, my := LimitsPts(a.Points())
	for _, p := range potentialPoints {
		x, y := p.XY()
		if lx <= x && x <= mx && ly <= y && y <= my {
			return []Pt{p}
		}
	}
	return nil
}

func IntersectionSegmentBezier(a Segment, b Bezier) []Pt {
	aLine := LineFromPt(a.Begin(), a.End())
	potentialPoints := IntersectionLineBezier(aLine, b)
	if len(potentialPoints) == 0 {
		return nil
	}

	lx, mx, ly, my := LimitsPts(a.Points())
	points := make([]Pt, 0, len(potentialPoints))
	for _, p := range potentialPoints {
		x, y := p.XY()
		if lx <= x && x <= mx && ly <= y && y <= my {
			points = append(points, p)
		}
	}
	return points
}

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
