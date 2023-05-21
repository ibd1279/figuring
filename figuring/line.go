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

func LineAbc(a, b, c Length) Line {
	return LineFromVec3(mgl64.Vec3{float64(a), float64(b), float64(c)})
}
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
func LineFromVector(p1 Pt, v Vector) Line {
	b, a := v.Units()
	c := p1.X()*a - p1.Y()*b
	return LineAbc(a, -b, c)
}
func LineFromPt(p1, p2 Pt) Line { return LineFromVector(p1, p1.VectorTo(p2)) }

func (le Line) Abc() (Length, Length, Length) {
	return Length(le.abc[0]), Length(le.abc[1]), Length(le.abc[2])
}
func (le Line) Coefficients() []float64 { return le.abc[:] }
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
func (le Line) IsHorizontal() bool { return le.s == LINE_DIRECTION_HORIZONTAL }
func (le Line) IsVertical() bool   { return le.s == LINE_DIRECTION_VERTICAL }
func (le Line) IsUnknown() bool    { return le.s == LINE_DIRECTION_UNKNOWN }
func (le Line) NormalizeX() Line {
	return LineFromVec3(mgl64.Vec3{1, le.abc[1] / le.abc[0], le.abc[2] / le.abc[0]})
}
func (le Line) NormalizeY() Line {
	return LineFromVec3(mgl64.Vec3{le.abc[0] / le.abc[1], 1, le.abc[2] / le.abc[1]})
}
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
func (le Line) Vector() Vector {
	ij := mgl64.Vec2{-le.abc[1], le.abc[0]}
	return VectorFromVec2(ij).Normalize()
}
func (le Line) Angle() Radians {
	return le.Vector().Angle()
}
func (le Line) Roots() []Pt {
	x := le.XForY(0)
	return []Pt{PtXy(x, 0)}
}

func RotateOrTranslateToXAxis(a Line, pts []Pt) []Pt {
	if a.IsHorizontal() {
		y := a.YForX(0)
		if !IsZero(y) {
			trans := PtXy(0, y).VectorTo(PtOrig)
			pts = TranslatePts(trans, pts)
		}
	} else {
		x := a.XForY(0)
		origin := PtXy(x, 0)
		theta := -a.Angle()
		pts = RotatePts(theta, origin, pts)
	}
	return pts
}

func IntersectionLineLine(a, b Line) []Pt {
	aTheta, bTheta := a.Angle(), b.Angle()
	if IsEqual(aTheta, bTheta) {
		// Parallel lines cannot meet in this geometry.
		// also catches the same line passed twice
		return nil
	}

	var p Pt
	switch {
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
		ma, mb = -ma, -mb

		x := Length((bb - ba) / (mb - ma))
		y := b.YForX(x)
		p = PtXy(x, y)
	}

	return []Pt{p}
}
func IntersectionLineSegment(a Line, b Segment) []Pt {
	bLine := LineFromPt(b.Begin(), b.End())
	potentialPoints := IntersectionLineLine(a, bLine)
	if len(potentialPoints) == 0 {
		return nil
	}

	lx, mx, ly, my := LimitsPts(b.Points())
	for _, p := range potentialPoints {
		x, y := p.XY()
		if lx <= x && x <= mx && ly <= y && y <= my {
			return []Pt{p}
		}
	}
	return nil
}
func IntersectionLineRectangle(a Line, b Rectangle) []Pt {
	min, max := b.MinPt(), b.MaxPt()

	var s Segment
	switch {
	case a.IsVertical():
		x := a.XForY(0)
		s = SegmentPt(PtXy(x, min.Y()), PtXy(x, max.Y()))
	case a.IsHorizontal():
		y := a.YForX(0)
		s = SegmentPt(PtXy(min.X(), y), PtXy(max.X(), y))
	default:
		miny := a.YForX(min.X())
		maxy := a.YForX(max.X())
		s = SegmentPt(PtXy(min.X(), miny), PtXy(max.X(), maxy))
	}
	return IntersectionSegmentRectangle(s, b)
}
func IntersectionLineBezier(a Line, b Bezier) []Pt {
	bb := b.BoundingBox()
	grossIntersections := IntersectionLineRectangle(a, bb)
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

// Segment represents a line with a fixed slope between two points.
type Segment struct {
	b, e Pt
}

// SegmentPt creates a new point using the provided points.
func SegmentPt(begin, end Pt) Segment {
	return Segment{
		b: begin,
		e: end,
	}
}

func (s Segment) Begin() Pt      { return s.b }
func (s Segment) End() Pt        { return s.e }
func (s Segment) Length() Length { return s.b.VectorTo(s.e).Magnitude() }
func (s Segment) Angle() Radians { return s.b.VectorTo(s.e).Angle() }
func (s Segment) Points() []Pt   { return []Pt{s.b, s.e} }
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

// Returns the start and stop points of the segment that exists inside the rectangle.
func IntersectionSegmentRectangle(a Segment, b Rectangle) []Pt {
	// based on Liang-Barsky: https://en.wikipedia.org/wiki/Liang%E2%80%93Barsky_algorithm
	// and https://www.skytopia.com/project/articles/compsci/clipping.html
	min, max := b.MinPt(), b.MaxPt()

	// Clip out a segment that is in the right X coordinate space. Min X value to max X value.

	pnt := a.Begin()
	vec := pnt.VectorTo(a.End())

	p1, p3 := vec.Invert().Units()
	p2, p4 := vec.Units()

	q1, q3 := pnt.X()-min.X(), pnt.Y()-min.Y()
	q2, q4 := max.X()-pnt.X(), max.Y()-pnt.Y()

	posarr, negarr := []Length{1}, []Length{0}

	r1, r2 := q1/p1, q2/p2
	if p1 < 0 {
		negarr = append(negarr, r1)
		posarr = append(posarr, r2)
	} else {
		negarr = append(negarr, r2)
		posarr = append(posarr, r1)
	}

	r3, r4 := q3/p3, q4/p4
	if p3 < 0 {
		negarr = append(negarr, r3)
		posarr = append(posarr, r4)
	} else {
		negarr = append(negarr, r4)
		posarr = append(posarr, r3)
	}

	rn1, rn2 := Maximum(negarr...), Minimum(posarr...)
	if rn1 > rn2 {
		return nil
	}

	return []Pt{
		PtXy(pnt.X()+p2*rn1, pnt.Y()+p4*rn1),
		PtXy(pnt.X()+p2*rn2, pnt.Y()+p4*rn2),
	}
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
