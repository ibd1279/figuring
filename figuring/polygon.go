package figuring

import (
	"fmt"
	"math"
	"strings"
)

// OrderedPtser is an interface for all the types that provide a Points()
// interface.
type OrderedPtser interface {
	Points() []Pt
}

// BoundingBoxer is the interface set of all types that provide a BoundingBox()
// interface.
type BoundingBoxer interface {
	BoundingBox() Rectangle
}

// Rectangle represents an axis aligned rectangle. The resulting rectangle will
// always be aligned with the X and Y axis.
type Rectangle struct {
	pts [2]Pt
}

// Create a rectangle based on points.
func RectanglePt(p1, p2 Pt) Rectangle {
	var (
		err   *FloatingPointError
		broke bool
	)
	if _, err = p1.OrErr(); err != nil {
		broke = err.IsNaN()
	}
	if _, err = p2.OrErr(); err != nil && !broke {
		broke = err.IsNaN()
	}

	if broke {
		return Rectangle{
			pts: [2]Pt{PtNaN, PtNaN},
		}
	}

	lx, mx, ly, my := LimitsPts([]Pt{p1, p2})
	return Rectangle{
		pts: [2]Pt{PtXy(lx, ly), PtXy(mx, my)},
	}
}

func RectangleAppend(r1, r2 Rectangle) Rectangle {
	var (
		err            *FloatingPointError
		broke1, broke2 bool
	)
	if _, err = r1.OrErr(); err != nil {
		broke1 = err.IsNaN()
	}
	if _, err = r2.OrErr(); err != nil {
		broke2 = err.IsNaN()
	}
	switch {
	case broke1 && broke2:
		return Rectangle{
			pts: [2]Pt{PtNaN, PtNaN},
		}
	case broke1:
		return r2
	case broke2:
		return r1
	}

	var pts [4]Pt
	copy(pts[:2], r1.pts[:])
	copy(pts[2:], r2.pts[:])

	lx, mx, ly, my := LimitsPts(pts[:])

	return Rectangle{
		pts: [2]Pt{PtXy(lx, ly), PtXy(mx, my)},
	}
}

func (r Rectangle) Dims() (Length, Length) {
	return r.pts[0].VectorTo(r.pts[1]).Units()
}
func (r Rectangle) Height() Length {
	_, h := r.Dims()
	return h
}
func (r Rectangle) MaxPt() Pt { return r.pts[1] }
func (r Rectangle) MinPt() Pt { return r.pts[0] }
func (r Rectangle) OrErr() (Rectangle, *FloatingPointError) {
	least, most := r.pts[0], r.pts[1]
	_, lerr := least.OrErr()
	_, merr := most.OrErr()
	if lerr != nil && lerr.IsNaN() {
		return r, lerr
	} else if merr != nil && merr.IsNaN() {
		return r, merr
	} else if lerr != nil {
		return r, lerr
	} else if merr != nil {
		return r, merr
	}
	return r, nil
}
func (r Rectangle) Points() []Pt { return r.pts[:] }
func (r Rectangle) Sides() []Segment {
	minmax, maxmin := PtXy(r.pts[0].X(), r.pts[1].Y()), PtXy(r.pts[1].X(), r.pts[0].Y())
	return []Segment{
		SegmentPt(r.pts[0], minmax),
		SegmentPt(minmax, r.pts[1]),
		SegmentPt(r.pts[1], maxmin),
		SegmentPt(maxmin, r.pts[0]),
	}
}
func (r Rectangle) String() string {
	minmax, maxmin := PtXy(r.pts[0].X(), r.pts[1].Y()), PtXy(r.pts[1].X(), r.pts[0].Y())
	return fmt.Sprintf("Rectangle[ Polygon(%v, %v, %v, %v) ]",
		r.pts[0], minmax, r.pts[1], maxmin)
}
func (r Rectangle) Width() Length {
	w, _ := r.Dims()
	return w
}

// ClipToRectangleSegment Clips the provided provided segment, keeping only the
// parts of the segment inside the rectangle. Returns an empty slice if the
// segment doesn't exist inside the rectangle.
func ClipToRectangleSegment(a Rectangle, b Segment) []Segment {
	// based on Liang-Barsky: https://en.wikipedia.org/wiki/Liang%E2%80%93Barsky_algorithm
	// and https://www.skytopia.com/project/articles/compsci/clipping.html
	min, max := a.MinPt(), a.MaxPt()

	pnt := b.Begin()
	vec := pnt.VectorTo(b.End())

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

	return []Segment{
		SegmentPt(
			PtXy(pnt.X()+p2*rn1, pnt.Y()+p4*rn1),
			PtXy(pnt.X()+p2*rn2, pnt.Y()+p4*rn2),
		),
	}
}

// Unit objects, including triangles and rectangles.
var (
	Half          = Length(0.5)
	HalfSqrtTwo   = Length(math.Sqrt(2.) / 2.)
	HalfSqrtThree = Length(math.Sqrt(3.) / 2.)

	TriangleThirtySixtyNinety = PolygonPt(PtOrig, PtXy(HalfSqrtThree, 0), PtXy(HalfSqrtThree, Half))
	TriangleIsoscelesRight    = PolygonPt(PtOrig, PtXy(HalfSqrtTwo, 0), PtXy(HalfSqrtTwo, HalfSqrtTwo))
	TriangleSixtyNinetyThirty = PolygonPt(PtOrig, PtXy(Half, 0), PtXy(Half, HalfSqrtThree))
	TriangleEquilateral       = PolygonPt(PtOrig, PtXy(HalfSqrtThree, -Half), PtXy(HalfSqrtThree, Half))

	Square = PolygonPt(PtOrig, PtXy(1, 0), PtXy(1, 1), PtXy(0, 1))
)

type Polygon struct {
	pts []Pt
}

func PolygonPt(pts ...Pt) Polygon {
	return Polygon{
		pts: pts,
	}
}
func PolygonFromRectangle(r Rectangle) Polygon {
	min, max := r.MinPt(), r.MaxPt()
	return PolygonPt(
		min,
		PtXy(max.X(), min.Y()),
		max,
		PtXy(min.X(), max.Y()),
	)
}

func (poly Polygon) Angles() []Radians {
	angles := make([]Radians, 0, len(poly.pts))
	sides := poly.Sides()
	prev := sides[len(sides)-1]
	for h := 0; h < len(sides); h++ {
		curr := sides[h]
		a0 := (prev.Angle() + math.Pi).Normalize()
		a1 := curr.Angle()
		var a Radians
		if a1 > a0 {
			a = 2.0*math.Pi - (a1 - a0)
		} else {
			a = a0 - a1
		}
		angles = append(angles, a.Normalize())
		prev = curr
	}
	return angles
}
func (poly Polygon) Perimeter() Length {
	var sum Length
	for _, side := range poly.Sides() {
		sum += side.Length()
	}
	return sum
}
func (poly Polygon) Points() []Pt { return poly.pts[:] }
func (poly Polygon) OrErr() (Polygon, *FloatingPointError) {
	var err *FloatingPointError
	for _, p := range poly.pts {
		_, perr := p.OrErr()
		if perr != nil && perr.IsNaN() {
			return poly, perr
		} else if perr != nil {
			err = perr
		}
	}
	return poly, err
}
func (poly Polygon) Rotate(theta Radians, origin Pt) Polygon {
	return PolygonPt(RotatePts(theta, origin, poly.pts[:])...)
}
func (poly Polygon) Scale(scalars Vector) Polygon {
	return PolygonPt(ScalePts(scalars, poly.pts[:])...)
}
func (poly Polygon) Sides() []Segment {
	sides := make([]Segment, 0, len(poly.pts))
	prev := poly.pts[0]
	for h := 1; h < len(poly.pts); h++ {
		curr := poly.pts[h]
		sides = append(sides, SegmentPt(prev, curr))
		prev = curr
	}
	sides = append(sides, SegmentPt(prev, poly.pts[0]))
	return sides
}
func (poly Polygon) String() string {
	var slice []string
	for _, p := range poly.pts {
		slice = append(slice, fmt.Sprintf("%v", p))
	}
	return fmt.Sprintf("Polygon(%s)",
		strings.Join(slice, ", "))
}
func (poly Polygon) Translate(direction Vector) Polygon {
	return PolygonPt(TranslatePts(direction, poly.pts[:])...)
}
