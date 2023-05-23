package figuring

import (
	"math"
	"testing"
)

func TestLine(t *testing.T) {
	identityTests := []struct {
		a           Line
		s           string
		xi, yi      Length
		horiz, vert bool
	}{
		{LineAbc(2, 0, 5), "2x=5", 2.5, Length(math.NaN()), false, true},
		{LineAbc(0, 2, 5), "2y=5", Length(math.NaN()), 2.5, true, false},
		{LineAbc(3, 5, 7), "3x+5y=7", 7. / 3., 1.4, false, false},
		{LineFromPt(PtXy(2, 3), PtXy(4, 4)), "1x-2y=-4", -4, 2, false, false},
		{LineFromVector(PtXy(1, 1), VectorIj(2, 5)), "5x-2y=3", 0.6, -1.5, false, false},
		{LineAbc(0, 0, 12), "0x+0y=12", Length(math.NaN()), Length(math.NaN()), false, false},
	}
	for h, test := range identityTests {
		a := test.a
		if s := a.String(); s != test.s {
			t.Errorf("[%d](%s).String() failed. %s != %s",
				h, a, s, test.s)
		}

		ti, terr := test.xi.OrErr()
		if xi, err := a.XForY(0).OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).XForY(0) failed (error). %v != %v",
				h, a, err, terr)
		} else if terr == nil && !IsEqual(xi, ti) {
			t.Errorf("[%d](%s).XForY(0) failed. %.8f != %.8f",
				h, a, xi, ti)
		}

		ti, terr = test.yi.OrErr()
		if yi, err := a.YForX(0).OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).YForX(0) failed (error). %v != %v",
				h, a, err, terr)
		} else if terr == nil && !IsEqual(yi, ti) {
			t.Errorf("[%d](%s).YForX(0) failed. %.8f != %.8f",
				h, a, yi, ti)
		}

		if ish := a.IsHorizontal(); ish != test.horiz {
			t.Errorf("[%d](%s).IsHorizontal() failed. %t != %t",
				h, a, ish, test.horiz)
		}
		if isv := a.IsVertical(); isv != test.vert {
			t.Errorf("[%d](%s).IsVertical() failed. %t != %t",
				h, a, isv, test.vert)
		}
	}

	normalizeTests := []struct {
		a     Line
		bxerr bool
		bx    Line
		byerr bool
		by    Line
		v     Vector
		rads  Radians
	}{
		{
			LineAbc(2, 2, 0),
			false, LineAbc(1, 1, 0),
			false, LineAbc(1, 1, 0),
			VectorIj(-1, 1).Normalize(),
			3 * math.Pi / 4,
		},
		{
			LineAbc(6, 2, 2),
			false, LineAbc(1, 1./3., 1./3.),
			false, LineAbc(3, 1, 1),
			VectorIj(-1, 3).Normalize(),
			1.8925468811868438,
		},
		{
			LineAbc(14, -42, 7),
			false, LineAbc(1, -3, 0.5),
			false, LineAbc(-1./3., 1, -5./30.),
			VectorIj(0.5, 5./30.).Normalize(),
			0.3217505543958438,
		},
	}
	for h, test := range normalizeTests {
		a, terr := test.a.OrErr()
		if b, err := a.NormalizeX().OrErr(); (err != nil) != test.bxerr {
			t.Errorf("[%d](%s).NormalizeX() failed (error). %v != %v",
				h, a, err, terr)
		} else if !IsEqualEquations(b, test.bx) {
			t.Errorf("[%d](%s).NormalizeX() failed. %v != %v",
				h, a, b, test.bx)
		}

		if b, err := a.NormalizeY().OrErr(); (err != nil) != test.byerr {
			t.Errorf("[%d](%s).NormalizeY() failed (error). %v != %v",
				h, a, err, terr)
		} else if !IsEqualEquations(b, test.by) {
			t.Errorf("[%d](%s).NormalizeY() failed. %v != %v",
				h, a, b, test.by)
		}

		if v, err := a.Vector().OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).Vector() failed (error). %v != %v",
				h, a, err, terr)
		} else if !IsEqualPair(v, test.v) {
			t.Errorf("[%d](%s).Vector() failed. %v != %v",
				h, a, v, test.v)
		}

		if rads, err := a.Angle().OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).Angle() failed (error). %v != %v",
				h, a, err, terr)
		} else if !IsEqual(rads, test.rads) {
			t.Errorf("[%d](%s).Angle() failed. %f != %f",
				h, a, rads, test.rads)
		}
	}

	errorTests := []struct {
		a     Line
		isErr bool
	}{
		{LineAbc(0, 0, 12), true},
		{LineAbc(3, 4, 12), false},
		{LineAbc(Length(math.NaN()), 4, 12), true},
		{LineAbc(3, Length(math.Inf(1)), 12), true},
		{LineAbc(Length(math.Inf(-1)), 4, 12), true},
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

func TestLineIntersection(t *testing.T) {
	lineLineTests := []struct {
		a, b Line
		pts  []Pt
	}{
		{
			//0
			LineAbc(0, 2, 5), LineAbc(0.5, 100, -10),
			[]Pt{PtXy(-520, 2.5)},
		}, {
			LineAbc(2, 0, 5), LineAbc(100, 0.5, -10),
			[]Pt{PtXy(2.5, -520)},
		}, {
			LineAbc(0.5, 100, -10), LineAbc(0, 2, 5),
			[]Pt{PtXy(-520, 2.5)},
		}, {
			LineAbc(100, 0.5, -10), LineAbc(2, 0, 5),
			[]Pt{PtXy(2.5, -520)},
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
			[]Pt{PtXy(0.8695652173913, -0.4347826086957)},
		}, {
			LineAbc(-10, 2, 0.123), LineAbc(0.012354343, -1020, 1000),
			[]Pt{PtXy(-0.2083789361539, -0.9803946807695)},
		}, {
			LineAbc(-10, 9, 0), LineAbc(9, 10, 1000),
			[]Pt{PtXy(49.7237569060774, 55.2486187845304)},
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

func TestSegment(t *testing.T) {
	identityTests := []struct {
		a          Segment
		s          string
		begin, end Pt
		length     Length
		angle      Radians
	}{
		{
			SegmentPt(PtXy(5, 0), PtXy(0, 5)),
			"Segment(Point({5, 0}), Point({0, 5}))",
			PtXy(5, 0), PtXy(0, 5),
			7.0710678118655, 3. * math.Pi / 4.,
		}, {
			SegmentPt(PtXy(5, 5), PtXy(0, 5)),
			"Segment(Point({5, 5}), Point({0, 5}))",
			PtXy(5, 5), PtXy(0, 5),
			5, math.Pi,
		}, {
			SegmentPt(PtXy(5, 5), PtXy(Length(math.NaN()), 5)),
			"Segment(Point({5, 5}), Point({NaN, 5}))",
			PtXy(5, 5), PtXy(Length(math.NaN()), 5),
			0, 0,
		},
	}
	for h, test := range identityTests {
		a := test.a
		if s := a.String(); s != test.s {
			t.Errorf("[%d](%s).String() failed. %s != %s",
				h, a, s, test.s)
		}

		tp, terr := test.begin.OrErr()
		if p, err := a.Begin().OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).Begin() failed (error). %v != %v",
				h, a, err, terr)
		} else if terr == nil && !IsEqualPair(p, tp) {
			t.Errorf("[%d](%s).Begin() failed. %v != %v",
				h, a, p, tp)
		}

		tp, terr = test.end.OrErr()
		if p, err := a.End().OrErr(); (err == nil) != (terr == nil) {
			t.Errorf("[%d](%s).End() failed (error). %v != %v",
				h, a, err, terr)
		} else if terr == nil && !IsEqualPair(p, tp) {
			t.Errorf("[%d](%s).End() failed. %v != %v",
				h, a, p, tp)
		}

		if _, err := a.OrErr(); err == nil {
			if length := a.Length(); !IsEqual(length, test.length) {
				t.Errorf("[%d](%s).Length() failed. %f != %f",
					h, a, length, test.length)
			}
			if angle := a.Angle(); !IsEqual(angle, test.angle) {
				t.Errorf("[%d](%s).Length() failed. %f != %f",
					h, a, angle, test.angle)
			}
		}
	}

	reverseTests := []struct {
		a Segment
		r Segment
	}{
		{SegmentPt(PtXy(0, 5), PtXy(5, 0)), SegmentPt(PtXy(5, 0), PtXy(0, 5))},
		{SegmentPt(PtXy(20, 5), PtXy(5, 2)), SegmentPt(PtXy(5, 2), PtXy(20, 5))},
	}
	for h, test := range reverseTests {
		a := test.a
		r := a.Reverse()
		if IsEqualPts(a, r) {
			t.Errorf("[%d](%s).Reverse() failed (matched source). %v == %v",
				h, a, r, test.r)
		}
		if !IsEqualPts(r, test.r) {
			t.Errorf("[%d](%s).Reverse() failed. %v == %v",
				h, a, r, test.r)
		}
	}

	errorTests := []struct {
		a     Segment
		isErr bool
	}{
		{SegmentPt(PtXy(0, 0), PtXy(0, 0)), false},
		{SegmentPt(PtXy(120, 12), PtXy(455, 30)), false},
		{SegmentPt(PtXy(0, Length(math.NaN())), PtXy(0, 0)), true},
		{SegmentPt(PtXy(Length(math.Inf(-1)), 3), PtXy(3, 3)), true},
		{SegmentPt(PtXy(3, 3), PtXy(3, Length(math.Inf(1)))), true},
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
		pts := IntersectionSegmentLine(a, b)
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
