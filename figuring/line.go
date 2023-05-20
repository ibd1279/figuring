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

// ParamCurve is a curve defined by a pair of parametric functions. It doesn't
// provide a lot of functionality, but does provide an easy way to recreate
// curves based on polynomial equations.
type ParamCurve struct {
	x, y     Derivable
	min, max float64
}

func ParamLinear(p1, p2 Pt) ParamCurve {
	ax, bx := -p1.X()+p2.X(), p1.X()
	ay, by := -p1.Y()+p2.Y(), p1.Y()
	return ParamCurve{
		x:   LinearAb(float64(ax), float64(bx)),
		y:   LinearAb(float64(ay), float64(by)),
		min: 0,
		max: 1.0,
	}
}
func ParamQuadratic(p1, p2, p3 Pt) ParamCurve {
	M := mgl64.Mat3{
		1, 0, 0,
		-2, 2, 0,
		1, -2, 1,
	}
	px := mgl64.Vec3{float64(p3.X()), float64(p2.X()), float64(p1.X())}
	py := mgl64.Vec3{float64(p3.Y()), float64(p2.Y()), float64(p1.Y())}
	xs, ys := M.Mul3x1(px), M.Mul3x1(py)
	return ParamCurve{
		x:   QuadraticFromVec3(xs),
		y:   QuadraticFromVec3(ys),
		min: 0,
		max: 1.0,
	}
}
func ParamCubic(p1, p2, p3, p4 Pt) ParamCurve {
	M := mgl64.Mat4{
		1, 0, 0, 0,
		-3, 3, 0, 0,
		3, -6, 3, 0,
		-1, 3, -3, 1,
	}
	px := mgl64.Vec4{float64(p4.X()), float64(p3.X()), float64(p2.X()), float64(p1.X())}
	py := mgl64.Vec4{float64(p4.Y()), float64(p3.Y()), float64(p2.Y()), float64(p1.Y())}
	xs, ys := M.Mul4x1(px), M.Mul4x1(py)
	return ParamCurve{
		x:   CubicFromVec4(xs),
		y:   CubicFromVec4(ys),
		min: 0,
		max: 1.0,
	}
}

func (pc ParamCurve) String() string {
	unknown := 't'
	return fmt.Sprintf("Curve(%s, %s, %c, 0, 1)",
		pc.x.Text(unknown, false),
		pc.y.Text(unknown, false),
		unknown,
	)
}
func (pc ParamCurve) PtAtT(t float64) Pt {
	x, y := pc.x.AtT(t), pc.y.AtT(t)
	return PtXy(Length(x), Length(y))
}
func (pc ParamCurve) TangentAtT(t float64) (Vector, Vector) {
	ieq, jeq := pc.x.Derivative(), pc.y.Derivative()
	i, j := ieq.AtT(t), jeq.AtT(t)
	tangent := VectorIj(Length(i), Length(j)).Normalize()
	return tangent, tangent.Rotate(0.5 * math.Pi)
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
