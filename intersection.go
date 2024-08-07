package figuring

import (
	"math"
)

// --- Line Dominant Intersections ---

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
		ba, bb = -ba, -bb

		x := Length((bb - ba) / (mb - ma))
		y := b.YForX(x)

		p = PtXy(x, y)
	}

	return []Pt{p}
}

// IntersectionLineRay returns the intersection points of a line and a
// ray. Returns an empty slice if the two do not intersect.
func IntersectionLineRay(a Line, b Ray) []Pt {
	bLine := b.Line()
	pts := FilterPtsRay(b, IntersectionLineLine(a, bLine))
	if len(pts) == 0 {
		return nil
	}

	return pts
}

// IntersectionLineSegment returns the intersection points of a line and a
// segment. Returns an empty slice if the two do not intersect.
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

// IntersectionLineBezier returns the intersection points of a line and a
// bezier. Returns an empty slice if the two do not intersect.
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

// IntersectionRayRay returns the intersection points of two rays
// Returns an empty slice if the two do not intersect.
func IntersectionRayRay(a Ray, b Ray) []Pt {
	aLine := a.Line()
	bLine := b.Line()
	pts := FilterPtsRay(a, FilterPtsRay(b, IntersectionLineLine(aLine, bLine)))
	if len(pts) == 0 {
		return nil
	}

	return pts
}

// --- Segment Dominant Intersections ---

// IntersectionSegmentSegment returns the intersection points of two segments.
// Returns an empty slice if the two do not intersect.
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

// IntersectionSegmentRay returns the intersection points of a segment and a
// ray. Returns an empty slice if the two do not intersect.
func IntersectionSegmentRay(a Segment, b Ray) []Pt {
	bLine := b.Line()
	pts := FilterPtsRay(b, IntersectionLineSegment(bLine, a))
	if len(pts) == 0 {
		return nil
	}

	return pts
}

// IntersectionSegmentBezier returns the intersection points of a segment and a
// bezier. Returns an empty slice if the two do not intersect.
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

// --- Rectangle Dominant Intersections ---

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

func IntersectionRectangleRectangle(a Rectangle, b Rectangle) []Rectangle {
	overlap := func(amax, bmax Length) Length {
		if bmax < amax {
			return bmax
		}
		return amax
	}

	var lx, mx Length
	switch {
	case IsEqual(a.MinPt().X(), b.MinPt().X()):
		lx = a.MinPt().X()
		mx = overlap(a.MaxPt().X(), b.MaxPt().X())
	case b.MinPt().X() < a.MinPt().X():
		a, b = b, a
		fallthrough
	case a.MinPt().X() < b.MinPt().X():
		if b.MinPt().X() > a.MaxPt().X() {
			return nil
		}
		lx = b.MinPt().X()
		mx = overlap(a.MaxPt().X(), b.MaxPt().X())
	}

	var ly, my Length
	switch {
	case IsEqual(a.MinPt().Y(), b.MinPt().Y()):
		ly = a.MinPt().Y()
		my = overlap(a.MaxPt().Y(), b.MaxPt().Y())
	case b.MinPt().Y() < a.MinPt().Y():
		a, b = b, a
		fallthrough
	case a.MinPt().Y() < b.MinPt().Y():
		if b.MinPt().Y() > a.MaxPt().Y() {
			return nil
		}
		ly = b.MinPt().Y()
		my = overlap(a.MaxPt().Y(), b.MaxPt().Y())
	}

	return []Rectangle{RectanglePt(PtXy(lx, ly), PtXy(mx, my))}
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

	SortPts(ptset)
	pts := make([]Pt, 1, len(ptset))
	pts[0] = ptset[0]
	for h := 1; h < len(ptset); h++ {
		if !IsEqualPair(pts[len(pts)-1], ptset[h]) {
			pts = append(pts, ptset[h])
		}
	}
	return pts
}

// --- Bezier Dominant Intersections ---

func IntersectionBezierBezier(a, b Bezier) []Pt {
	type combination struct {
		a, b Pt
	}

	var xsectfunc func(Bezier, Bezier) []combination
	xsectfunc = func(a, b Bezier) []combination {
		var combos []combination
		abox, bbox := a.FastBox(), b.FastBox()
		xsect := IntersectionRectangleRectangle(abox, bbox)
		if len(xsect) > 0 {
			aw, ah := abox.Dims()
			bw, bh := bbox.Dims()
			if aw < 0.005 && ah < 0.005 && bw < 0.005 && bh < 0.005 {
				return []combination{combination{
					a.PtAtT(0.5), b.PtAtT(0.5),
				}}
			}
			a1, a2 := a.SplitAtT(0.5)
			b1, b2 := b.SplitAtT(0.5)
			combos = append(combos, xsectfunc(a1, b1)...)
			combos = append(combos, xsectfunc(a1, b2)...)
			combos = append(combos, xsectfunc(a2, b1)...)
			combos = append(combos, xsectfunc(a2, b2)...)
		}
		return combos
	}

	buffer := xsectfunc(a, b)

	if len(buffer) == 0 {
		return nil
	}

	ap := buffer[0].a
	bp := buffer[0].b
	lastx, lasty := ap.XY()
	dist := ap.VectorTo(bp).Magnitude()
	ret := []Pt{ap}
	for _, pair := range buffer {
		ap = pair.a
		bp = pair.b
		x, y := ap.XY()
		if math.Abs(float64(x-lastx)) < 0.05 && math.Abs(float64(y-lasty)) < 0.05 {
			newDist := ap.VectorTo(bp).Magnitude()
			if newDist < dist {
				lastx, lasty = x, y
				ret[len(ret)-1] = ap
				dist = newDist
			}
		} else {
			lastx, lasty = x, y
			dist = ap.VectorTo(bp).Magnitude()
			ret = append(ret, ap)
		}
	}
	SortPts(ret)
	return ret
}
