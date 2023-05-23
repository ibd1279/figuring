package figuring

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

var (
	MatrixBezierCubic mgl64.Mat4 = mgl64.Mat4{
		1, 0, 0, 0,
		-3, 3, 0, 0,
		3, -6, 3, 0,
		-1, 3, -3, 1,
	}
)

// deCasteljau Performs de Casteljau to return all the subpoints for the
// derivitive of points represented by the bezierCurve.
//
// https://pomax.github.io/bezierinfo/
func deCasteljau(s []Pt, tf float64) ([]Pt, []Pt) {
	t := Length(tf)
	pts := make([]Pt, len(s))
	copy(pts, s)

	left, right := []Pt{pts[0]}, []Pt{pts[len(pts)-1]}
	for len(pts) > 1 {
		newpts := make([]Pt, len(pts)-1)
		for h := 0; h < len(newpts); h++ {
			x := (1-t)*pts[h].X() + t*pts[h+1].X()
			y := (1-t)*pts[h].Y() + t*pts[h+1].Y()
			newpts[h] = PtXy(x, y)
		}
		left = append(left, newpts[0])
		right = append(right, newpts[len(newpts)-1])
		pts = newpts
	}

	return left, right
}

// ParamCurve is a curve defined by a pair of parametric functions. It doesn't
// provide a lot of functionality, but does provide an easy way to recreate
// curves based on polynomial equations.
//
// The ParamCurve does not keep track of the points that created the curve and
// only retains the polynomial equations required for the curve. The min and
// max constraints allow for some range checking on the curve.
type ParamCurve struct {
	X, Y     Derivable
	Min, Max float64
}

// ParamCubic creates a cubic bezier curve based on the four provided points.
// If more cubic bezier curve features are needed, use the Bezier type instead.
func ParamCubic(p1, p2, p3, p4 Pt) ParamCurve {
	px := mgl64.Vec4{float64(p4.X()), float64(p3.X()), float64(p2.X()), float64(p1.X())}
	py := mgl64.Vec4{float64(p4.Y()), float64(p3.Y()), float64(p2.Y()), float64(p1.Y())}
	xs, ys := MatrixBezierCubic.Mul4x1(px), MatrixBezierCubic.Mul4x1(py)
	return ParamCurve{
		X:   CubicFromVec4(xs),
		Y:   CubicFromVec4(ys),
		Min: 0,
		Max: 1.0,
	}
}

// ParamLinear creates a ParamCurve based on Linear equations between two
// points. Effectively Lerp or a linear bezier.
func ParamLinear(p1, p2 Pt) ParamCurve {
	ax, bx := -p1.X()+p2.X(), p1.X()
	ay, by := -p1.Y()+p2.Y(), p1.Y()
	return ParamCurve{
		X:   LinearAb(float64(ax), float64(bx)),
		Y:   LinearAb(float64(ay), float64(by)),
		Min: 0,
		Max: 1.0,
	}
}

// ParamQuadratic creates a quadratic bezier curve based on the three provided
// points.
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
		X:   QuadraticFromVec3(xs),
		Y:   QuadraticFromVec3(ys),
		Min: 0,
		Max: 1.0,
	}
}

// PtAtT returns the point for the provided value of \c t.
func (pc ParamCurve) PtAtT(t float64) Pt {
	t = Clamp(pc.Min, t, pc.Max)
	x, y := pc.X.AtT(t), pc.Y.AtT(t)
	return PtXy(Length(x), Length(y))
}

// Roots returns the roots for the current curve. This is a helper function
// that filters the root values between the \c Min and \c Max values before
// returning them.
func (pc ParamCurve) Roots() ([]float64, []float64) {
	xr := pc.X.Roots()
	xroots := make([]float64, 0, len(xr))
	for h := 0; h < len(xr); h++ {
		r := xr[h]
		if IsZero(r) {
			r = 0
		} else if IsZero(1.0 - r) {
			r = 1
		}
		if pc.Min <= r && r <= pc.Max {
			xroots = append(xroots, r)
		}
	}
	yr := pc.Y.Roots()
	yroots := make([]float64, 0, len(yr))
	for h := 0; h < len(yr); h++ {
		r := yr[h]
		if IsZero(r) {
			r = 0
		} else if IsZero(1.0 - r) {
			r = 1
		}
		if pc.Min <= r && r <= pc.Max {
			yroots = append(yroots, r)
		}

	}
	return xroots, yroots
}

// String returns a string representation of the ParamCurve. Format allows the
// curve to be pasted into Geogebra.
func (pc ParamCurve) String() string {
	unknown := 't'
	return fmt.Sprintf("Curve(%s, %s, %c, %s, %s)",
		pc.X.Text(unknown, false),
		pc.Y.Text(unknown, false),
		unknown,
		HumanFormat(9, pc.Min),
		HumanFormat(9, pc.Max),
	)
}

// TangentAtT returns the tangent and the normal of the curve for the given
// value of \c t.
func (pc ParamCurve) TangentAtT(t float64) (Vector, Vector) {
	t = Clamp(pc.Min, t, pc.Max)
	ieq, jeq := pc.X.Derivative(), pc.Y.Derivative()
	i, j := ieq.AtT(t), jeq.AtT(t)
	tangent := VectorIj(Length(i), Length(j)).Normalize()
	return tangent, tangent.Rotate(0.5 * math.Pi)
}

type BezierCurveType uint

const (
	BEZIER_CURVE_TYPE_PLAIN BezierCurveType = iota
	BEZIER_CURVE_TYPE_LOOP
	BEZIER_CURVE_TYPE_CUSP
	BEZIER_CURVE_TYPE_LOOPEND
	BEZIER_CURVE_TYPE_LOOPBEGIN
	BEZIER_CURVE_TYPE_SINGLEINFLECTION
	BEZIER_CURVE_TYPE_DOUBLEINFLECTION
)

// Represents a Cubic Bezier Curve.
type Bezier struct {
	pts  [4]Pt
	x, y Cubic
}

// BezierPt creates a new Bezier curve based on the provided points.
func BezierPt(p1, p2, p3, p4 Pt) Bezier {
	px := mgl64.Vec4{float64(p4.X()), float64(p3.X()), float64(p2.X()), float64(p1.X())}
	py := mgl64.Vec4{float64(p4.Y()), float64(p3.Y()), float64(p2.Y()), float64(p1.Y())}
	xs, ys := MatrixBezierCubic.Mul4x1(px), MatrixBezierCubic.Mul4x1(py)
	return Bezier{
		pts: [4]Pt{p1, p2, p3, p4},
		x:   CubicFromVec4(xs),
		y:   CubicFromVec4(ys),
	}
}

// AlignOnX rotates, translates, and scales the Bezier to the X-Axis, with the
// first point on the origin and the last point (1,0).  If the last point is at
// zero on the x-axis, it skips the scale operation.
func (curve Bezier) AlignOnX() (Vector, Radians, Length, Bezier) {
	translate := curve.pts[0].VectorTo(PtOrig)
	pts := TranslatePts(translate, curve.Points())
	theta := -PtOrig.VectorTo(pts[3]).Angle()
	pts = RotatePts(theta, PtOrig, pts)
	scale := pts[3].X()
	if !IsZero(scale) {
		pts = ScalePts(VectorIj(1/scale, 1/scale), pts)
	}

	return translate, theta, scale, BezierPt(pts[0], pts[1], pts[2], pts[3])
}

// ApproxLength treats the curve \c steps number of line segments and returns
// the sum of the length of all the line segments. It isn't as accurate as the
// \c Length() function, but can be much faster for smaller values of \c steps.
func (curve Bezier) ApproxLength(steps int) Length {
	prev := curve.PtAtT(0)
	var sum Length
	for h := 1; h <= steps; h++ {
		t := 1.0 / float64(steps) * float64(h)
		curr := curve.PtAtT(t)
		sum += prev.VectorTo(curr).Magnitude()
		prev = curr
	}
	return sum
}

// BoundingBox returns an axis-aligned rectangle that encompasses all the
// points of the curve.
func (curve Bezier) BoundingBox() Rectangle {
	ieq, jeq := curve.x.FirstDerivative(), curve.y.FirstDerivative()
	roots := []float64{0.0, 1.0}
	roots = append(roots, ieq.Roots()...)
	roots = append(roots, jeq.Roots()...)
	pts := make([]Pt, 0, len(roots))
	for h := 0; h < len(roots); h++ {
		if 0 <= roots[h] && roots[h] <= 1.0 {
			pts = append(pts, curve.PtAtT(roots[h]))
		}
	}
	lx, mx, ly, my := LimitsPts(pts)
	return RectanglePt(PtXy(lx, ly), PtXy(mx, my))
}

// CurveType returns the type of curve this is. See BezierCurveType for more
// details on return values.
func (curve Bezier) CurveType() BezierCurveType {
	// See https://pomax.github.io/bezierinfo/#canonical
	translate := curve.pts[0].VectorTo(PtOrig)
	pts := TranslatePts(translate, curve.Points())

	x2, y2 := pts[1].XY()
	x3, y3 := pts[2].XY()
	x4, y4 := pts[3].XY()

	y42 := y4 / y2
	y32 := y3 / y2

	x43 := (x4 - x2*y42) / (x3 - x2*y32)
	x := float64(x43)
	y := float64(y42 + x43*(1-y32))

	if y > 1 {
		return BEZIER_CURVE_TYPE_SINGLEINFLECTION
	}

	if y <= 1 && x <= 1 {
		c := (-x*x + 2*x + 3) / 4

		if x <= 0 {
			t0loop := (-x*x + 3*x) / 3
			if IsEqual(y, t0loop) {
				return BEZIER_CURVE_TYPE_LOOPBEGIN
			}
			if t0loop < y && y < c {
				return BEZIER_CURVE_TYPE_LOOP
			}
		}

		if 0 <= x && x <= 1.0 {
			t1loop := (math.Sqrt(3)*math.Sqrt(4*x-x*x) - x) / 2
			if IsEqual(y, t1loop) {
				return BEZIER_CURVE_TYPE_LOOPEND
			}
			if t1loop < y && y < c {
				return BEZIER_CURVE_TYPE_LOOP
			}
		}

		if IsEqual(y, c) {
			return BEZIER_CURVE_TYPE_CUSP
		}
		if y > c {
			return BEZIER_CURVE_TYPE_DOUBLEINFLECTION
		}
	}
	return BEZIER_CURVE_TYPE_PLAIN
}

// InflectionPts returns the points where the curvature of the curve switches
// directions.
func (curve Bezier) InflectionPts() []float64 {
	_, _, _, ac := curve.AlignOnX()
	// https://pomax.github.io/bezierinfo/#inflections
	a := ac.pts[2].X() * ac.pts[1].Y()
	b := ac.pts[3].X() * ac.pts[1].Y()
	c := ac.pts[1].X() * ac.pts[2].Y()
	d := ac.pts[3].X() * ac.pts[2].Y()

	x := -3*a + 2*b + 3*c - d
	y := 3*a - b - 3*c
	z := c - a

	eq := QuadraticAbc(float64(x), float64(y), float64(z))
	roots := eq.Roots()

	validRoots := make([]float64, 0, len(roots))
	for h := 0; h < len(roots); h++ {
		if 0 <= roots[h] && roots[h] <= 1.0 {
			validRoots = append(validRoots, roots[h])
		}
	}

	return validRoots
}

// Length returns a more accurate approximation than ApproxLength.
func (curve Bezier) Length() Length {
	// see https://pomax.github.io/bezierinfo/legendre-gauss.html
	z := 1.
	var sum float64
	for h := 0; h < len(legendregauss_weight); h++ {
		C := legendregauss_weight[h]
		T := legendregauss_abscissa[h]
		t := (z/2)*T + (z / 2)

		x := curve.x.FirstDerivative().AtT(t)
		y := curve.y.FirstDerivative().AtT(t)

		sum += C * math.Sqrt(x*x+y*y)
	}

	return Length(sum * (z / 2))
}

// Points provides access to the individual points of this curve. Consider the
// points readonly.
func (curve Bezier) Points() []Pt { return curve.pts[:] }

// PtAtT returns the point for the provided value of \c t.
func (curve Bezier) PtAtT(t float64) Pt {
	x, y := curve.x.AtT(t), curve.y.AtT(t)
	return PtXy(Length(x), Length(y))
}

// Roots returns the roots for the current curve. See Also Bezier.AlignOnX()
// and RotateOrTranslateToXAxis()
func (curve Bezier) Roots() ([]float64, []float64) {
	xr := curve.x.Roots()
	xroots := make([]float64, 0, len(xr))
	for h := 0; h < len(xr); h++ {
		r := xr[h]
		if IsZero(r) {
			r = 0
		} else if IsZero(1.0 - r) {
			r = 1
		}
		if 0 <= r && r <= 1.0 {
			xroots = append(xroots, r)
		}
	}
	yr := curve.y.Roots()
	yroots := make([]float64, 0, len(yr))
	for h := 0; h < len(yr); h++ {
		r := yr[h]
		if IsZero(r) {
			r = 0
		} else if IsZero(1.0 - r) {
			r = 1
		}
		if 0 <= r && r <= 1.0 {
			yroots = append(yroots, r)
		}
	}
	return xroots, yroots
}

// SplitAtT splits the current Bezier into 2 distinct bezier that have the are
// the same curvature.
func (curve Bezier) SplitAtT(t float64) (Bezier, Bezier) {
	px := mgl64.Vec4{
		float64(curve.pts[0].X()),
		float64(curve.pts[1].X()),
		float64(curve.pts[2].X()),
		float64(curve.pts[3].X()),
	}
	py := mgl64.Vec4{
		float64(curve.pts[0].Y()),
		float64(curve.pts[1].Y()),
		float64(curve.pts[2].Y()),
		float64(curve.pts[3].Y()),
	}

	z := t - 1
	qa := mgl64.Mat4{
		1, -z, z * z, -(z * z * z),
		0, t, -2 * z * t, 3 * (z * z) * t,
		0, 0, t * t, -3 * z * (t * t),
		0, 0, 0, t * t * t,
	}
	qb := mgl64.Mat4{
		-(z * z * z), 0, 0, 0,
		3 * (z * z) * t, z * z, 0, 0,
		-3 * z * (t * t), -2 * z * t, -z, 0,
		t * t * t, t * t, t, 1,
	}
	pax := qa.Mul4x1(px)
	pay := qa.Mul4x1(py)
	pbx := qb.Mul4x1(px)
	pby := qb.Mul4x1(py)

	return BezierPt(
			PtXy(Length(pax[0]), Length(pay[0])),
			PtXy(Length(pax[1]), Length(pay[1])),
			PtXy(Length(pax[2]), Length(pay[2])),
			PtXy(Length(pax[3]), Length(pay[3])),
		),
		BezierPt(
			PtXy(Length(pbx[0]), Length(pby[0])),
			PtXy(Length(pbx[1]), Length(pby[1])),
			PtXy(Length(pbx[2]), Length(pby[2])),
			PtXy(Length(pbx[3]), Length(pby[3])),
		)
}

// String returns a string representation of the bezier. Format allows the
// curve to be pasted into Geogebra.
func (curve Bezier) String() string {
	unknown := 't'
	return fmt.Sprintf("Bezier[ Curve(%s, %s, %c, 0, 1) ]",
		curve.x.Text(unknown, false),
		curve.y.Text(unknown, false),
		unknown,
	)
}

// TangentAtT returns the tangent and the normal of the curve for the given
// value of \c t.
func (curve Bezier) TangentAtT(t float64) (Vector, Vector) {
	ieq, jeq := curve.x.FirstDerivative(), curve.y.FirstDerivative()
	i, j := ieq.AtT(t), jeq.AtT(t)
	tangent := VectorIj(Length(i), Length(j))
	normal := VectorIj(-Length(j), Length(i))
	return tangent, normal
}
