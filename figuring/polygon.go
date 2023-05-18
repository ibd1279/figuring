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
