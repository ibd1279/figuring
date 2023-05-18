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
	LineXAxis Linear = LinearAbc(0, 1, 0)
	LineYAxis Linear = LinearAbc(1, 0, 0)
)

type SlopeType uint

const (
	LINE_DIRECTION_UNKNOWN SlopeType = iota
	LINE_DIRECTION_HORIZONTAL
	LINE_DIRECTION_VERTICAL
	LINE_DIRECTION_SAME
	LINE_DIRECTION_MIXED
)

// Linear Equation that represents a line. ax+by+c=0.
// See https://scholarsarchive.byu.edu/cgi/viewcontent.cgi?article=1000&context=facpub
type Linear struct {
	abc mgl64.Vec3
	s   SlopeType
}

func LinearAbc(a, b, c Length) Linear {
	return LinearFromVec3(mgl64.Vec3{float64(a), float64(b), float64(c)})
}

func LinearFromVec3(abc mgl64.Vec3) Linear {
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
	return Linear{
		abc: abc,
		s:   s,
	}
}

func LinearFromVector(p1 Pt, v Vector) Linear {
	b, a := v.Units()
	c := p1.X()*a - p1.Y()*b
	return LinearAbc(a, -b, c)
}

func LinearFromPt(p1, p2 Pt) Linear {
	return LinearFromVector(p1, p1.VectorTo(p2))
}

func (le Linear) Abc() (Length, Length, Length) {
	return Length(le.abc[0]), Length(le.abc[1]), Length(le.abc[2])
}

func (le Linear) Coefficients() []float64 {
	return le.abc[:]
}

func (le Linear) OrErr() (Linear, *FloatingPointError) {
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

func (le Linear) String() string {
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

func (le Linear) IsHorizontal() bool { return le.s == LINE_DIRECTION_HORIZONTAL }
func (le Linear) IsVertical() bool   { return le.s == LINE_DIRECTION_VERTICAL }
func (le Linear) NormalizeX() Linear {
	return LinearFromVec3(mgl64.Vec3{1, le.abc[1] / le.abc[0], le.abc[2] / le.abc[0]})
}
func (le Linear) NormalizeY() Linear {
	return LinearFromVec3(mgl64.Vec3{le.abc[0] / le.abc[1], 1, le.abc[2] / le.abc[1]})
}
func (le Linear) XForY(y Length) Length {
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
func (le Linear) YForX(x Length) Length {
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
func (le Linear) Vector() Vector {
	ij := mgl64.Vec2{-le.abc[1], le.abc[0]}
	return VectorFromVec2(ij).Normalize()
}
func (le Linear) Angle() Radians {
	return le.Vector().Angle()
}

type OrderedPtser interface {
	Points() []Pt
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

// Reverse swaps the begin and end points of the segment.
func (s Segment) Reverse() Segment { return SegmentPt(s.e, s.b) }

func IsEqualEquations[T Coefficienter](a, b T) bool {
	as, bs := a.Coefficients(), b.Coefficients()
	if len(as) != len(bs) {
		return false
	}
	for h := 0; h < len(as); h++ {
		if !IsEqual(as[h], bs[h]) {
			return false
		}
	}
	return true
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
