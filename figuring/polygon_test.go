package figuring

import (
	"math"
	"testing"
)

func TestRectangle(t *testing.T) {
	identityTests := []struct {
		a        Rectangle
		s        string
		min, max Pt
		w, h     Length
	}{
		{
			//0
			RectanglePt(PtXy(2, -2), PtXy(-2, 2)),
			"Rectangle[ Polygon(Point({-2, -2}), Point({-2, 2}), Point({2, 2}), Point({2, -2})) ]",
			PtXy(-2, -2), PtXy(2, 2),
			4, 4,
		},
	}
	for h, test := range identityTests {
		a := test.a
		if s := a.String(); s != test.s {
			t.Errorf("[%d](%s).String() failed. %s != %s",
				h, a, s, test.s)
		}
		if min := a.MinPt(); !IsEqualPair(min, test.min) {
			t.Errorf("[%d](%s).MinPt() failed. %v != %v",
				h, a, min, test.min)
		}
		if max := a.MaxPt(); !IsEqualPair(max, test.max) {
			t.Errorf("[%d](%s).MaxPt() failed. %v != %v",
				h, a, max, test.max)
		}
		if width := a.Width(); !IsEqual(width, test.w) {
			t.Errorf("[%d](%s).Width() failed. %f != %f",
				h, a, width, test.w)
		}
		if height := a.Height(); !IsEqual(height, test.h) {
			t.Errorf("[%d](%s).Height() failed. %f != %f",
				h, a, height, test.h)
		}
		if width, height := a.Dims(); !IsEqual(width, test.w) || !IsEqual(height, test.h) {
			t.Errorf("[%d](%s).Dims() failed. (%f, %f) != (%f, %f)",
				h, a, width, height, test.w, test.h)
		}

	}

	errorTests := []struct {
		a     Rectangle
		isErr bool
	}{
		{RectanglePt(PtXy(1, 1), PtXy(5, 5)), false},
		{RectanglePt(PtXy(-1, -1), PtXy(-5, -5)), false},
		{RectanglePt(PtXy(Length(math.NaN()), 1), PtXy(5, 5)), true},
		{RectanglePt(PtXy(1, 1), PtXy(5, Length(math.NaN()))), true},
		{RectanglePt(PtXy(1, Length(math.Inf(1))), PtXy(5, 5)), true},
		{RectanglePt(PtXy(1, 1), PtXy(Length(math.Inf(-1)), 5)), true},
	}
	for h, test := range errorTests {
		a := test.a
		_, err := a.OrErr()
		if (err != nil) != test.isErr {
			t.Errorf("[%d](%v).OrErr() failed. %t != %t. %v",
				h, test.a, (err != nil), test.isErr, err)
		}
	}
}

func TestIntersectionRectangle(t *testing.T) {
	rectangleLineTests := []struct {
		a   Rectangle
		b   Line
		pts []Pt
	}{
		{
			//0
			RectanglePt(PtXy(1, 1), PtXy(5, 5)),
			LineFromPt(PtOrig, PtXy(6, 6)),
			[]Pt{PtXy(1, 1), PtXy(5, 5)},
		}, {
			RectanglePt(PtXy(1, 1), PtXy(5, 5)),
			LineFromPt(PtXy(2, 0), PtXy(4, 6)),
			[]Pt{PtXy(7./3., 1), PtXy(11./3., 5)},
		},
	}
	for h, test := range rectangleLineTests {
		a := test.a
		b := test.b
		pts := IntersectionRectangleLine(a, b)
		if len(pts) != len(test.pts) {
			t.Fatalf("[%d]IntersectionRectangleLine(%v, %v) (length) failed. %v != %v",
				h, a, b, pts, test.pts)
		}
		for i := 0; i < len(pts); i++ {
			if !IsEqualPair(pts[i], test.pts[i]) {
				t.Errorf("[%d][%d]IntersectionRectangleLine(%v, %v) failed. %v != %v",
					h, i, a, b, pts[i], test.pts[i])
			}
		}
	}

	rectangleSegmentTests := []struct {
		a   Rectangle
		b   Segment
		pts []Pt
	}{
		{
			//0
			RectanglePt(PtXy(1, 1), PtXy(5, 5)),
			SegmentPt(PtOrig, PtXy(6, 6)),
			[]Pt{PtXy(1, 1), PtXy(5, 5)},
		}, {
			RectanglePt(PtXy(1, 1), PtXy(5, 5)),
			SegmentPt(PtXy(2, 0), PtXy(4, 6)),
			[]Pt{PtXy(7./3., 1), PtXy(11./3., 5)},
		},
	}
	for h, test := range rectangleSegmentTests {
		a := test.a
		b := test.b
		pts := IntersectionRectangleSegment(a, b)
		if len(pts) != len(test.pts) {
			t.Fatalf("[%d]IntersectionRectangleSegment(%v, %v) (length) failed. %v != %v",
				h, a, b, pts, test.pts)
		}
		for i := 0; i < len(pts); i++ {
			if !IsEqualPair(pts[i], test.pts[i]) {
				t.Errorf("[%d][%d]IntersectionRectangleSegment(%v, %v) failed. %v != %v",
					h, i, a, b, pts[i], test.pts[i])
			}
		}
	}
}

func BenchmarkRectangleLine(b *testing.B) {
	a := RectanglePt(PtXy(1, 1), PtXy(5, 5))
	c := LineFromPt(PtOrig, PtXy(6, 6))
	for h := 0; h < b.N; h++ {
		IntersectionRectangleLine(a, c)
	}
}

func BenchmarkRectangleSegment(b *testing.B) {
	a := RectanglePt(PtXy(1, 1), PtXy(5, 5))
	c := SegmentPt(PtOrig, PtXy(6, 6))
	for h := 0; h < b.N; h++ {
		IntersectionRectangleSegment(a, c)
	}
}
