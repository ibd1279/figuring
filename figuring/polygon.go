package figuring

import (
	"fmt"
)

type OrderedPtser interface {
	Points() []Pt
}

type BoundingBoxer interface {
	BoundingBox() Rectangle
}

// Rectangle represents an axis aligned rectangle. The resulting rectangle will
// always be aligned with the X and Y axis.
type Rectangle struct {
	pts [2]Pt
}

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

func IntersectionRectangleLine(a Rectangle, b Line) []Pt {
	min, max := a.MinPt(), a.MaxPt()

	var s Segment
	switch {
	case b.IsVertical():
		x := b.XForY(0)
		s = SegmentPt(PtXy(x, min.Y()), PtXy(x, max.Y()))
	case b.IsHorizontal():
		y := b.YForX(0)
		s = SegmentPt(PtXy(min.X(), y), PtXy(max.X(), y))
	default:
		ly, lerr := b.YForX(min.X()).OrErr()
		my, merr := b.YForX(max.X()).OrErr()
		if lerr == nil && merr == nil {
			s = SegmentPt(PtXy(min.X(), ly), PtXy(max.X(), my))
		} else {
			// Don't check for errors here since there is no fall
			// back. let the Segment carry the error.
			lx := b.XForY(min.Y())
			mx := b.XForY(max.Y())
			s = SegmentPt(PtXy(lx, min.Y()), PtXy(mx, max.Y()))
		}
	}
	clipped := ClipToRectangleSegment(a, s)
	if len(clipped) == 0 {
		return nil
	}
	pts := make([]Pt, 0, len(clipped)*2)
	for h := 0; h < len(clipped); h++ {
		pts = append(pts, clipped[h].Points()...)
	}
	return pts
}

func IntersectionRectangleSegment(a Rectangle, b Segment) []Pt {
	min, max := a.MinPt(), a.MaxPt()

	clipped := ClipToRectangleSegment(a, b)
	if len(clipped) == 0 {
		return nil
	}
	pts := make([]Pt, 0, len(clipped)*2)
	for h := 0; h < len(clipped); h++ {
		x, y := clipped[h].Begin().XY()
		xequal := IsEqual(x, min.X()) || IsEqual(x, max.X())
		yequal := IsEqual(y, min.Y()) || IsEqual(y, max.Y())
		if xequal || yequal {
			pts = append(pts, clipped[h].Begin())
		}
		x, y = clipped[h].End().XY()
		xequal = IsEqual(x, min.X()) || IsEqual(x, max.X())
		yequal = IsEqual(y, min.Y()) || IsEqual(y, max.Y())
		if xequal || yequal {
			pts = append(pts, clipped[h].End())
		}
	}
	return pts
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

// Triangle represents a three sided polygon. This is independent from the
// polygon to support angluar math on it.
type Triangle struct {
	pts [3]Pt
}

func TrianglePt(p1, p2, p3 Pt) Triangle {
	return Triangle{
		pts: [3]Pt{p1, p2, p3},
	}
}
func (tri Triangle) Points() []Pt { return tri.pts[:] }
func (tri Triangle) Sides() []Segment {
	return []Segment{
		SegmentPt(tri.pts[0], tri.pts[1]),
		SegmentPt(tri.pts[1], tri.pts[2]),
		SegmentPt(tri.pts[2], tri.pts[0]),
	}
}
func (tri Triangle) String() string {
	return fmt.Sprintf("Triangle[ Polygon(%v, %v, %v) ]",
		tri.pts[0], tri.pts[1], tri.pts[2])
}

type Quadrilateral struct {
	pts [4]Pt
}

func QuadrilateralPt(p1, p2, p3, p4 Pt) Quadrilateral {
	return Quadrilateral{
		pts: [4]Pt{p1, p2, p3, p4},
	}
}
func (quad Quadrilateral) Points() []Pt { return quad.pts[:] }
func (quad Quadrilateral) Sides() []Segment {
	return PolygonPt(quad.Points()...).Sides()
}
func (quad Quadrilateral) String() string {
	return fmt.Sprintf("Quadrilateral[ Polygon(%v, %v, %v, %v) ]",
		quad.pts[0], quad.pts[1], quad.pts[2], quad.pts[3])
}

type Polygon struct {
	pts []Pt
}

func PolygonPt(pts ...Pt) Polygon {
	return Polygon{
		pts: pts,
	}
}
func (poly Polygon) Points() []Pt { return poly.pts[:] }
func (poly Polygon) Sides() []Segment {
	sides := make([]Segment, 0, len(poly.pts))
	prev := poly.pts[0]
	for h := 1; h < len(poly.pts); h++ {
		curr := poly.pts[h]
		sides = append(sides, SegmentPt(prev, curr))
		prev = curr
	}
	sides = append(sides, SegmentPt(prev, poly.pts[len(poly.pts)-1]))
	return sides
}

func IntersectionPolygonSegment(a Polygon, b Segment) []Pt {
	sides := a.Sides()
	ptset := make([]Pt, 0, 4)
	for _, aside := range sides {
		ptset = append(ptset, IntersectionSegmentSegment(aside, b)...)
	}
	if len(ptset) == 0 {
		return nil
	}

	ptset = SortPts(ptset)
	pts := make([]Pt, 1, len(ptset))
	pts[0] = ptset[0]
	for h := 1; h < len(ptset); h++ {
		if !IsEqualPair(pts[len(pts)-1], ptset[h]) {
			pts = append(pts, ptset[h])
		}
	}
	return pts
}
