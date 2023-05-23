package figuring

import "fmt"

type OrderedPtser interface {
	Points() []Pt
}

// Rectangle represents an axis aligned rectangle. The resulting rectangle will
// always be aligned with the X and Y axis.
type Rectangle struct {
	pts [2]Pt
}

func RectanglePt(p1, p2 Pt) Rectangle {
	lx, mx, ly, my := LimitsPts([]Pt{p1, p2})
	return Rectangle{
		pts: [2]Pt{PtXy(lx, ly), PtXy(mx, my)},
	}
}
func (r Rectangle) MinPt() Pt    { return r.pts[0] }
func (r Rectangle) MaxPt() Pt    { return r.pts[1] }
func (r Rectangle) Points() []Pt { return r.pts[:] }
func (r Rectangle) Dims() (Length, Length) {
	return r.pts[0].VectorTo(r.pts[1]).Units()
}
func (r Rectangle) Width() Length {
	w, _ := r.Dims()
	return w
}
func (r Rectangle) Height() Length {
	_, h := r.Dims()
	return h
}
func (r Rectangle) OrErr() (Rectangle, *FloatingPointError) {
	if _, err := r.pts[0].OrErr(); err != nil {
		return r, err
	} else if _, err = r.pts[1].OrErr(); err != nil {
		return r, err
	}
	return r, nil
}
func (r Rectangle) String() string {
	minmax, maxmin := PtXy(r.pts[0].X(), r.pts[1].Y()), PtXy(r.pts[1].X(), r.pts[0].Y())
	return fmt.Sprintf("rect=Polygon(%v, %v, %v, %v)",
		r.pts[0], minmax, r.pts[1], maxmin)
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
func (t Triangle) Points() []Pt { return t.pts[:] }

type Quadrilateral struct {
	pts [4]Pt
}

func QuadrilateralPt(p1, p2, p3, p4 Pt) Quadrilateral {
	return Quadrilateral{
		pts: [4]Pt{p1, p2, p3, p4},
	}
}
func (quad Quadrilateral) Points() []Pt { return quad.pts[:] }

type Polygon struct {
	pts []Pt
}
