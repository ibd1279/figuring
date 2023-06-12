package figuring

import (
	"math"
	"testing"
)

func TestLineIntersection(t *testing.T) {
	lineLineTests := []struct {
		a, b Line
		pts  []Pt
	}{
		{
			//0
			LineAbc(0, 2, 5), LineAbc(0.5, 100, -10),
			[]Pt{PtXy(520, -2.5)},
		}, {
			LineAbc(2, 0, 5), LineAbc(100, 0.5, -10),
			[]Pt{PtXy(-2.5, 520)},
		}, {
			LineAbc(0.5, 100, -10), LineAbc(0, 2, 5),
			[]Pt{PtXy(520, -2.5)},
		}, {
			LineAbc(100, 0.5, -10), LineAbc(2, 0, 5),
			[]Pt{PtXy(-2.5, 520)},
		}, {
			LineAbc(0, 0, 120), LineAbc(9, 10, 1000),
			nil,
		}, {
			//5
			LineAbc(9, 10, 1000), LineAbc(0, 0, 120),
			nil,
		}, {
			LineAbc(-10, 9, 0), LineAbc(-1, 0.9, 100),
			nil,
		}, {
			LineAbc(1, 2, 0), LineAbc(100, -30, 100),
			[]Pt{PtXy(-0.8695652173913, 0.4347826086957)},
		}, {
			LineAbc(-10, 2, 0.123), LineAbc(0.012354343, -1020, 1000),
			[]Pt{PtXy(0.2083789361539, 0.9803946807695)},
		}, {
			LineAbc(-10, 9, 0), LineAbc(9, 10, 1000),
			[]Pt{PtXy(-49.7237569060774, -55.2486187845304)},
		},
	}
	for h, test := range lineLineTests {
		a, b := test.a, test.b
		pts := IntersectionLineLine(a, b)
		if len(pts) != len(test.pts) {
			t.Fatalf("[%d]IntersectionLineLine(%v, %v) (length) failed. %v != %v",
				h, a, b, pts, test.pts)
		}
		for i := 0; i < len(pts); i++ {
			if !IsEqualPair(pts[i], test.pts[i]) {
				t.Errorf("[%d][%d]IntersectionLineLine(%v, %v) failed. %v != %v",
					h, i, a, b, pts[i], test.pts[i])
			}
		}
	}
}

func TestRayIntersection(t *testing.T) {
	rayRayTests := []struct {
		a, b         Ray
		rrpts, rlpts []Pt
	}{
		{
			//0
			RayFromVector(PtOrig, VectorFromTheta(1./2.)),
			RayFromVector(PtOrig, VectorFromTheta(-1./2.)),
			[]Pt{PtOrig},
			[]Pt{PtOrig},
		}, {
			RayFromVector(PtXy(5, 0), VectorFromTheta(3*math.Pi/4)),
			RayFromVector(PtOrig, VectorFromTheta(-1./2.)),
			[]Pt{},
			[]Pt{},
		}, {
			RayFromVector(PtXy(-5, 0), VectorFromTheta(3*math.Pi/4)),
			RayFromVector(PtOrig, VectorFromTheta(-1./2.)),
			[]Pt{},
			[]Pt{PtXy(-11.0205586151371, 6.0205586151371)},
		},
	}
	for h, test := range rayRayTests {
		a, b := test.a, test.b
		rrpts := IntersectionRayRay(a, b)
		if len(rrpts) != len(test.rrpts) {
			t.Fatalf("[%d]IntersectionRayRay(%v, %v) (length) failed. %v != %v",
				h, a, b, rrpts, test.rrpts)
		}
		for i := 0; i < len(rrpts); i++ {
			if !IsEqualPair(rrpts[i], test.rrpts[i]) {
				t.Errorf("[%d][%d]IntersectionRayRay(%v, %v) failed. %v != %v",
					h, i, a, b, rrpts[i], test.rrpts[i])
			}
		}

		rlpts := IntersectionLineRay(b.Line(), a)
		if len(rlpts) != len(test.rlpts) {
			t.Fatalf("[%d]IntersectionRayLine(%v, %v) (length) failed. %v != %v",
				h, a, b.Line(), rlpts, test.rlpts)
		}
		for i := 0; i < len(rlpts); i++ {
			if !IsEqualPair(rlpts[i], test.rlpts[i]) {
				t.Errorf("[%d][%d]IntersectionRayLine(%v, %v) failed. %v != %v",
					h, i, a, b.Line(), rlpts[i], test.rlpts[i])
			}
		}
	}

	raySegmentTests := []struct {
		a   Ray
		b   Segment
		pts []Pt
	}{
		{
			RayFromVector(PtOrig, VectorFromTheta(1./2.)),
			SegmentPt(PtXy(5, 0), PtXy(6, 15)),
			[]Pt{PtXy(5.1889836458308, 2.8347546874622)},
		}, {
			RayFromVector(PtXy(5, 0), VectorFromTheta(1./2.)),
			SegmentPt(PtXy(4, 0), PtXy(6, 15)),
			[]Pt{},
		},
	}
	for h, test := range raySegmentTests {
		a, b := test.a, test.b
		pts := IntersectionSegmentRay(b, a)
		if len(pts) != len(test.pts) {
			t.Fatalf("[%d]IntersectionRaySegment(%v, %v) (length) failed. %v != %v",
				h, a, b, pts, test.pts)
		}
		for i := 0; i < len(pts); i++ {
			if !IsEqualPair(pts[i], test.pts[i]) {
				t.Errorf("[%d][%d]IntersectionRaySegment(%v, %v) failed. %v != %v",
					h, i, a, b, pts[i], test.pts[i])
			}
		}
	}
}

func TestSegmentIntersection(t *testing.T) {
	segmentSegmentTests := []struct {
		a, b Segment
		pts  []Pt
	}{
		{
			//0
			SegmentPt(PtXy(0, 0), PtXy(51, 51)), SegmentPt(PtXy(100, 0), PtXy(49, 51)),
			[]Pt{PtXy(50, 50)},
		},
		{
			SegmentPt(PtXy(-10, 0), PtXy(100, 40)), SegmentPt(PtXy(100, 0), PtXy(49, 51)),
			[]Pt{PtXy(70.+2./3., 29.+1./3.)},
		},
		{
			SegmentPt(PtXy(-10, -100), PtXy(102, 1)), SegmentPt(PtXy(100, 0), PtXy(49, 51)),
			nil,
		},
		{
			SegmentPt(PtXy(-10, 100), PtXy(102, 100)), SegmentPt(PtXy(90, 100), PtXy(10, 100)),
			nil,
		},
		{
			SegmentPt(PtXy(-10, 10), PtXy(10, -10)), SegmentPt(PtXy(-15, 15), PtXy(15, -15)),
			nil,
		},
	}
	for h, test := range segmentSegmentTests {
		a, b := test.a, test.b
		pts := IntersectionSegmentSegment(a, b)
		if len(pts) != len(test.pts) {
			t.Fatalf("[%d]IntersectionSegmentSegment(%v, %v) (length) failed. %v != %v",
				h, a, b, pts, test.pts)
		}
		for i := 0; i < len(pts); i++ {
			if !IsEqualPair(pts[i], test.pts[i]) {
				t.Errorf("[%d][%d]IntersectionSegmentSegment(%v, %v) failed. %v != %v",
					h, i, a, b, pts[i], test.pts[i])
			}
		}

	}

	segmentLineTests := []struct {
		a   Segment
		b   Line
		pts []Pt
	}{
		{
			//0
			SegmentPt(PtXy(40, 60), PtXy(60, 40)),
			LineAbc(-10, 9, 0),
			[]Pt{PtXy(47.3684210526316, 52.6315789473684)},
		}, {
			SegmentPt(PtXy(20, 30), PtXy(40, 40)),
			LineAbc(-10, 9, 0),
			[]Pt{PtXy(32.7272727272727, 36.3636363636364)},
		}, {
			SegmentPt(PtXy(20, 60), PtXy(65, 80)),
			LineAbc(-10, 9, 0),
			nil,
		},
	}
	for h, test := range segmentLineTests {
		a, b := test.a, test.b
		pts := IntersectionLineSegment(b, a)
		if len(pts) != len(test.pts) {
			t.Fatalf("[%d]IntersectionSegmentLine(%v, %v) (length) failed. %v != %v",
				h, a, b, pts, test.pts)
		}
		for i := 0; i < len(pts); i++ {
			if !IsEqualPair(pts[i], test.pts[i]) {
				t.Errorf("[%d][%d]IntersectionSegmentLine(%v, %v) failed. %v != %v",
					h, i, a, b, pts[i], test.pts[i])
			}
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

	rectangleRectangleTests := []struct {
		a, b Rectangle
		c    []Rectangle
	}{
		{
			RectanglePt(PtXy(2, 2), PtXy(12, 4)),
			RectanglePt(PtXy(4, 1), PtXy(10, 5)),
			[]Rectangle{RectanglePt(PtXy(4, 2), PtXy(10, 4))},
		}, {
			RectanglePt(PtXy(2, 2), PtXy(12, 4)),
			RectanglePt(PtXy(4, 1), PtXy(15, 5)),
			[]Rectangle{RectanglePt(PtXy(4, 2), PtXy(12, 4))},
		}, {
			RectanglePt(PtXy(5, 2), PtXy(14, 4)),
			RectanglePt(PtXy(4, 1), PtXy(15, 5)),
			[]Rectangle{RectanglePt(PtXy(5, 2), PtXy(14, 4))},
		}, {
			RectanglePt(PtXy(5, 2), PtXy(14, 4)),
			RectanglePt(PtXy(24, 1), PtXy(15, 5)),
			nil,
		}, {
			RectanglePt(PtXy(5, 2), PtXy(14, 4)),
			RectanglePt(PtXy(4, 11), PtXy(15, 15)),
			nil,
		},
	}
	for h, test := range rectangleRectangleTests {
		a, b := test.a, test.b
		c := IntersectionRectangleRectangle(a, b)
		if len(c) != len(test.c) {
			t.Fatalf("[%d]IntersectionRectangleRectangle(%v, %v) (length) failed. %v != %v",
				h, a, b, c, test.c)
		}
		for i := 0; i < len(c); i++ {
			if !IsEqualPts(c[i], test.c[i]) {
				t.Errorf("[%d][%d]IntersectionRectangleRectangle(%v, %v) failed. %v != %v",
					h, i, a, b, c[i], test.c[i])
			}
		}
	}
}

/*

	BezierPt(PtXy(51, 113), PtXy(37, 245), PtXy(138, 245), PtXy(152, 150)),
	BezierPt(PtXy(110, 150), PtXy(25, 190), PtXy(210, 250), PtXy(210, 30)),
	BezierPt(PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80)),
	BezierPt(PtXy(285, 39), PtXy(129, 126), PtXy(248, 201), PtXy(127, 32)),
	BezierPt(PtXy(70, 250), PtXy(120, 15), PtXy(20, 95), PtXy(225, 80)),
*/

func TestIntersectionBezier(t *testing.T) {
	bezierBezierTests := []struct {
		a, b Bezier
		c    []Pt
	}{
		{
			BezierPt(PtXy(-10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10)),
			BezierPt(PtXy(-10, -10), PtXy(100, 400), PtXy(500, 450), PtXy(450, -100)),
			[]Pt{PtXy(-1.0417906819186, 20.6265054390535)},
		}, {
			BezierPt(PtXy(-10, 10), PtXy(10, 40), PtXy(50, 45), PtXy(45, -10)),
			BezierPt(PtXy(-0.10, -0.10), PtXy(1.2, 4.1), PtXy(0.5, 4.50), PtXy(-5.45, -0.1)),
			nil,
		}, {
			BezierPt(PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80)),
			BezierPt(PtXy(170, 140), PtXy(85, 180), PtXy(280, 250), PtXy(250, 30)),
			[]Pt{
				PtXy(193.5006244749729, 181.7476687150318),
				PtXy(217.2945210065356, 170.6953332130412),
				PtXy(249.9001783235052, 111.1241514682179),
				PtXy(252.032480333132, 96.4711903599781),
			},
		},
	}
	for h, test := range bezierBezierTests {
		a, b := test.a, test.b
		c := IntersectionBezierBezier(a, b)
		if len(c) != len(test.c) {
			t.Fatalf("[%d]IntersectionBezierBezier(%v, %v) (length) failed. %v != %v",
				h, a, b, c, test.c)
		}
		// Scale by 10, to round to the closest 10th.
		scalar := VectorIj(10, 10)
		c = ScalePts(scalar, c)
		test.c = ScalePts(scalar, test.c)
		for i := 0; i < len(c); i++ {
			if !IsEqual(c[i].X().Round(), test.c[i].X().Round()) || !IsEqual(c[i].Y().Round(), test.c[i].Y().Round()) {
				t.Errorf("[%d][%d]IntersectionBezierBezier(%v, %v) failed. %v != %v",
					h, i, a, b, c[i], test.c[i])
			}
		}
	}
}

func BenchmarkIntersectionBezierBezier(b *testing.B) {
	b1 := BezierPt(PtXy(396, 34), PtXy(89, 120), PtXy(199, 295), PtXy(260, 80))
	b2 := BezierPt(PtXy(170, 140), PtXy(85, 180), PtXy(280, 250), PtXy(250, 30))
	for h := 0; h < b.N; h++ {
		IntersectionBezierBezier(b1, b2)
	}
}
